package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

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
	//log.Log.Info("收到了一条来自客户端的消息: " + msg)
	var data AutoGenerated
	Grpcagent, err := s.Grpc(string(msg))
	if err != nil {
		log.Log.Info("grpc连接失败! %v", err)
	}
	log.Log.Info("grpc连接ok! %v", Grpcagent)
	if err := json.Unmarshal([]byte(msg), &data); err == nil {
		log.Log.Info("收到了一条来自客户端的pod事件：" + data.Event.InvolvedObject.Name)
		log.Log.Info("收到了一条来自客户端的Message: " + data.Event.Message)
	} else {
		fmt.Println(err)
	}
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

type AutoGenerated struct {
	Verb  string `json:"verb"`
	Event struct {
		Metadata struct {
			Name              string    `json:"name"`
			Namespace         string    `json:"namespace"`
			UID               string    `json:"uid"`
			ResourceVersion   string    `json:"resourceVersion"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
			ManagedFields     []struct {
				Manager    string    `json:"manager"`
				Operation  string    `json:"operation"`
				APIVersion string    `json:"apiVersion"`
				Time       time.Time `json:"time"`
			} `json:"managedFields"`
		} `json:"metadata"`
		InvolvedObject struct {
			Kind            string `json:"kind"`
			Namespace       string `json:"namespace"`
			Name            string `json:"name"`
			UID             string `json:"uid"`
			APIVersion      string `json:"apiVersion"`
			ResourceVersion string `json:"resourceVersion"`
		} `json:"involvedObject"`
		Reason  string `json:"reason"`
		Message string `json:"message"`
		Source  struct {
			Component string `json:"component"`
		} `json:"source"`
		FirstTimestamp     time.Time   `json:"firstTimestamp"`
		LastTimestamp      time.Time   `json:"lastTimestamp"`
		Count              int         `json:"count"`
		Type               string      `json:"type"`
		EventTime          interface{} `json:"eventTime"`
		ReportingComponent string      `json:"reportingComponent"`
		ReportingInstance  string      `json:"reportingInstance"`
	} `json:"event"`
}
