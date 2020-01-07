package controller

import (
	"vinfast-project/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"vinfast-project/model"
)

func Login(c *gin.Context) {
	var user model.User
	var Res interface{}
	if err := c.BindJSON(&user); err != nil {
		fmt.Println(err)
		c.JSON(400, Res)
		return

	c.JSON(200, user)
		return	
	}

	if err := database.AddUser(user); err != nil {
		c.JSON(400, Res)
	} else {
		c.JSON(200, Res)
	}


func  GetItems(c *gin.Context) {
		var Vehicle []model.Vehicle
		var itemDao database.ItemDao
		itemDao = controller.dao
		items, err := itemDao.FetchItems()
		if err != nil {
			c.JSON(http.StatusInternalServerError, utility.MakeResponse(500, "Internal server error!", nil))
			return
		}
		c.JSON(http.StatusOK, utility.MakeResponse(200, "Request successful", items))
	
	
	
	
	


	// if user2, err := mysql.GetUserFromUP(user.ID); err != nil {
	// 	c.JSON(400, Res)
	// } else {
	// 	c.JSON(200, user2)
	// }
}

// AddUser add a user into database
func AddUser(c *gin.Context) {
	var user model.User
	var Res interface{}
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, Res)
		return
	}
	c.JSON(200, nil)
	return
	// if err := mysql.AddUser(user); err != nil {
	// 	c.JSON(400, Res)
	// } else {
	// 	c.JSON(200, Res)
	// }
}

//  point into wallet

func EditWallet(c *gin.Context) {
	var res interface{}
	var OneWallet model.Wallet
	if err := c.BindJSON(&OneWallet); err != nil {
		c.JSON(400, res)
		return
	}
	c.JSON(200, nil)
	return
	// if err := mysql.EditBill(OneWallet); err != nil {
	// 	c.JSON(400, res)
	// } else {
	// 	c.JSON(200, res)
	// }
}

func VehicleList(c *gin.Context) {
	var res interface{}
	var OneVehicle model.Vehicle
	if err := c.BindJSON(&OneVehicle); err != nil {
		c.JSON(400, res)
		return
	}
	c.JSON(200, nil)
	return
	// if err := mysql.Vehicle(OneVehicle); err != nil {
	// 	c.JSON(400, res)
	// } else {
	// 	c.JSON(200, res)
	// }
}
