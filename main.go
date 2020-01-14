package main

import (
	"vinfast-project/controller"
	"vinfast-project/model"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	router := gin.Default()
	// c := controller.NewController()
	ginController := controller.NewController()

	config := model.SetupConfig()
	db := model.ConnectDb(config.Database.User, config.Database.Password, config.Database.Database, config.Database.Address)
	defer db.Close()

	ginController.DB = db
	ginController.Config = config

	// //api USER
	router.POST("/login", ginController.Login)
	router.GET("/userinfo/:id", ginController.UserInfo)
	router.POST("/topup",ginController.TopUpWallet)

	// //api vehicle
	router.GET("/listvehicle", ginController.ListVehicle)
	router.GET("/detailvehicle/:id", ginController.DetailVehicle)

	// //api booking
	router.POST("/booking", ginController.Booking)
	
	

	router.Run(":8089")
}
