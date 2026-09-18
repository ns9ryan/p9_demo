package config

type SwaggerConfig struct {
	Protocol string `yaml:"Protocol" json:"Protocol"`
}

type KafkaConfig struct {
	Brokers []string `json:"brokers" mapstructure:"brokers" yaml:"brokers"`
}
