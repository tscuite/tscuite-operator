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
	"encoding/json"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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

	// your logic here
	nginx := &tscuitev1.Nginx{}
	if err := r.Get(ctx, req.NamespacedName, nginx); err != nil {
		fmt.Println(err)
	} else {
		err = Deployment(r.Client, nginx)
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
func Operator(client client.Client, deployment *appsv1.Deployment) error {
	//新建
	//err := client.Create(context.Background(), deployment)
	//更新
	err := client.Update(context.Background(), deployment)
	return err
}
func Deployment(client client.Client, nginx *tscuitev1.Nginx) error {
	var memory corev1.ResourceRequirements
	var replicas int32 = nginx.Spec.Replicas
	data := `{"limits": {"cpu":"2000m", "memory": "1Gi"}, "requests": {"cpu":"2000m", "memory": "1Gi"}}`
	json.Unmarshal([]byte(data), &memory)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nginx.Name,
			Namespace: nginx.Namespace,
			Labels: map[string]string{
				"app": nginx.Name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
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
							Image:           nginx.Spec.Images,
							Resources:       memory,
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
									HTTPGet: &corev1.HTTPGetAction{
										Path:   "/ready",
										Scheme: "HTTP",
										Port:   intstr.FromInt(int(nginx.Spec.Port)),
									},
								},
							},
							ReadinessProbe: &corev1.Probe{
								FailureThreshold: 3,
								PeriodSeconds:    10,
								SuccessThreshold: 1,
								TimeoutSeconds:   1,
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path:   "/ready",
										Scheme: "HTTP",
										Port:   intstr.FromInt(int(nginx.Spec.Port)),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return Operator(client, deployment)
}

// SetupWithManager sets up the controller with the Manager.
func (r *NginxReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tscuitev1.Nginx{}).
		Complete(r)
}
