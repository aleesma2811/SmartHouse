package kafkax

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

func ensureTopicOnce(broker, topic string, partitions int) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	ctrl, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer ctrl.Close()

	err = ctrl.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: 1,
	})

	if err != nil && !errors.Is(err, kafka.TopicAlreadyExists) {
		return err
	}
	return nil
}

func EnsureTopic(broker, topic string, partitions int) error {
	var err error
	for i := 1; i <= 15; i++ {
		if err = ensureTopicOnce(broker, topic, partitions); err == nil {
			return nil
		}
		log.Printf("kafka: tópico %s, intento %d/15: %v", topic, i, err)
		time.Sleep(time.Duration(i) * time.Second)
	}
	return err
}

type Publisher struct {
	w *kafka.Writer
}

func NewPublisher(brokers []string, topic string) *Publisher {
	return &Publisher{w: &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
		BatchTimeout: 10 * time.Millisecond,
	}}
}

func (p *Publisher) Publish(ctx context.Context, key, eventType string, payload any) error {
	value, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return p.w.WriteMessages(ctx, kafka.Message{
		Key:     []byte(key),
		Value:   value,
		Headers: []kafka.Header{{Key: "event-type", Value: []byte(eventType)}},
	})
}

func (p *Publisher) Close() error { return p.w.Close() }

func EventType(msg kafka.Message) string {
	for _, h := range msg.Headers {
		if h.Key == "event-type" {
			return string(h.Value)
		}
	}
	return ""
}

type Handler func(ctx context.Context, msg kafka.Message) error

func Consume(ctx context.Context, brokers []string, topic, groupID string, handle Handler) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.FirstOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
	})
	defer r.Close()

	for {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}

		if err := processWithRetry(ctx, handle, msg, 3); err != nil {

			log.Printf("kafka: se abandona partition=%d offset=%d: %v", msg.Partition, msg.Offset, err)
		}

		if err := r.CommitMessages(ctx, msg); err != nil {
			return err
		}
	}
}

func processWithRetry(ctx context.Context, h Handler, msg kafka.Message, attempts int) error {
	var err error
	for i := 1; i <= attempts; i++ {
		if err = h(ctx, msg); err == nil {
			return nil
		}
		log.Printf("kafka: intento %d/%d falló: %v", i, attempts, err)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(i) * 500 * time.Millisecond):
		}
	}
	return err
}
