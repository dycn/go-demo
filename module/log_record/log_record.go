package log_record

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"go.uber.org/zap/zapcore"
	"log"
	"math/rand"
	"time"
	"zdy_demo/utility"
)

type LogRecord struct {
	Uid      int           `json:"uid"`
	Content  string        `json:"content"`
	LogLevel zapcore.Level `json:"log_level"`
}

func NewConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion

	//配置作用于 sarama.Config.Producer struct
	config.Producer.Partitioner = sarama.NewRandomPartitioner //随机分区
	config.Producer.Return.Successes = true                   //消息发送成功后通知 success channel

	//配置作用于 sarama.Config.Consumer struct
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = true //是否自动提交偏移量给broker
	//Range的分配策略是 将分片数与消费者数量取模,然后分配  消费多个topic时 第一个消费者总是多分配到分片
	//RoundRobin的策略是 将消费者所有的topic、分片轮询分配给各个消费者 组内的消费者订阅的topic不同时 容易导致订阅多的消费者消费压力大
	//Sticky的策略是 分配尽量平均，尽量与上次分配相同
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	return config
}

func NewClient(addrs []string, config *sarama.Config) (sarama.Client, error) {
	consumer, err := sarama.NewClient(addrs, config)
	if err != nil {
		return nil, err
	}
	return consumer, err
}

func NewGroup(addrs []string, groupId string, config *sarama.Config) (sarama.ConsumerGroup, error) {
	group, err := sarama.NewConsumerGroup(addrs, groupId, config)
	if err != nil {
		return nil, err
	}
	return group, err
}

func OnPushLogMsg(client sarama.SyncProducer, topic string, logLevel zapcore.Level) error {
	logRecord := newLogRecord(logLevel)
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	logRecordJson, err := json.Marshal(logRecord)
	if err != nil {
		return err
	}
	msg.Value = sarama.StringEncoder(logRecordJson)

	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		return err
	}

	fmt.Printf("发送一条消息成功 topic:%s partition:%d offset:%v\n", topic, pid, offset)
	return nil
}

func newLogRecord(logLevel zapcore.Level) *LogRecord {
	logRecord := new(LogRecord)
	rand.Seed(time.Now().Local().Unix())
	logRecord.Uid = rand.Int()
	logRecord.LogLevel = logLevel
	logRecord.Content = getLogContent(logLevel)
	return logRecord
}

func getLogContent(logLevel zapcore.Level) string {
	switch logLevel {
	case zapcore.InfoLevel:
		return fmt.Sprintf("[%s] 这是一条info级别日志", utility.GetNowTimeStr())
	case zapcore.WarnLevel:
		return fmt.Sprintf("[%s] 这是一条warnning级别日志", utility.GetNowTimeStr())
	case zapcore.ErrorLevel:
		return fmt.Sprintf("[%s] 这是一条error级别日志", utility.GetNowTimeStr())
	case zapcore.DebugLevel:
		fallthrough
	default:
		return fmt.Sprintf("[%s] 这是一条debug级别日志", utility.GetNowTimeStr())
	}
}

func OnHandlerGroup1Client1RecvMsg(msg *sarama.ConsumerMessage, logger *log.Logger) {
	logger.Println(time.Now().UnixNano())
	fmt.Printf("[group1.client1] 消费 topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)

	logRecord := new(LogRecord)
	err := json.Unmarshal(msg.Value, logRecord)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//switch logRecord.LogLevel {
	//case zapcore.ErrorLevel:
	//	fmt.Printf("发生了错误日志! 日志内容为:\n%s\n", logRecord.Content)
	//case zapcore.WarnLevel:
	//	fmt.Printf("发生了警告日志, 请检查:\n%s\n", logRecord.Content)
	//default:
	//
	//}
}

func OnHandlerGroup1Client2RecvMsg(msg *sarama.ConsumerMessage, logger *log.Logger) {
	logger.Println(time.Now().UnixNano())
	fmt.Printf("[group1.client2] 消费 topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
	logRecord := new(LogRecord)
	err := json.Unmarshal(msg.Value, logRecord)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func OnHandlerGroup1Client3SlowRecvMsg(msg *sarama.ConsumerMessage, logger *log.Logger) {
	if logger != nil {
		logger.Println(time.Now().UnixNano())
	}
	fmt.Printf("[group1.client3] 消费 topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
	logRecord := new(LogRecord)
	err := json.Unmarshal(msg.Value, logRecord)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func OnHandlerGroup2RecvMsg(msg *sarama.ConsumerMessage, logger *log.Logger) {
	logger.Println(time.Now().UnixNano())
	//fmt.Printf("[group2] 消费 topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
	logRecord := new(LogRecord)
	err := json.Unmarshal(msg.Value, logRecord)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//switch logRecord.LogLevel {
	//case zapcore.ErrorLevel:
	//	fmt.Printf("发生了错误日志! 日志内容为:\n%s\n", logRecord.Content)
	//case zapcore.WarnLevel:
	//	fmt.Printf("发生了警告日志, 请检查:\n%s\n", logRecord.C	ontent)
	//default:
	//
	//}
}
