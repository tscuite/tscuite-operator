/*
Copyright 2017 Heptio Inc.

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

package sinks

import (
	"context"
	"encoding/json"

	proto "gitee.com/tscuite/tscuite-operator/operator-proto"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	v1 "k8s.io/api/core/v1"
)

// GlogSink is the most basic sink
// Useful when you already have ELK/EFK Stack
type GlogSink struct {
	// TODO: create a channel and buffer for scaling
}

// NewGlogSink will create a new
func NewGlogSink() EventSinkInterface {
	return &GlogSink{}
}

// UpdateEvents implements the EventSinkInterface
func (gs *GlogSink) UpdateEvents(eNew *v1.Event, eOld *v1.Event) {
	eData := NewEventData(eNew, eOld)

	if eJSONBytes, err := json.Marshal(eData); err == nil {
		Grpcagent, err := gs.Grpc(string(eJSONBytes))
		if err != nil {
			glog.Warningf("grpc连接失败!: %v", err)
		}
		glog.Info(Grpcagent)
		//glog.Info(string(eJSONBytes))

	} else {
		glog.Warningf("Failed to json serialize event: %v", err)
	}
}

const PORT = "192.168.6.229:8082"

func (gs *GlogSink) Grpc(request string) (string, error) {
	conn, err := grpc.Dial(PORT, grpc.WithInsecure())
	if err != nil {

		glog.Warningf("grpc连接失败! %v", err)
	}
	defer conn.Close()
	client := proto.NewSearchServiceClient(conn)
	req, err := client.Search(context.Background(), &proto.SearchRequest{
		Request: request,
	})
	if err != nil {

		glog.Warningf("grpc发送消息失败!: %v", err)
	}
	return req.GetResponse(), err
}
