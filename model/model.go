package model

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"user_name,omitempty"`
}
type UserInfo struct {
	ID       int
	UserName string `json:"user_name"`
	Wallet   Wallet    `json:"wallet"`
	
}
type Vehicle struct {
	ID          int      `json:"ID"`
	Name        string   `json:"Name"`
	Type        string   `json:"Type"`
	Status      string   `json:"Status"`
	Description string   `json:"Description"`
	Image       string   `json:"Image"`
	ImageList   []string `json:"image_list"`
	DailyPrice  int      `json:"DailyPrice"`
}

type Booking struct {
	ID        int `json:"id"`
	IDVehical int `json:"vehicle_id"`
	IDUser    int `json:"user_id"`
	TotalCost int `json:"total_cost"`
	StartDate int `json:"start_date"`
	EndDate   int `json:"end_date"`
}

type Wallet struct {
	IDWallet int `json:"id"`
	IDUser   int `json:"user_id"`
	Money    int `json:"money"`
}

type UserJSON struct {
	ID       int    `json:"id"`
	UserName string `json:"name"`
}

type BookingJSON struct {
	UserID    int `json:"iduser"`
	VehicleID int `json:"vehicleid"`
	StartDate int `json:"startdate"`
	EndDate   int `json:"enddate"`
}

type TopUpWalletJSON struct {
	UserID int `json:"userid"`
	Amount	int `json:"amount"`
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


type Meta struct {
	Code    int
	Message string
}
type Respond struct {
	Data interface{}
	Meta Meta
}

func MakeRespond(data interface{}, code int, msg string) Respond {
	return Respond{
		Data: data,
		Meta: Meta{
			Code:    code,
			Message: msg,
		},
	}
}