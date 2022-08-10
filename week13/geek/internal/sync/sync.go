package sync

import (
	"context"
	"geek/global"
	"geek/pkg/kafka"
	"geek/pkg/log"
	"github.com/Shopify/sarama"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"strings"
	"sync"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type SyncService struct {
	kafkaProducer kafka.Producer
	syncProducer  sarama.SyncProducer
	producerTopic string
}

func NewSyncService(kafkaProducer kafka.Producer, syncProducer sarama.SyncProducer, producerTopic string) *SyncService {
	return &SyncService{
		kafkaProducer,
		syncProducer,
		producerTopic,
	}
}

func (s *SyncService) Sync(msgValue []byte) error {
	s.kafkaProducer.SendData2(s.syncProducer, s.producerTopic, msgValue)
	return nil
}

type ConsumerHandler struct {
	ready        chan bool
	transIncrSvc *SyncService
}

func (ch *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	close(ch.ready)
	return nil
}

func (ch *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (ch *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for message := range claim.Messages() {
		log.Logger.Debug("Recv", zap.ByteString("msgValue", message.Value))
		session.MarkMessage(message, "")
		//数据转换

		err := ch.transIncrSvc.Sync(message.Value)
		if err != nil {
			log.Logger.Error("data_ext trans or send data error", zap.Error(err))
		}

	}
	return nil
}

func SyncData(ctx context.Context, exit chan struct{}) {
	//数据同步 程序

	consumerTopics := strings.Split(global.KafkaSetting.ConsumerTopic, ",")
	consumerBrokers := strings.Split(global.KafkaSetting.ConsumerBrokers, ",")
	producerBrokers := strings.Split(global.KafkaSetting.ProducerBrokers, ",")
	consumerGroupID := global.KafkaSetting.ConsumerGroupID
	consumerVersion := global.KafkaSetting.ConsumerVersion
	consumerOffset := global.KafkaSetting.ConsumerOffset
	producerTopic := global.KafkaSetting.ProducerTopic
	producerVersion := global.KafkaSetting.ProducerVersion

	kafkaProducer := kafka.NewProducer(producerBrokers, producerVersion)
	kafkaConsumer := kafka.NewConsumer(consumerBrokers, consumerTopics, consumerGroupID, consumerVersion)
	consumer, err := kafkaConsumer.MakeComsumer(consumerOffset, true)

	defer func() {
		exit <- struct{}{}
		log.Logger.Info("sync incr process exit!...")
	}()

	if err != nil {
		log.Logger.Error("error creating consumer group client", zap.Error(err))
		return
	}
	defer func() {
		if err = consumer.Close(); err != nil {
			log.Logger.Error("error closing client: %v", zap.Error(err))
			return
		}
		log.Logger.Info("sarama consumer close and exit!...")
	}()
	syncProducer, err := kafkaProducer.MakeSyncProducer()
	if err != nil {
		log.Logger.Error("make sync producer error", zap.Error(err))
		return
	}
	defer func() {
		err = syncProducer.Close()
		if err != nil {
			log.Logger.Error("close sync producer error", zap.Error(err))
		}
	}()
	transIncrSvc := NewSyncService(*kafkaProducer, syncProducer, producerTopic)
	consumerHandler := ConsumerHandler{
		ready:        make(chan bool),
		transIncrSvc: transIncrSvc,
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := consumer.Consume(ctx, consumerTopics, &consumerHandler); err != nil {
				log.Logger.Error("Error from consumer", zap.Error(err))
			}
			if ctx.Err() != nil {
				return
			}
			consumerHandler.ready = make(chan bool)
		}
	}()
	<-consumerHandler.ready
	log.Logger.Info("Sarama consumer up and running!...")

	select {
	case <-ctx.Done():
		log.Logger.Info("terminating: context cancelled")
	}
	wg.Wait()
}
