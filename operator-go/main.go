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

package main

import (
	"context"
	"flag"
	"net"
	"os"
	"strconv"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	pb "gitee.com/tscuite/tscuite-operator/operator-go/grpc"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	tscuitev1 "gitee.com/tscuite/tscuite-operator/operator-go/api/v1"
	"gitee.com/tscuite/tscuite-operator/operator-go/controllers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"sigs.k8s.io/controller-runtime/pkg/log"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(tscuitev1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	//end := make(chan bool, 1)
	go RunGrpc()
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "f6d5007b.registry.cn-hangzhou.aliyuncs.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.NginxReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Nginx")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
		//<-end
	}
}

func RunGrpc() {
	list, err := net.Listen("tcp", ":8082")
	log.Log.Info("grpc:8082")
	if err != nil {
		log.Log.Info("grpc err=%s", err)
	}

	s := grpc.NewServer()
	pb.RegisterSearchServiceServer(s, &SearchService{})
	reflection.Register(s)
	if err := s.Serve(list); err != nil {
		log.Log.Info("failed to serve: %v", err)
	}
}

type SearchService struct{}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	log.Log.Info("收到了一条来自客户端的消息: " + r.Request)
	log.Log.Info("收到了一条来自客户端的消息: " + strconv.FormatInt(int64(r.XXX_sizecache), 10))

	return &pb.SearchResponse{Response: r.GetRequest() + " HTTP 服务端给你返回的消息"}, nil
}
func (s *SearchService) Saarch(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	log.Log.Info("收到了一条来自客户端的消息2: " + r.Request)
	log.Log.Info("收到了一条来自客户端的消息2: " + strconv.FormatInt(int64(r.XXX_sizecache), 10))

	return &pb.SearchResponse{Response: r.GetRequest() + " HTTP 2服务端给你返回的消息"}, nil
}
