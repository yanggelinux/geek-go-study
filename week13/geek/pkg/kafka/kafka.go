package kafka

import (
	"geek/pkg/log"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

type Producer struct {
	brokers []string
	version string
}

func NewProducer(brokers []string, version string) *Producer {

	return &Producer{brokers, version}
}

func (p *Producer) MakeAsyncProducer() (sarama.AsyncProducer, error) {

	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewHashPartitioner
	version, err := sarama.ParseKafkaVersion(p.version)
	config.Version = version
	producer, err := sarama.NewAsyncProducer(p.brokers, config)
	return producer, err
}

func (p *Producer) AsyncSendData(producer sarama.AsyncProducer, topic string, data []byte, partition int32) {

	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.ByteEncoder(data)
	msg.Partition = partition
	producer.Input() <- msg
}

func (p *Producer) AsyncSendData2(producer sarama.AsyncProducer, topic string, data []byte) {

	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.ByteEncoder(data)
	producer.Input() <- msg
}

func (p *Producer) MakeSyncProducer() (sarama.SyncProducer, error) {

	config := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//自定义partition发送消息
	config.Producer.Partitioner = sarama.NewHashPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//注意，版本设置不对的话，kafka会返回很奇怪的错误，并且无法成功发送消息
	version, err := sarama.ParseKafkaVersion(p.version)
	config.Version = version
	producer, err := sarama.NewSyncProducer(p.brokers, config)
	return producer, err
}

func (p *Producer) SendData(producer sarama.SyncProducer, topic string, data []byte, key string) {

	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.ByteEncoder(data)
	msg.Key = sarama.StringEncoder(key)

	pid, offset, err := producer.SendMessage(msg)
	_ = pid
	_ = offset
	if err != nil {
		log.Logger.Error("send message failed", zap.Error(err))
	}
}

func (p *Producer) SendData2(producer sarama.SyncProducer, topic string, data []byte) {

	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.ByteEncoder(data)

	pid, offset, err := producer.SendMessage(msg)
	_ = pid
	_ = offset
	if err != nil {
		log.Logger.Error("send message failed", zap.Error(err))
	}
}

type Consumer struct {
	brokers []string
	groupID string
	topics  []string
	version string
}

func NewConsumer(brokers, topics []string, groupID, version string) *Consumer {
	return &Consumer{
		brokers,
		groupID,
		topics,
		version,
	}
}

func (c *Consumer) MakeComsumer(offset string, autoCommit bool) (sarama.ConsumerGroup, error) {
	// init (custom) configs, set mode to ConsumerModePartitions
	config := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(c.version)

	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Group.Rebalance.Retry.Max = 20
	config.Consumer.Offsets.AutoCommit.Enable = autoCommit
	if offset == "new" {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	} else {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	config.Version = version

	consumer, err := sarama.NewConsumerGroup(c.brokers, c.groupID, config)
	return consumer, err
}
