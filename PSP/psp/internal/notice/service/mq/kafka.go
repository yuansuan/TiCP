package mq

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
	"github.com/yuansuan/ticp/common/go-kit/logging"
)

// KafkaProducer 消息生产者
type KafkaProducer struct {
	writer *kafka.Writer
}

// KafkaConsumer 消息消费者
type KafkaConsumer struct {
	reader *kafka.Reader
}

// TestMsg ...
type TestMsg struct {
	Key string `json:"key"`
	Msg any    `json:"msg"`
}

// NewKafkaProducer 创建消息生产者
func NewKafkaProducer(topic string) (*KafkaProducer, error) {

	writer := middleware.NewKafkaWriterWithConfig(kafka.WriterConfig{
		Topic:     topic,
		BatchSize: 1,
		Balancer:  &kafka.LeastBytes{},
	})
	if writer == nil {
		return nil, errors.New("failed to create kafka writer")
	}

	return &KafkaProducer{
		writer: writer,
	}, nil
}

// NewKafkaConsumer 创建消费者
func NewKafkaConsumer(topic string, group string) (*KafkaConsumer, error) {

	reader := middleware.NewKafkaReader(topic, group)
	if reader == nil {
		return nil, errors.New("failed to create kafka reader")
	}

	return &KafkaConsumer{
		reader: reader,
	}, nil
}

// SendMessage 发送消息
func (s *KafkaProducer) SendMessage(ctx context.Context, key string, value any) error {
	logger := logging.GetLogger(ctx)

	message, err := json.Marshal(value)
	if err != nil {
		return err
	}

	kafkaMessage := kafka.Message{
		Key:   []byte(key),
		Value: message,
	}
	if err = s.writer.WriteMessages(ctx, kafkaMessage); err != nil {
		return errors.Wrap(err, "send kafka message err")
	}

	logger.Infof("kafka send message: key=[%v], value=[%v]", key, string(message))

	return nil
}

// ReadByteMessage 读取字节消息
func (s *KafkaConsumer) ReadByteMessage(ctx context.Context) ([]byte, error) {
	logger := logging.GetLogger(ctx)

	msg, err := s.reader.ReadMessage(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "read kafka message err")
	}

	logger.Infof("read kafka message: key=[%v], value=[%v]", string(msg.Key), string(msg.Value))

	return msg.Value, nil
}

// ReadMessage 读取消息
func (s *KafkaConsumer) ReadMessage(ctx context.Context, v any) error {
	logger := logging.GetLogger(ctx)

	msg, err := s.reader.ReadMessage(ctx)
	if err != nil {
		return errors.Wrap(err, "read kafka message err")
	}

	err = json.Unmarshal(msg.Value, v)
	if err != nil {
		return errors.Wrap(err, "unmarshal kafka message err")
	}

	logger.Infof("kafka read message: key=[%v], value=[%v]", string(msg.Key), string(msg.Value))

	return nil
}

// CloseMQService 关闭消息服务
//func (s *KafkaService) CloseMQService() {
//	s.writer.Close()
//	s.reader.Close()
//}
