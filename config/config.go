package config

import (
	"bytes"

	tio "iothub"
	"iothub/pkg/log"

	_ "embed"

	"iothub/db/mysql"
	"iothub/db/sqlite"

	"github.com/spf13/viper"
)

var (
	Version   = ""
	GitCommit = ""

	defaultConfigYaml []byte
)

func init() {
	defaultConfigYaml = tio.DefaultConfigYaml
}

const (
	// DBSqlite DB type
	DBSqlite = "sqlite"
	DBMySQL  = "mysql"

	// Connector type

	// ConnectorEmqx EMQX MQTT broker
	ConnectorEmqx = "emqx"
	// ConnectorMqttEmbed MQTT broker embedded in iothub
	ConnectorMqttEmbed = "embed"
)

type UserPassword struct {
	Name     string
	Password string
}

type InnerMqttBroker struct {
	TcpPort    int
	TcpSslPort int
	WsPort     int
	WssPort    int
	CertFile   string
	KeyFile    string
	SuperUsers []UserPassword
}
type Config struct {
	Log struct {
		Level string
	}
	API struct {
		Port      int
		BasicAuth UserPassword
	}
	DB struct {
		Typ    string `mapstructure:"type"`
		Mysql  mysql.Config
		Sqlite sqlite.Config
	}
	Connector Connector
}

func ReadConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/iothub/")
	viper.AddConfigPath("$HOME/.iothub")
	viper.AddConfigPath(".")

	err := viper.ReadConfig(bytes.NewReader(defaultConfigYaml))
	if err != nil {
		log.Fatalf("Error read default config file: %v", err)
	}

	err = viper.MergeInConfig()
	if err != nil {
		log.Fatalf("Error read config file: %v", err)
	}
	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("Error read config file content: %v", err)
	}
	return cfg
}

type Connector struct {
	Typ        string `mapstructure:"type"`
	MqttClient MqttClientConfig
	MqttBroker InnerMqttBroker
	Emqx       EmqxAdapterConfig
}

type MqttClientConfig struct {
	ClientId     string `json:"clientId"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	User         string `json:"user"`
	Password     string `json:"password"`
	WillTopic    string `json:"WillTopic"`
	WillPayload  string `json:"willPayload"`
	CleanSession *bool  `json:"cleanSession"`
}

type EmqxAdapterConfig struct {
	ApiPrefix   string // like http://localhost:18083
	ApiUser     string
	ApiPassword string
}
