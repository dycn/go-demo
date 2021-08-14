package main

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"os"
	"zdy_demo/module/log_record"
)

//此时又新增一个消费者, 原来对group1组对分配策略是sarama.BalanceStrategySticky 会重新平衡各个消费者消费对分片
func main() {
	config := log_record.NewConfig()
	client, err := log_record.NewClient([]string{"127.0.0.1:9092"}, config)
	groupAc3, err := sarama.NewConsumerGroupFromClient("group1", client)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer func() { _ = groupAc3.Close() }()

	ctx := context.Background()
	err = groupAc3.Consume(ctx, []string{"log_record"}, consumerGroupAClient3Handler{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

type consumerGroupAClient3Handler struct{}

func (c consumerGroupAClient3Handler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}
func (c consumerGroupAClient3Handler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 循环消费消息 claim.Message() chan
func (c consumerGroupAClient3Handler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	fmt.Printf("上次消费到 offset:%d \n", claim.InitialOffset()) //打印上次消费的位置

	for msg := range claim.Messages() {
		log_record.OnHandlerGroup1Client3SlowRecvMsg(msg, nil)
		sess.MarkMessage(msg, "") //标记已消费
	}
	return nil
}
