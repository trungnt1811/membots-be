package msgqueue

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/util/log"

	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

const (
	KAFKA_NEW_DATA_MAX_WAIT          = time.Minute * 1
	KAFKA_TIME_OUT                   = time.Second * 3
	KAFKA_READ_BATCH_TIME_OUT        = time.Second * 60
	KAFKA_TOPIC_SHIPPING_BATCH       = "affiliate-system-batch"
	KAFKA_GROUP_ID                   = "affiliate-system-backend"
	KAFKA_TOPIC_PEDNING_TX           = "affiliate-system-pending-tx"
	KAFKA_TOPIC_IMPORT_RECEIPT_TX    = "affiliate-system-receipt"
	KAFKA_TOPIC_AFF_ORDER_APPROVE    = "aff-order-approved"
	KAFKA_NOTI_GROUP_ID              = "reward-notification-backend"
	KAFKA_TOPIC_NOTI_SMS             = "notification-SMS"
	KAFKA_TOPIC_NOTI_EMAIL           = "notification-EMAIL"
	KAFKA_TOPIC_NOTI_STATUS_RESPONSE = "notification-status"
	KAFKA_TOPIC_NOTI_APP_MESSAGE     = "notification-worker-app-noti-message"
)

type QueueBasic struct {
	Brokers []string
	User    string
	Sigchan chan os.Signal
}

type QueueReader struct {
	QueueBasic
	*kafka.Reader
}

type QueueWriter struct {
	QueueBasic
	*kafka.Writer
}

type MessageQueue struct {
	Reader *kafka.Reader
	Writer *kafka.Writer
	QueueBasic
}

func NewMsgQueue(topic, groupId string) *MessageQueue {
	mq := &MessageQueue{}
	config := conf.GetConfiguration()

	mq.User = config.Kafka.User
	mq.Brokers = []string{config.Kafka.KafkaURL}
	mq.Writer = newKafkaWriter(config, topic, mq.Brokers)
	mq.Reader = newKafkaReader(config, topic, groupId, mq.Brokers)

	return mq
}

func NewKafkaProducer(topic string) *QueueWriter {
	mq := &QueueWriter{}
	config := conf.GetConfiguration()

	mq.User = config.Kafka.User
	mq.Brokers = []string{config.Kafka.KafkaURL}
	mq.Writer = newKafkaWriter(config, topic, mq.Brokers)

	return mq
}

func NewKafkaConsumer(topic, groupId string) *QueueReader {
	mq := &QueueReader{}
	config := conf.GetConfiguration()

	mq.User = config.Kafka.User
	mq.Brokers = []string{config.Kafka.KafkaURL}
	mq.Reader = newKafkaReader(config, topic, groupId, mq.Brokers)

	return mq
}

func newKafkaWriter(config *conf.Configuration, topic string, brokers []string) *kafka.Writer {
	sharedTransport, err := GetTransport(config)
	if err != nil {
		panic(err)
	}
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		Transport:    sharedTransport,
		BatchTimeout: 10 * time.Millisecond,
	}
}

func newKafkaReader(config *conf.Configuration, topic, groupId string, brokers []string) *kafka.Reader {
	dialer, err := GetDialer(config)
	if err != nil {
		panic(err)
	}

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupId,
		MinBytes:    1,        // same value of Shopify/sarama
		MaxBytes:    57671680, // java client default
		Dialer:      dialer,
		ErrorLogger: kafka.LoggerFunc(logf),
		// Logger:      kafka.LoggerFunc(logf),
	})
}

func GetDialer(c *conf.Configuration) (*kafka.Dialer, error) {
	if c.Env == "dev" {
		return &kafka.Dialer{}, nil
	}

	switch c.Kafka.AuthenType {
	case "SASL":
		caCert, err := os.ReadFile(c.Kafka.CaCertPath)
		if err != nil {
			return nil, fmt.Errorf("error reading ca cert file: %v", err)
		}

		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			return nil, errors.New("failed to parse CA Certificate file")
		}
		tlsConfig := &tls.Config{
			RootCAs:            caCertPool,
			InsecureSkipVerify: true,
		}
		mechanism, err := scram.Mechanism(scram.SHA256, c.Kafka.User, c.Kafka.Password)
		if err != nil {
			return nil, err
		}
		dialer := &kafka.Dialer{
			Timeout:       KAFKA_TIME_OUT,
			KeepAlive:     time.Hour,
			DualStack:     true,
			TLS:           tlsConfig,
			SASLMechanism: mechanism,
		}
		return dialer, nil
	case "SSL":
		keypair, err := tls.LoadX509KeyPair(
			c.Kafka.TlsCertPath, c.Kafka.TlsKeyPath,
		)
		if err != nil {
			return nil, err
		}
		caCert, err := os.ReadFile(c.Kafka.CaCertPath)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			return nil, errors.New("failed to parse CA Certificate file")
		}
		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{keypair},
			RootCAs:            caCertPool,
			InsecureSkipVerify: true,
		}
		dialer := &kafka.Dialer{
			Timeout:   KAFKA_TIME_OUT,
			KeepAlive: time.Hour,
			DualStack: true,
			TLS:       tlsConfig,
		}
		return dialer, nil
	default:
		return nil, errors.New(c.Kafka.AuthenType + ": Kafka authentication type is not supported")
	}
}

func GetTransport(c *conf.Configuration) (*kafka.Transport, error) {
	if c.Env == "dev" {
		return &kafka.Transport{}, nil
	}

	switch c.Kafka.AuthenType {
	case "SASL":
		caCert, err := os.ReadFile(c.Kafka.CaCertPath)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			return nil, errors.New("failed to parse CA Certificate file")
		}
		tlsConfig := &tls.Config{
			RootCAs:            caCertPool,
			InsecureSkipVerify: true,
		}
		mechanism, err := scram.Mechanism(scram.SHA256, c.Kafka.User, c.Kafka.Password)
		if err != nil {
			return nil, err
		}
		return &kafka.Transport{
			TLS:  tlsConfig,
			SASL: mechanism,
		}, nil
	case "SSL":
		keypair, err := tls.LoadX509KeyPair(
			c.Kafka.TlsCertPath, c.Kafka.TlsKeyPath,
		)
		if err != nil {
			return nil, err
		}
		caCert, err := os.ReadFile(c.Kafka.CaCertPath)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			return nil, errors.New("failed to parse CA Certificate file")
		}
		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{keypair},
			RootCAs:            caCertPool,
			InsecureSkipVerify: true,
		}

		return &kafka.Transport{
			TLS: tlsConfig,
		}, nil
	default:
		return nil, errors.New(c.Kafka.AuthenType + ": Kafka authentication type is not supported")
	}
}

func logf(msg string, a ...interface{}) {
	log.LG.Fatalf(msg, a...)
}
