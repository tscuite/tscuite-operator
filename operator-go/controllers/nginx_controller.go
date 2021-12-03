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

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	tscuitev1 "github.com/tscuite/crd/operator-go/api/v1"
)

// NginxReconciler reconciles a Nginx object
type NginxReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tscuite.registry.cn-hangzhou.aliyuncs.com,resources=nginxes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tscuite.registry.cn-hangzhou.aliyuncs.com,resources=nginxes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tscuite.registry.cn-hangzhou.aliyuncs.com,resources=nginxes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Nginx object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *NginxReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	NginxV1 := &tscuitev1.Nginx{}
	if err := r.Get(ctx, req.NamespacedName, NginxV1); err != nil {
		log.Log.Error(err, "ns", req.NamespacedName)
	} else {
		log.Log.Info("Info", "ns", req.NamespacedName)
		if err := r.NginxOperatorGet(ctx, req, r.NginxDeployment(NginxV1)); err != nil {
			return ctrl.Result{}, r.NginxOperatorCreate(r.Client, req, r.NginxDeployment(NginxV1))
		} else {
			return ctrl.Result{}, r.NginxOperatorUpdate(r.Client, req, r.NginxDeployment(NginxV1))
		}
	}
	return ctrl.Result{}, nil
}

func (r *NginxReconciler) NginxOperatorGet(ctx context.Context, req ctrl.Request, nginxdeployment *appsv1.Deployment) error {
	return r.Client.Get(ctx, req.NamespacedName, nginxdeployment)
}
func (r *NginxReconciler) NginxOperatorCreate(client client.Client, req ctrl.Request, nginxdeployment *appsv1.Deployment) error {
	log.Log.Info("Create", "ns", req.NamespacedName)
	return client.Create(context.Background(), nginxdeployment)
}
func (r *NginxReconciler) NginxOperatorUpdate(client client.Client, req ctrl.Request, nginxdeployment *appsv1.Deployment) error {
	log.Log.Info("Update", "ns", req.NamespacedName)
	return client.Update(context.Background(), nginxdeployment)
}
func (r *NginxReconciler) NginxDeployment(nginx *tscuitev1.Nginx) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nginx.Name,
			Namespace: nginx.Namespace,
			Labels: map[string]string{
				"app": nginx.Name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr2(nginx.Spec.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": nginx.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": nginx.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: nginx.Name,
							Image: nginx.Spec.Images,
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("10m"),
									corev1.ResourceMemory: resource.MustParse("32Mi"),
								},
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("10m"),
									corev1.ResourceMemory: resource.MustParse("32Mi"),
								},
							},
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: nginx.Spec.Port,
									Name:          nginx.Name,
									Protocol:      "TCP",
								},
							},
							LivenessProbe: &corev1.Probe{
								FailureThreshold:    5,
								InitialDelaySeconds: 60,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								TimeoutSeconds:      5,
								Handler: corev1.Handler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.FromInt(int(nginx.Spec.Port)),
									},
								},
							},
							ReadinessProbe: &corev1.Probe{
								FailureThreshold: 3,
								PeriodSeconds:    10,
								SuccessThreshold: 1,
								TimeoutSeconds:   1,
								Handler: corev1.Handler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.FromInt(int(nginx.Spec.Port)),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
func int32Ptr2(i int32) *int32 { return &i }

// SetupWithManager sets up the controller with the Manager.
func (r *NginxReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tscuitev1.Nginx{}).
		Complete(r)
}
