package process

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"zdy_demo/module/log_record"
)

type consumerGroupAClient1Handler struct{}

func (c consumerGroupAClient1Handler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}
func (c consumerGroupAClient1Handler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 循环消费消息 claim.Message() chan
func (c consumerGroupAClient1Handler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	fmt.Printf("client1上次消费到 offset:%d \n", claim.InitialOffset()) //打印上次消费的位置
	logger := getLogger("groupAClient1")

	for msg := range claim.Messages() {
		log_record.OnHandlerGroup1Client1RecvMsg(msg, logger)
		sess.MarkMessage(msg, "") //标记已消费
	}
	return nil
}

type consumerGroupAClient2Handler struct{}

func (c consumerGroupAClient2Handler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}
func (c consumerGroupAClient2Handler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 循环消费消息 claim.Message() chan
func (c consumerGroupAClient2Handler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	fmt.Printf("client2上次消费到 offset:%d \n", claim.InitialOffset()) //打印上次消费的位置
	logger := getLogger("groupAClient2")

	for msg := range claim.Messages() {
		log_record.OnHandlerGroup1Client2RecvMsg(msg, logger)
		sess.MarkMessage(msg, "") //标记已消费
	}
	return nil
}

type consumerGroupAClient3SlowHandler struct{}

func (c consumerGroupAClient3SlowHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}
func (c consumerGroupAClient3SlowHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 循环消费消息 claim.Message() chan
func (c consumerGroupAClient3SlowHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	fmt.Printf("client3上次消费到 offset:%d \n", claim.InitialOffset()) //打印上次消费的位置
	logger := getLogger("groupAClient3")

	for msg := range claim.Messages() {
		log_record.OnHandlerGroup1Client3SlowRecvMsg(msg, logger)
		sess.MarkMessage(msg, "") //标记已消费
	}
	return nil
}

type consumerGroupBHandler struct{}

func (c consumerGroupBHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}
func (c consumerGroupBHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 循环消费消息 claim.Message() chan
func (c consumerGroupBHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	fmt.Printf("clientB上次消费到 offset:%d \n", claim.InitialOffset()) //打印上次消费的位置

	//groupB写到日志
	logger := getLogger("groupB")

	for msg := range claim.Messages() {
		log_record.OnHandlerGroup2RecvMsg(msg, logger)
		sess.MarkMessage(msg, "") //标记已消费
	}
	return nil
}

func getLogger(filename string) *log.Logger {
	dir, _ := os.Getwd()
	file := dir + "/" + filename + ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	logger := log.New(logFile, "", log.Lshortfile) // 将文件设置为loger作为输出
	return logger
}
