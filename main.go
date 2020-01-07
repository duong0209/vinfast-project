package main

import (
	"vinfast-project/model"
	"vinfast-project/controller"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	router := gin.Default()
	// c := controller.NewController()

	config := model.SetupConfig()
	db := model.ConnectDb(config.Database.User, config.Database.Password, config.Database.Database, config.Database.Address)
	defer db.Close()

	// var mysql mysql.sql

	// //api USER
	router.GET("/login", controller.Login)
	// router.POST("/users", AddUser)

	// //api vehicle
	// router.GET("/listvehicle", ListVehicle)
	// router.GET("/detailvehcle/:id", DetailVehicle)
	// router.PUT("/vehicle/status/", VehicleStatus)

	// //api booking
	// router.POST("/booking", Booking)

	// //api wallet
	// router.PUT("/wallet/point/:id", Point)

	router.Run("0.0.0.0:8089")
}
