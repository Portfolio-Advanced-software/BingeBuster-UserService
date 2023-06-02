package config

import "github.com/spf13/viper"

type Config struct {
	Port              string `mapstructure:"PORT"`
	MongoDBUser       string `mapstructure:"MONGODB_USER"`
	MongoDBPwd        string `mapstructure:"MONGODB_PWD"`
	MongoDBCluster    string `mapstructure:"MONGODB_CLUSTER"`
	MongoDBDb         string `mapstructure:"MONGODB_DB"`
	MongoDBCollection string `mapstructure:"MONGODB_COLLECTION"`
	RabbitMQUser      string `mapstructure:"RABBITMQ_USER"`
	RabbitMQPwd       string `mapstructure:"RABBITMQ_PWD"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./config/")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
