package main

import (
	"database/sql"
	"fmt"

	"jatis/storage"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func main() {
	conf := LoadConfig()
	db := DBInit(conf)
	storage.New(db)
}

func LoadConfig() *viper.Viper {
	viper.SetConfigFile("config.env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return viper.GetViper()
}

func DBInit(conf *viper.Viper) *sql.DB {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s", conf.Get("DB_USERNAME"), conf.Get("DB_PASSWORD"), conf.Get("DB_URL"),
		conf.Get("DB_SCHEMA"))
	db, err := sql.Open(conf.Get("DB_DRIVER").(string), dataSource)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
