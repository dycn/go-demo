package process

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"os"
	"time"
	pb "zdy_demo/common/rpc/protobuf"
)

func service() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(":20000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	c := pb.NewGreeterClient(conn)
	c2 := pb.NewOwnerClient(conn)

	// Contact the server and print out its response.
	name := "world"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("greet服务错误: %v", err)
	}
	log.Printf("Greeting服务返回: %s", r.GetMessage())

	r2, err := c2.CheckGuest(ctx, &pb.CheckRequest{Passwd: "1password"})
	if err != nil {
		log.Fatalf("CheckGuest服务错误: %v", err)
	}
	log.Printf("CheckGuest服务返回: %s", r2.GetGuest())
}

func demo() {
	//WithBlock Dial连接是异步对,设置会导致阻塞等待握手成功  当未设置此设置项时 ctx对超时控制无任何效果
	//WithInsecure  disables transport security
	//WithDisAbleRetry  disables retries,even if the service config enables them.
	//WithChainStreamInterceptor 设置流式拦截器
	//WithUnaryInterceptor 一元拦截器
	conn, err := grpc.Dial(":50000", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithDisableRetry())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	client := pb.NewStreamServiceClient(conn)
	req := &pb.Request{
		Name: "zdy",
	}
	//http1.1可以设置header传递数据 grpc基于http2采用metadata传递数据
	//
	// 一元模式下
	// 客户端消息由 头部-内容1-EOS(end of stream) 构成，服务端消息由 头部-内容1-trailer 构成 header和trailer是同时的
	//
	// 双向流RPC模式通过多发多收模式
	// 客户端消息由 头部-内容3-内容1-内容2-EOS 构成，
	// 服务端消息由 头部-内容1-内容2-内容3-trailer   trailer是在流最后才能收到
	//
	//metadata实质上是map[string][]string 类型 key值不区分大小写
	var ctx context.Context
	md := metadata.Pairs(
		"key1", "val1",
		"Key1", "val1-2", // "key1" will have map value []string{"val1", "val1-2"}
		"key2", "val2",
	)
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	stream, err := client.Eat(ctx, req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	//接收server的header信息
	fmt.Println(stream.Header()) //map[cc:[dd] content-type:[application/grpc]]
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("eof")
			break
		}
		if err != nil {
			log.Fatalln(err.Error())
		}
		log.Printf("resp: Result %s", resp.Result)
	}
	//在一元rpc中header和trailer是一起到达的，在流式中是在流结束后到达的
	fmt.Println(stream.Trailer()) //map[cc1:[dd1]]

	return
	fmt.Println("客户端流式 start")
	stream2, err := client.Work(context.Background())
	if err != nil {
		log.Fatalln(err.Error())
	}
	for i := 0; i < 5; i++ {
		//fmt.Println(r)
		err := stream2.Send(req)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err.Error())
		}
		time.Sleep(time.Second)
	}
	//注意这个header是设置不了的
	//fmt.Println(stream2.Header())
	resp, err := stream2.CloseAndRecv()
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("resp: result %s", resp.Result)
	fmt.Println(stream.Trailer()) //map[cc1:[dd1]]

	fmt.Println("双端流式")
	stream3, err := client.Sleep(context.Background())
	if err != nil {
		log.Fatalln(err.Error())
	}
	for n := 0; n < 5; n++ {
		err = stream3.Send(req)
		if err != nil {
			log.Fatalln(err.Error())
		}
		resp, err := stream3.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err.Error())
		}
		log.Printf("resp: result: %s", resp.Result)
	}
	if err = stream3.CloseSend(); nil != err {
		log.Println(err)
	}
}
