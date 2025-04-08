package conf_type

// Kafka Cluster Configure
type Kafka struct {
	KafkaClusterStartUp bool `yaml:"startup"`
	// example: 10.0.1.89:9092
	KafkaClusterURL []string `yaml:"cluster_addr"`
}
