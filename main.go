package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
)

type Configuration struct {
	Database struct {
		Username       string            `yaml:"username"`
		Password       string            `yaml:"password"`
		Host           string            `yaml:"host"`
		Port           string            `yaml:"port"`
		Schema         string            `yaml:"schema"`
		ConnectOptions map[string]string `yaml:"connectoptions"`
	} `yaml:"database"`
}

func main() {
	viper.SetConfigName("mysql-ssl-test")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s\n", err.Error())
	}
	var config Configuration
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshaling config: %s\n", err.Error())
	}
	dbconfig := config.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbconfig.Username, dbconfig.Password, dbconfig.Host, dbconfig.Port, dbconfig.Schema)
	if len(dbconfig.ConnectOptions) > 0 {
		dsn += "?"
		params := url.Values{}
		for k, v := range config.Database.ConnectOptions {
			params.Add(k, v)
		}
		dsn += params.Encode()
	}
	fmt.Printf("%+v\n", dbconfig)
	fmt.Printf("[%s]\n", dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %s\n", err.Error())
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %s\n", err.Error())
	}
	fmt.Printf("Successfully connected to database [%s:%s]\n", dbconfig.Host, dbconfig.Port)
}
