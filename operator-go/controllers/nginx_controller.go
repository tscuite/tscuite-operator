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
	"fmt"

	tscuitev1 "github.com/tscuite/crd/operator-go/api/v1"
	kubernetesv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
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
		err = CreatePod(r.Client, nginx)
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func CreatePod(client client.Client, nginx *tscuitev1.Nginx) error {
	newPod := &kubernetesv1.Pod{}
	newPod.Name = nginx.Name
	newPod.Namespace = nginx.Namespace
	newPod.Spec.Containers = []kubernetesv1.Container{
		{
			Name:            nginx.Name,
			Image:           nginx.Spec.Images,
			ImagePullPolicy: kubernetesv1.PullIfNotPresent,
			Ports: []kubernetesv1.ContainerPort{
				{
					ContainerPort: nginx.Spec.Prot,
				},
			},
		},
	}
	return client.Create(context.Background(), newPod)
}

// SetupWithManager sets up the controller with the Manager.
func (r *NginxReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tscuitev1.Nginx{}).
		Complete(r)
}
