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
	"github.com/HenryGits/hadoop-operator/controllers/tools"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
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

var HadoopTpl = tools.Parser{
	Directory: tools.EnvVar("GT_TEMPLATE_PATH", "/etc/operator/templates"),
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
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods;services;persistentvolumes;persistentvolumeclaims;configmaps;secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.istio.io,resources=destinationrules,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=policy,resources=poddisruptionbudgets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=*

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *HadoopReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("==========HadoopReconciler begin to Reconcile==========")

	// 获取 Hadoop 实例
	var originHadoop = &hadoopv1.Hadoop{}
	if err := r.Get(ctx, req.NamespacedName, originHadoop); err != nil {
		logger.Error(err, "unable to fetch hadoop")
		// 停止协调
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	hadoop := originHadoop.DeepCopy()
	hadoop.Status.Phase = hadoopv1.Reconciling

	// examine DeletionTimestamp to determine if object is under deletion
	if hadoop.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !tools.ContainsString(hadoop.GetFinalizers(), Finalizer) {
			controllerutil.AddFinalizer(hadoop, Finalizer)
			if err := r.Update(ctx, hadoop); err != nil {
				logger.Error(err, "unable to update hadoop")
				return ctrl.Result{}, err
			}
		}

		runtimeObjects, err := r.generateRuntimeObjects(hadoop)
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
			object.SetNamespace(hadoop.Namespace)
			// set controller reference
			if err := controllerutil.SetControllerReference(hadoop, object, r.Scheme); err != nil {
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
		if tools.ContainsString(hadoop.GetFinalizers(), Finalizer) {
			// our finalizer is present, so lets handle any external dependency
			//if err := r.deleteExternalResources(ctx, hadoop); err != nil {
			//	// if fail to delete the external dependency here, return with error
			//	// so that it can be retried
			//	return ctrl.Result{}, err
			//}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(hadoop, Finalizer)
			if err := r.Update(ctx, hadoop); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	if !reflect.DeepEqual(hadoop.TypeMeta, hadoop.TypeMeta) || !reflect.DeepEqual(hadoop.ObjectMeta, hadoop.ObjectMeta) || !reflect.DeepEqual(hadoop.Spec, hadoop.Spec) {
		if err := r.Update(ctx, hadoop); err != nil {
			logger.Error(err, "update without status error")
			return ctrl.Result{}, err
		}
	}

	if !reflect.DeepEqual(hadoop.Status, hadoop.Status) {
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
		Owns(&corev1.ConfigMap{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&appsv1.DaemonSet{}).
		Complete(r)
}

func (r *HadoopReconciler) generateRuntimeObjects(hadoop *hadoopv1.Hadoop) (runtimeObjects []*runtime.Object, err error) {
	templates, err := HadoopTpl.ParseTemplate("hadoop.dameng.com_hadoop.gotmpl", hadoop)
	if err != nil {
		klog.Errorf("generate hadoop runtime error: %v", err)
		return runtimeObjects, err
	}
	klog.V(8).Infof("template: %s", templates)
	return tools.ParseYaml(templates)
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

func objectEqual(hadoop, newer *unstructured.Unstructured) (bool, error) {
	newSpec, found, err := unstructured.NestedFieldNoCopy(newer.Object, "spec")
	if err != nil {
		klog.Errorf("nested spec error: %v", err)
		return false, err
	}
	if !found {
		return false, nil
	}

	originSpecJSON, ok := hadoop.GetAnnotations()["operator.dameng.com/spec"]
	if !ok {
		return false, nil
	}

	klog.V(6).Infof("hadoop spec json: %v", originSpecJSON)

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
