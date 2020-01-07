package model

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/json-iterator/go"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	ID   int `json:"user_id,omitempty"`
	Name string `json:"user_name,omitempty"`
}

type Vehicle struct {
	ID     int
	Name   string
	Status bool
}

type VehicalType struct {
	ID          int
	Name        string
	Description string
	Images      string
	DailyPrice  int
}

type Booking struct {
	ID        int
	IDVehical int
	IDUser    int
	TotalCost int
	StartDate int
	EndDate   int
}

type Wallet struct {
	ID     int
	IDUser int
	Point  int
}

type Config struct {
	Database struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
		Address  string `json:"address"`
	} `json:"database"`
}

func DecodeDataFromJsonFile(f *os.File, data interface{}) error {
	jsonParser := jsoniter.NewDecoder(f)
	err := jsonParser.Decode(&data)
	if err != nil {
		return err
	}

	return nil
}

func SetupConfig() Config {
	var conf Config

	// Đọc file config.dev.json
	configFile, err := os.Open("config.local.json")
	if err != nil {
		// Nếu không có file config.dev.json thì đọc file config.default.json
		configFile, err = os.Open("config.default.json")
		if err != nil {
			panic(err)
		}
		defer configFile.Close()
	}
	defer configFile.Close()

	// Parse dữ liệu JSON và bind vào conf
	err = DecodeDataFromJsonFile(configFile, &conf)
	if err != nil {
		log.Println("cannot read config file", err)
		panic(err)
	}

	return conf
}

func ConnectDb(user string, password string, database string, address string) *gorm.DB {
	connectionInfo := fmt.Sprintf(`%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local`, user, password, address, database)

	db, err := gorm.Open("mysql", connectionInfo)
	if err != nil {
		panic(err)
	}
	return db
}
