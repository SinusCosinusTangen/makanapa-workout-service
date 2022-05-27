package config

import (
	"fmt"
	"net"
	nurl "net/url"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	Port        int
	Environment string
	Debug       bool

	AuthHost string
	MQHost   string

	DBHost     string
	DBPort     int
	DBDatabase string
	DBUsername string
	DBPassword string
	SSLMode    string

	ESHost     string
	ESUsername string
	ESPassword string
}

func InitializeConfig() {
	viper.SetConfigName(".env") // allow directly reading from .env file
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	// Application environment
	AppConfig.Port = viper.GetInt("PORT")
	AppConfig.Environment = viper.GetString("ENVIRONMENT")
	AppConfig.Debug = viper.GetBool("DEBUG")

	if AppConfig.Environment == "dev" {
		AppConfig.SSLMode = "disable"
	} else {
		AppConfig.SSLMode = "require"
	}

	// Database environment
	dbUrl := viper.GetString("DATABASE_URL")

	if dbUrl != "" {
		dbInfo, err := ParseURL(dbUrl)

		if err == nil {
			AppConfig.DBUsername = dbInfo[0]
			AppConfig.DBPassword = dbInfo[1]
			AppConfig.DBHost = dbInfo[2]
			AppConfig.DBPort, err = strconv.Atoi(dbInfo[3])
			AppConfig.DBDatabase = dbInfo[4]

			if err == nil {
				return
			}
		}
	}

	// Authentication
	AppConfig.AuthHost = viper.GetString("AUTH_HOST")

	// Message Queue
	AppConfig.MQHost = viper.GetString("MQ_HOST")

	// Database
	AppConfig.DBUsername = viper.GetString("DATABASE_USERNAME")
	AppConfig.DBPassword = viper.GetString("DATABASE_PASSWORD")
	AppConfig.DBHost = viper.GetString("DATABASE_HOST")
	AppConfig.DBPort = viper.GetInt("DATABASE_PORT")
	AppConfig.DBDatabase = viper.GetString("DATABASE_DATABASE")

	// Elasticsearch
	AppConfig.ESHost = viper.GetString("ELASTIC_SEARCH_URL")
	AppConfig.ESUsername = viper.GetString("ELASTIC_SEARCH_USERNAME")
	AppConfig.ESPassword = viper.GetString("ELASTIC_SEARCH_PASSWORD")
}

// Example:
//	"postgres://bob:secret@1.2.3.4:5432/mydb?sslmode=verify-full"
//
// converts to:
//	"user=bob password=secret host=1.2.3.4 port=5432 dbname=mydb sslmode=verify-full"
//
// This function is modified from github.com/lib/pq

func ParseURL(url string) ([]string, error) {
	var kvs []string

	u, err := nurl.Parse(url)
	if err != nil {
		return kvs, err
	}

	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return kvs, fmt.Errorf("invalid connection protocol: %s", u.Scheme)
	}

	escaper := strings.NewReplacer(`'`, `\'`, `\`, `\\`)
	accrue := func(v string) {
		if v != "" {
			kvs = append(kvs, escaper.Replace(v))
		}
	}

	if u.User != nil {
		v := u.User.Username()
		accrue(v)

		v, _ = u.User.Password()
		accrue(v)
	}

	if host, port, err := net.SplitHostPort(u.Host); err != nil {
		accrue(u.Host)
	} else {
		accrue(host)
		accrue(port)
	}

	if u.Path != "" {
		accrue(u.Path[1:])
	}

	q := u.Query()
	for k := range q {
		accrue(q.Get(k))
	}

	return kvs, nil
}
