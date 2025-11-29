package messaging

const (
	Kafka Vendor = "kafka"
)

type InputConfig struct {
	VendorKey string
}

type OutputConfig struct {
}

type RoutingConfig struct {
	Input  map[RoutingKey]InputConfig
	Output map[RoutingKey]OutputConfig
}

type Config struct {
	Routing map[Vendor]RoutingConfig `mapstructure:"routing"`
}
