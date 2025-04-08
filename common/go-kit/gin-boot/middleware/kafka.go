package middleware

import (
	"context"
	"os"
	"strings"

	"github.com/segmentio/kafka-go"
	"github.com/yuansuan/ticp/common/go-kit/logging"
)

// kafkaT kafkaT
type kafkaT struct {
	brokers []string
}

func (mw *Middleware) initKafka() {
	conf := &mw.conf.App.Middleware.Kafka
	if !conf.KafkaClusterStartUp {
		return
	}
	// if env YS_KAFKA_BROKER_URL is set, use $(YS_KAFKA_BROKER_URL) as kafka broker address
	// if not, use from config file
	kafkaenv := os.Getenv("YS_KAFKA_BROKER_URL")
	logging.GetLogger(context.Background()).Infof("YS_KAFKA_BROKER_URL is %v\n", kafkaenv)
	if kafkaenv != "" {
		mw.kafka = newKafka(strings.Split(kafkaenv, ","))
	} else {
		mw.kafka = newKafka(conf.KafkaClusterURL)
	}
	for k, v := range mw.kafka.brokers {
		logging.GetLogger(context.Background()).Infof("kafka broker %v is %v", k, v)
	}
}

// newKafka newKafka
func newKafka(brokers []string) *kafkaT {
	return &kafkaT{
		brokers: brokers,
	}
}

// NewKafkaReader NewKafkaReader
func NewKafkaReader(topic string, group string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  Instance.kafka.brokers,
		Topic:    topic,
		GroupID:  group,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
}

// NewKafkaWriter NewKafkaWriter
func NewKafkaWriter(topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  Instance.kafka.brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}

// NewKafkaWriterWithConfig NewKafkaWriterWithConfig
func NewKafkaWriterWithConfig(config kafka.WriterConfig) *kafka.Writer {
	config.Brokers = Instance.kafka.brokers
	return kafka.NewWriter(config)
}
