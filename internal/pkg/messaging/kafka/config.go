package kafka

import msg "github.com/unmei211/notifyme/internal/pkg/messaging"

type ProducerConfig struct {
	Topic        Topic `mapstructure:"topic"`
	Async        bool  `mapstructure:"async"`
	BatchSize    int   `mapstructure:"batchSize"`
	BatchTimeout int   `mapstructure:"batchTimeout"`
}

//type ConsumerConfig struct {
//	Topic Topic `mapstructure:"topic"`
//}

//type ConsumeConfig struct {
//	GroupId   string                            `mapstructure:"groupId"`
//	Consumers map[msg.RoutingKey]ConsumerConfig `mapstructure:"consumers"`
//}

type FetchingConfig struct {
	GroupId string `mapstructure:"groupId"`
}

type Config struct {
	Addr      []string                          `mapstructure:"addr"`
	Producers map[msg.RoutingKey]ProducerConfig `mapstructure:"producers"`
	Fetching  FetchingConfig                    `mapstructure:"fetching"`
	//Consume   ConsumeConfig                     `mapstructure:"consume"`
}
