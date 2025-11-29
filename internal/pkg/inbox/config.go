package inbox

type PutConfig struct {
}

type UnboxConfig struct {
	MaxWorkers int `mapstructure:"maxWorkers"`
}

type Config struct {
	Put   PutConfig   `mapstructure:"put"`
	Unbox UnboxConfig `mapstructure:"unbox"`
}
