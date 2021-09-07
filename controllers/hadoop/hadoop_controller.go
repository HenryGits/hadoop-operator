/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hadoop

import (
	"context"
	"github.com/HenryGits/hadoop-operator/pkg/tools"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/third_party/forked/golang/template"
	"k8s.io/klog/v2"
	"os"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	hadoopv1 "github.com/HenryGits/hadoop-operator/apis/hadoop/v1"
)

const (
	FieldManager  string = "hadoop-operator"
	Finalizer     string = "finalizer.hadoop.operator.dameng.com"
	ReleasePrefix string = "hadoop-"
)

var HadoopTpl = template.Parser{
	Directory: os.EnvVar("TEMPLATES_PATH", "/etc/dmcca/templates"),
	Pattern:   "\\.gotmpl$",
}

// HadoopReconciler reconciles a Hadoop object
type HadoopReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=hadoop.dameng.com,resources=hadoops,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=hadoop.dameng.com,resources=hadoops/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=hadoop.dameng.com,resources=hadoops/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Hadoop object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *HadoopReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("==========HadoopReconciler begin to Reconcile==========")

	var origin = &hadoopv1.Hadoop{}
	if err := r.Get(ctx, req.NamespacedName, origin); err != nil {
		logger.Error(err, "unable to fetch hadoop")
		// 停止协调
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	hadoop := origin.DeepCopy()
	hadoop.Status.Phase = hadoopv1.Reconciling

	// examine DeletionTimestamp to determine if object is under deletion
	if origin.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !tools.ContainsString(origin.GetFinalizers(), Finalizer) {
			controllerutil.AddFinalizer(origin, Finalizer)
			if err := r.Update(ctx, origin); err != nil {
				logger.Error(err, "unable to update hadoop")
				return ctrl.Result{}, err
			}
		}

		runtimeObjects, err := r.generateRuntimeObjects(origin)
		if err != nil {
			logger.Error(err, "decode error")
			return ctrl.Result{}, err
		}

		for _, object := range runtimeObjects {
			// set spec annotation
			object, err := setSpecAnnotation((*object).(*unstructured.Unstructured))
			if err != nil {
				logger.Error(err, "set annotation error")
				return ctrl.Result{}, err
			}
			// set namespace
			object.SetNamespace(origin.Namespace)
			// set controller reference
			if err := controllerutil.SetControllerReference(origin, object, r.Scheme); err != nil {
				logger.Error(err, "maintain hadoop controller reference error")
				return ctrl.Result{}, err
			}

			logger.V(6).Info("object content:", "object", object)

			// retrieve origin
			var originObject unstructured.Unstructured
			originObject.SetGroupVersionKind(object.GroupVersionKind())
			if err := r.Get(ctx, client.ObjectKey{Namespace: hadoop.Namespace, Name: object.GetName()}, &originObject); err != nil {
				// create object if not found
				if errors.IsNotFound(err) {
					if err := r.Create(ctx, object, &client.CreateOptions{FieldManager: FieldManager}); err != nil {
						logger.Error(err, "Object create error", "Object", object)

						return ctrl.Result{}, err
					}
				} else {
					logger.Error(err, "get origin Object error")
					return ctrl.Result{}, err
				}
			} else {
				// continue if origin equal new
				equal, err := objectEqual(&originObject, object)
				if err != nil {
					logger.Error(err, "Object equal error")
					return ctrl.Result{}, err
				}
				logger.V(4).Info("deep equal", "kind", object.GroupVersionKind(), "result", equal)
				if equal {
					continue
				}

				// patch if origin not equal new
				if err := r.Patch(ctx, object, client.Merge, &client.PatchOptions{FieldManager: FieldManager}); err != nil {
					logger.Error(err, "Object patch error", "Object", object)

					return ctrl.Result{}, err
				}
			}
		}
	} else {
		// The object is being deleted
		if tools.ContainsString(origin.GetFinalizers(), Finalizer) {
			// our finalizer is present, so lets handle any external dependency
			if err := r.deleteExternalResources(ctx, origin); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(origin, Finalizer)
			if err := r.Update(ctx, origin); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	if !reflect.DeepEqual(hadoop.TypeMeta, origin.TypeMeta) || !reflect.DeepEqual(hadoop.ObjectMeta, origin.ObjectMeta) || !reflect.DeepEqual(hadoop.Spec, origin.Spec) {
		if err := r.Update(ctx, hadoop); err != nil {
			logger.Error(err, "update without status error")
			return ctrl.Result{}, err
		}
	}

	if !reflect.DeepEqual(hadoop.Status, origin.Status) {
		if err := r.Status().Update(ctx, hadoop); err != nil {
			logger.Error(err, "update status error")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HadoopReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hadoopv1.Hadoop{}).
		Complete(r)
}

func (r *HadoopReconciler) deleteExternalResources(ctx context.Context, hadoop *hadoopv1.Hadoop) error {
	//
	// delete any external resources associated with the cronJob
	//
	// Ensure that delete implementation is idempotent and safe to invoke
	// multiple times for same object.

	// volumeClaimTemplates 声明的存储不能主动删除，需要手动删除
	if hadoop.Spec.Persistence.DataNode.Enabled {
		name := hadoop.ObjectMeta.Name
		//
		//pvcList := &corev1.PersistentVolumeClaimList{}
		//listOpts := []client.ListOption{
		//	client.InNamespace(hadoop.Namespace),
		//	client.MatchingLabels(map[string]string{
		//		"app":       "hadoop",
		//		"component": "hdfs-dn",
		//		"release":   ReleasePrefix + name,
		//	}),
		//}
		//if err := r.List(ctx, pvcList, listOpts...); err != nil {
		//	klog.Errorf("get pvc of hadoop datanode %s error: %v", name, err)
		//	return err
		//}

		pvc := &corev1.PersistentVolumeClaim{}
		opts := []client.DeleteAllOfOption{
			client.InNamespace(hadoop.Namespace),
			client.MatchingLabels{
				"app":       "hadoop",
				"component": "hdfs-dn",
				"release":   ReleasePrefix + name,
			},
			client.GracePeriodSeconds(5),
		}
		if err := r.DeleteAllOf(ctx, pvc, opts...); err != nil {
			klog.Errorf("delete pvc of hadoop datanode %s error: %v", name, err)
			return err
		}
	}

	return nil
}

func (r *HadoopReconciler) generateRuntimeObjects(hadoop *hadoopv1.Hadoop) (runtimeObjects []*runtime.Object, err error) {
	templates, err := HadoopTpl.ParseTemplate("hadoop.dameng.com_hadoop.gotmpl", hadoop)
	if err != nil {
		klog.Errorf("generate hadoop runtime error: %v", err)
		return runtimeObjects, err
	}
	klog.V(8).Infof("template: %s", templates)
	return kubernetes.ParseYaml(templates)
}

func setSpecAnnotation(object *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	spec, found, err := unstructured.NestedFieldNoCopy(object.Object, "spec")
	if err != nil {
		klog.Errorf("nested spec error: %v", err)
		return nil, err
	}
	if !found {
		return object, nil
	}
	annotation, err := json.Marshal(spec)
	if err != nil {
		klog.Errorf("marshal spec error: %v", err)
		return nil, err
	}
	annotations := object.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations["operator.dameng.com/spec"] = string(annotation)
	object.SetAnnotations(annotations)
	return object, nil
}

func objectEqual(origin, newer *unstructured.Unstructured) (bool, error) {
	newSpec, found, err := unstructured.NestedFieldNoCopy(newer.Object, "spec")
	if err != nil {
		klog.Errorf("nested spec error: %v", err)
		return false, err
	}
	if !found {
		return false, nil
	}

	originSpecJSON, ok := origin.GetAnnotations()["operator.dameng.com/spec"]
	if !ok {
		return false, nil
	}

	klog.V(6).Infof("origin spec json: %v", originSpecJSON)

	var originSpec map[string]interface{}
	if err := json.Unmarshal([]byte(originSpecJSON), &originSpec); err != nil {
		klog.Errorf("unmarshal error: %v", err)
		return false, err
	}

	if reflect.DeepEqual(newSpec, originSpec) {
		return true, nil
	}
	return false, nil
}
