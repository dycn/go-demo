package process

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
	pb "zdy_demo/common/rpc/protobuf"
)

var guests []*pb.Guest

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Printf("Received: %v", in.GetName())

	guests = append(guests, &pb.Guest{
		Name:      in.GetName(),
		Time:      timestamppb.Now(),
		Timestamp: time.Now().Unix(),
	})

	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

type owner struct {
	pb.UnimplementedOwnerServer
}

func (o *owner) CheckGuest(_ context.Context, in *pb.CheckRequest) (*pb.CheckResponse, error) {
	if in.Passwd == "password" {
		return &pb.CheckResponse{Guest: guests}, nil
	}
	return nil, errors.New("password wrong")
}
