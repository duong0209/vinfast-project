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
	ginController := controller.NewController()

	config := model.SetupConfig()
	db := model.ConnectDb(config.Database.User, config.Database.Password, config.Database.Database, config.Database.Address)
	defer db.Close()

	ginController.DB = db
	ginController.Config = config


	

	// //api USER
	router.GET("/login", ginController.Login)
	// router.POST("/users", AddUser)

	// //api vehicle
	 router.GET("/listvehicle",ginController.ListVehicle)
	 router.GET("/detailvehicle/:id",ginController.DetailVehicle)
	// router.PUT("/vehicle/status/", VehicleStatus)

	// //api booking
	// router.POST("/booking", Booking)

	// //api wallet
	// router.PUT("/wallet/point/:id", Point)

	router.Run(":8089")
}
