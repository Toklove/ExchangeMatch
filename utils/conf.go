package utils

type MeConfig struct {
	Debug     bool     `yaml:"Debug"`
	LogLevel  string   `yaml:"LogLevel"`
	MqConf    RabbitMQ `yaml:"rabbitmq"`
	CacheConf Cache    `yaml:"redis"`
	GRPCConf  GRPC     `yaml:"grpc"`
	MeConf    MEngine  `yaml:"MeConf"`
}

type Cache struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
}

type RabbitMQ struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type GRPC struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type MEngine struct {
	Accuracy int `yaml:"accuracy"`
}
