package msgqueue

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/astraprotocol/affiliate-system/conf"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/util/log"
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
	KAFKA_NOTI_GROUP_ID              = "reward-notification-backend"
	KAFKA_TOPIC_NOTI_SMS             = "notification-SMS"
	KAFKA_TOPIC_NOTI_EMAIL           = "notification-EMAIL"
	KAFKA_TOPIC_NOTI_STATUS_RESPONSE = "notification-status"
	KAFKA_TOPIC_NOTI_APP_MESSAGE     = "notification-worker-app-noti-message"
)

type MessageQueue struct {
	Reader  *kafka.Reader
	Writer  *kafka.Writer
	Brokers []string
	User    string
	Sigchan chan os.Signal
}

func GetDialer(conf *conf.Configuration) (*kafka.Dialer, error) {
	if conf.Env == "dev" {
		return nil, nil
	}

	switch conf.Kafka.AuthenType {
	case "SASL":

		caCert, err := os.ReadFile(conf.Kafka.CaCertPath)
		if err != nil {
			return nil, fmt.Errorf("error reading ca cert file: %v", err)
		}
		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			return nil, fmt.Errorf("error appending ca cert")
		}
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			RootCAs:            caCertPool,
		}

		mechanism, err := scram.Mechanism(scram.SHA256, conf.Kafka.User, conf.Kafka.Password)
		if err != nil {
			return nil, fmt.Errorf("error setup scram mechanism: %v", err)
		}

		return &kafka.Dialer{
			Timeout:       KAFKA_TIME_OUT,
			DualStack:     true,
			TLS:           tlsConfig,
			SASLMechanism: mechanism,
		}, nil
	case "SSL":
		keypair, err := tls.LoadX509KeyPair(
			conf.Kafka.TlsCertPath, conf.Kafka.TlsKeyPath,
		)
		if err != nil {
			return nil, err
		}
		caCert, err := os.ReadFile(conf.Kafka.CaCertPath)
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
		return &kafka.Dialer{
			Timeout:   KAFKA_TIME_OUT,
			KeepAlive: time.Hour,
			DualStack: true,
			TLS:       tlsConfig,
		}, nil
	default:
		return nil, errors.New(conf.Kafka.AuthenType + ": Kafka authentication type is not supported")
	}
}

func NewMsgQueue(topic, groupId string) *MessageQueue {
	mq := &MessageQueue{}

	config := conf.GetConfiguration()

	mq.User = config.Kafka.User
	mq.Brokers = []string{config.Kafka.KafkaURL}

	dialer, err := GetDialer(config)
	if err != nil {
		panic(err)
	}

	mq.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:          mq.Brokers,
		Topic:            topic,
		GroupID:          groupId,
		MinBytes:         1,        // same value of Shopify/sarama
		MaxBytes:         57671680, // java client default
		MaxWait:          KAFKA_NEW_DATA_MAX_WAIT,
		ReadBatchTimeout: KAFKA_READ_BATCH_TIME_OUT,
		Dialer:           dialer,
		ErrorLogger:      kafka.LoggerFunc(logf),
		RebalanceTimeout: 5 * time.Minute,
		// Logger: kafka.LoggerFunc(logf),
	})

	mq.Writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  mq.Brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Dialer:   dialer,
	})

	return mq
}

func logf(msg string, a ...interface{}) {
	log.LG.Fatalf(msg, a...)
}
