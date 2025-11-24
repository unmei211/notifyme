package inbox

type PutConfig struct {
}

type PollConfig struct {
}

type Config struct {
	put  PutConfig  `mapstructure:"put"`
	poll PollConfig `mapstructure:"poll"`
}
