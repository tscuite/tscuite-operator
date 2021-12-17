package grpc

import (
	"context"
	"fmt"
	"net"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	proto "gitee.com/tscuite/tscuite-operator/operator-proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"sigs.k8s.io/controller-runtime/pkg/log"
	//+kubebuilder:scaffold:imports
)

func RunGrpc() {
	list, err := net.Listen("tcp", ":8082")
	log.Log.Info("grpc:8082")
	if err != nil {
		log.Log.Info("grpc err=%s", err)
	}

	s := grpc.NewServer()
	proto.RegisterSearchServiceServer(s, &SearchService{})
	reflection.Register(s)
	if err := s.Serve(list); err != nil {
		log.Log.Info("failed to serve: %v", err)
	}
}

type SearchService struct{}

func (s *SearchService) Search(ctx context.Context, r *proto.SearchRequest) (*proto.SearchResponse, error) {
	msg := r.Request
	Grpcagent, err := s.Grpc(string(msg))
	if err != nil {
		fmt.Println(err)
	}
	log.Log.Info("grpc转发成功: " + Grpcagent)
	return &proto.SearchResponse{Response: r.GetRequest() + "服务端给你返回的消息"}, nil
}

const PORT string = "127.0.0.1:8085"

func (s *SearchService) Grpc(request string) (string, error) {
	conn, err := grpc.Dial(PORT, grpc.WithInsecure())
	if err != nil {
		log.Log.Info("grpc连接失败! %v", err)
	}
	defer conn.Close()
	client := proto.NewSearchServiceClient(conn)
	req, err := client.Search(context.Background(), &proto.SearchRequest{
		Request: request,
	})
	if err != nil {
		log.Log.Info("grpc发送消息失败! %v", err)
	}
	return req.GetResponse(), err
}
