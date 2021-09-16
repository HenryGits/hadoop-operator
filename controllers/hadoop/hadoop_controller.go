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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
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
	Directory: tools.EnvVar("TEMPLATES_PATH", "/etc/operator/templates"),
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

	// 获取 Hadoop 实例
	var originHadoop = &hadoopv1.Hadoop{}
	if err := r.Get(ctx, req.NamespacedName, originHadoop); err != nil {
		logger.Error(err, "unable to fetch hadoop")
		// 停止协调
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	hadoop := originHadoop.DeepCopy()
	hadoop.Status.Phase = hadoopv1.Reconciling

	// 检查 DeletionTimestamp 以确定对象是否正在删除
	if originHadoop.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !tools.ContainsString(originHadoop.GetFinalizers(), Finalizer) {
			controllerutil.AddFinalizer(originHadoop, Finalizer)
			if err := r.Update(ctx, originHadoop); err != nil {
				logger.Error(err, "unable to update hadoop")
				return ctrl.Result{}, err
			}
		}

		// update hadoop nn&dn daemonSet
		err := r.ensureDaemonSet(ctx, hadoop)
		if err != nil {
			logger.Error(err, "failed reconciling hadoop Master StatefulSet")
			return ctrl.Result{}, err
		}
		logger.Info("hadoop daemonSet is in sync")

		err = r.ensureService(ctx, originHadoop)
		if err != nil {
			return ctrl.Result{}, err
		}
		logger.Info("hadoop headless service is in sync")

	} else {
		// 正在删除对象
		if tools.ContainsString(originHadoop.GetFinalizers(), Finalizer) {
			// our finalizer is present, so lets handle any external dependency
			//if err := r.deleteExternalResources(ctx, originHadoop); err != nil {
			//	// if fail to delete the external dependency here, return with error
			//	// so that it can be retried
			//	return ctrl.Result{}, err
			//}
			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(originHadoop, Finalizer)
			if err := r.Update(ctx, originHadoop); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	if !reflect.DeepEqual(hadoop.TypeMeta, originHadoop.TypeMeta) || !reflect.DeepEqual(hadoop.ObjectMeta, originHadoop.ObjectMeta) || !reflect.DeepEqual(hadoop.Spec, originHadoop.Spec) {
		if err := r.Update(ctx, hadoop); err != nil {
			logger.Error(err, "update without status error")
			return ctrl.Result{}, err
		}
	}

	if !reflect.DeepEqual(hadoop.Status, originHadoop.Status) {
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

func (r *HadoopReconciler) ensureService(ctx context.Context, h *hadoopv1.Hadoop) error {
	// create headless service if it doesn't exist
	foundService := &corev1.Service{}
	if err := r.Get(ctx, types.NamespacedName{
		Name:      h.Name,
		Namespace: h.Namespace,
	}, foundService); err != nil {
		if errors.IsNotFound(err) {
			// Define and create a new headless service.
			//svc := r.headlessService(hb)
			//if err := r.Create(context.TODO(), svc); err != nil {
			//	klog.Error(err, "failed creating service")
			//	return false, err
			//}
			klog.Info("created Hadoop Service")
			return nil
		}
		klog.Error(err, "failed getting service")
		return err
	}
	return nil
}

/*
	@title    ensureDaemonSet
	@description	确保守护进程集正常
	@param: context.Context, *hadoopv1.Hadoop
	@return: bool, error
**/
func (r *HadoopReconciler) ensureDaemonSet(ctx context.Context, h *hadoopv1.Hadoop) error {
	// create if it doesn't exist
	ds := &appsv1.DaemonSet{}
	if err := r.Get(ctx, types.NamespacedName{
		Name:      h.Name,
		Namespace: h.Namespace,
	}, ds); err != nil {
		if errors.IsNotFound(err) {
			ds := r.handleDs(h)

			if err := r.Create(ctx, ds); err != nil {
				klog.Error(err, "failed creating DaemonSet")
				return err
			}
			klog.Info("Created Hadoop DaemonSet Success")
			return nil
		}
		klog.Error(err, "failed getting DaemonSet")
		return err
	}
	return nil
}

/*
	@title    handleDs
	@description	处理ds
	@param: []string
	@return: err error
**/
func (r *HadoopReconciler) handleDs(h *hadoopv1.Hadoop) *appsv1.DaemonSet {
	hadoopDs := (&h.Spec.Hdfs.DaemonSet).DeepCopy()
	dss := appsv1.DaemonSetSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"app": h.Name},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Name:        h.Name,
				Labels:      map[string]string{"app": h.Name},
				Annotations: hadoopDs.Template.Annotations,
			},
			Spec: hadoopDs.Template.Spec,
		},
	}

	ds := appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        h.Name,
			Namespace:   h.Namespace,
			Labels:      h.Labels,
			Annotations: h.Annotations,
		},
		Spec: dss,
	}
	err := controllerutil.SetControllerReference(h, &ds, r.Scheme)
	if err != nil {
		klog.Errorf("nested spec error: %v", err)
		return nil
	}
	return &ds
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

func objectEqual(originHadoop, newer *unstructured.Unstructured) (bool, error) {
	newSpec, found, err := unstructured.NestedFieldNoCopy(newer.Object, "spec")
	if err != nil {
		klog.Errorf("nested spec error: %v", err)
		return false, err
	}
	if !found {
		return false, nil
	}

	originSpecJSON, ok := originHadoop.GetAnnotations()["operator.dameng.com/spec"]
	if !ok {
		return false, nil
	}

	klog.V(6).Infof("originHadoop spec json: %v", originSpecJSON)

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
