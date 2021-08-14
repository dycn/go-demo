package process

import (
	"fmt"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"time"
	pb "zdy_demo/common/rpc/protobuf"
)

type StreamService struct {
	pb.UnimplementedStreamServiceServer
}

func (s *StreamService) Eat(r *pb.Request, stream pb.StreamService_EatServer) error {
	//设置header信息 与sendHeader不可同时用，否则SendHeader会覆盖前一个
	// SetHeader sets the header metadata. It may be called multiple times.
	// When call multiple times, all the provided metadata will be merged.
	// All the metadata will be sent out when one of the following happens:
	//  - ServerStream.SendHeader() is called;
	//  - The first response is sent out;
	//  - An RPC status is sent out (error or success).
	if err := stream.SetHeader(metadata.MD{"debug": []string{"this is server resp header"}}); nil != err {
		return err
	}
	//设置header信息
	// SendHeader sends the header metadata.
	// The provided md and headers set by SetHeader() will be sent.
	// It fails if called multiple times.
	//if err := stream.SendHeader(metadata.MD{"cc":[]string{"dd"}});err != nil{
	//	return err
	//}

	//设置metadata
	stream.SetTrailer(metadata.MD{"info": []string{"this is server metadata resp trailer"}})

	//metadata
	data, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		log.Fatalln("metadata from incoming err")
	}
	fmt.Println("收到的ctx是:", data)

	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)
		//这里发送的流有最大字节限制 默认为4M
		//可以通过Server.Builder.SetMaxReceiveMessageSize(int)设置,参数是int, 意味着可设置的最大值是2147483674(2G)
		err := stream.Send(&pb.Response{
			Result: "服务端流式返回,接收到" + r.Name,
		})

		if err != nil {
			return err
		}

	}

	return nil
}

func (s *StreamService) Work(stream pb.StreamService_WorkServer) error {
	//设置header信息 sendHeader不可同时用，否则SendHeader会覆盖前一个
	if err := stream.SetHeader(metadata.MD{"cc2": []string{"dd2"}}); nil != err {
		return err
	}
	//设置header信息
	//if err := stream.SendHeader(metadata.MD{"cc":[]string{"dd"}});err != nil{
	//	return err
	//}

	stream.SetTrailer(metadata.MD{"client": []string{"stream"}})

	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Response{
				Result: "client-stream",
			})
		}
		if err != nil {
			return err
		}
		log.Printf("client stream Recv Name: %s", r.Name)
	}
}
func (s *StreamService) Sleep(stream pb.StreamService_SleepServer) error {
	n := 0
	for n < 5 {
		err := stream.Send(&pb.Response{
			Result: "双向stream-value",
		})
		if err != nil {
			return err
		}
		r, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		n++
		log.Printf("stream.Recv Name: %s", r.Name)
		time.Sleep(time.Second)
	}
	return nil
}
