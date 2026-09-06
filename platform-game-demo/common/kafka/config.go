package kafka

type Config struct {
	Brokers []string `json:"brokers" mapstructure:"brokers" yaml:"brokers"`
	// Consumers []TopicConfig `json:"consumers" mapstructure:"consumers" yaml:"consumers"`
}

// type TopicConfig struct {
// 	Topic       string `json:"topic" mapstructure:"topic" yaml:"topic"`
// 	GroupId     string `json:"group_id" mapstructure:"group_id" yaml:"group_id"`
// 	Concurrency int    `json:"concurrency" mapstructure:"concurrency" yaml:"concurrency"`
// }
