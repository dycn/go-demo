package process

import (
	"fmt"
	"github.com/Shopify/sarama"
	"go.uber.org/zap/zapcore"
	"time"
	"zdy_demo/module/log_record"
)

func OnProcessLog(client sarama.SyncProducer, topic string) {
	now := time.Now().Local().Unix()
	key := now % 10
	switch {
	case key > 8:
		produceErrorLog(client, topic)
	case key > 7:
		produceWarnLog(client, topic)
	case key > 5:
		produceInfoLog(client, topic)
	case key >= 0:
		fallthrough
	default:
		produceDebugLog(client, topic)
	}
}

func produceDebugLog(client sarama.SyncProducer, topic string) {
	if err := log_record.OnPushLogMsg(client, topic, zapcore.DebugLevel); err != nil {
		fmt.Println(err.Error())
	}
}

func produceInfoLog(client sarama.SyncProducer, topic string) {
	if err := log_record.OnPushLogMsg(client, topic, zapcore.InfoLevel); err != nil {
		fmt.Println(err.Error())
	}
}

func produceWarnLog(client sarama.SyncProducer, topic string) {
	fmt.Println("准备生成一条警告日志")
	if err := log_record.OnPushLogMsg(client, topic, zapcore.WarnLevel); err != nil {
		fmt.Println(err.Error())
	}
}

func produceErrorLog(client sarama.SyncProducer, topic string) {
	fmt.Println("准备生成一条错误日志！")
	if err := log_record.OnPushLogMsg(client, topic, zapcore.ErrorLevel); err != nil {
		fmt.Println(err.Error())
	}
}
