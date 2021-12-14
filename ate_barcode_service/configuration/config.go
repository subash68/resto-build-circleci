package configuration

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type ServicePortConfig struct {
	GRPCPort int `json:"GRPC_PORT"`
	HTTPPort int `json:"HTTP_PORT"`
}

type LoggingConfig struct {
	LogLevel      int    `json:"LOG_LEVEL"`
	LogTimeFormat string `json:"LOG_TIME_FORMAT"`
}

type DatabaseConfig struct {
	DatabaseHost     string `json:"DB_HOST"`
	DatabaseUser     string `json:"DB_USER_NAME"`
	DatabasePassword string `json:"DB_USER_PASSWORD"`
	DatabaseSchema   string `json:"DB_NAME"`
}

var servicePortConfig *ServicePortConfig
var loggingConfig *LoggingConfig
var databaseConfig *DatabaseConfig

func PortConfig() ServicePortConfig { return *servicePortConfig }
func LogConfig() LoggingConfig      { return *loggingConfig }
func DbConfig() DatabaseConfig      { return *databaseConfig }

func LoadConfig() {
	file, err := os.Open("./configuration/config.json")

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	byteValue, _ := ioutil.ReadAll(file)

	err = json.Unmarshal(byteValue, &servicePortConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = json.Unmarshal(byteValue, &loggingConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = json.Unmarshal(byteValue, &databaseConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
}
