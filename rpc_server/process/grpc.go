package process

import (
	"google.golang.org/grpc"
	"log"
	"net"
	pb "zdy_demo/common/rpc/protobuf"
)

func service() {
	lis, err := net.Listen("tcp", ":20000")
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	pb.RegisterOwnerServer(s, &owner{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("服务开启失败: %v", err)
	}
}

func demo() {
	server := grpc.NewServer()
	pb.RegisterStreamServiceServer(server, &StreamService{})
	lis, err := net.Listen("tcp", ":50000")
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	if err := server.Serve(lis); err != nil {
		log.Fatalf(err.Error())
	}
}
