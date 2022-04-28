package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	"github.com/gumper23/dbstuff/dbhelper"
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
	config, err := GetConfig()
	if err != nil {
		log.Fatalf("Error reading configuration: %s\n", err.Error())
	}

	db, err := config.GetMySQLDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %s\n", err.Error())
	}
	defer db.Close()
	dbconfig := config.Database
	fmt.Printf("Successfully connected to database [%s:%s]\n", dbconfig.Host, dbconfig.Port)

	row, cols, err := dbhelper.QueryRow(db, "select 22 as val")
	if err != nil {
		log.Fatalf("Error querying db: %s\n", err.Error())
	}
	for _, col := range cols {
		fmt.Printf("Col [%s] = [%s]\n", col, row[col])
	}
}

func GetConfig() (config *Configuration, err error) {
	viper.SetConfigName("mysql-ssl-test")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}

func (config *Configuration) GetMySQLDSN() (dsn string) {
	dbconfig := config.Database
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbconfig.Username, dbconfig.Password, dbconfig.Host, dbconfig.Port, dbconfig.Schema)
	if len(dbconfig.ConnectOptions) > 0 {
		dsn += "?"
		params := url.Values{}
		for k, v := range config.Database.ConnectOptions {
			params.Add(k, v)
		}
		dsn += params.Encode()
	}
	return
}

func (config *Configuration) GetMySQLDB() (db *sql.DB, err error) {
	dsn := config.GetMySQLDSN()
	if db, err = sql.Open("mysql", dsn); err != nil {
		return
	}
	err = db.Ping()
	return
}
