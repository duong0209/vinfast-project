package controller

import (
	//"vinfast-project/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"vinfast-project/model"
	"log"
)

func (gc *Controller)Login(c *gin.Context) {
	var user model.User
	var Res interface{}
	if err := c.BindJSON(&user); err != nil {
		fmt.Println(err)
		c.JSON(400, Res)
		return

	c.JSON(200, user)
		return	
	}

	

}


func (gc *Controller)ListVehicle(c *gin.Context) {
	var vehicle []model.Vehicle
	errGetVehicle := gc.DB.Raw(`
		SELECT id, name, images , status
		 FROM vehicle
	`).Scan(&vehicle).Error
	if errGetVehicle != nil {
		log.Println(errGetVehicle)
		
			
		return
	}
	c.JSON(200, vehicle)
}	


func (gc *Controller)DetailVehicle(c *gin.Context) {
	VehicleId := c.Param("id")

	var detailvehicle model.Vehicle
	//log.Println(VehicleId)
	errGetDetailVehicle := gc.DB.Raw(`
		SELECT * 
		FROM vehicle
		WHERE id = ?
	`,  VehicleId).Scan(&detailvehicle).Error
	if errGetDetailVehicle != nil {
		log.Println(errGetDetailVehicle)
		return
	}
	//log.Println(detailvehicle)
	 c.JSON(200, detailvehicle)

}

