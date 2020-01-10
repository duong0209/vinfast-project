package controller

import (
	//"vinfast-project/database"
	// "fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"log"
	"strings"
	"vinfast-project/model"
	// "net/http"
)

const DayTime = 10 * 10 * 10 * 60 * 60 * 24

// dang nhap
func (gc *Controller) Login(c *gin.Context) {
	var user model.UserJSON
	c.BindJSON(&user)
	err := gc.DB.Raw(`
	      SELECT id, user_name
	      FROM user
	      WHERE user_name = ?
    `, user.UserName).Scan(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {

			gc.DB.Exec(`
		INSERT INTO user (user_name) values (?)
	`, user.UserName).Scan(&user)
			gc.DB.Raw(`
	    SELECT id, user_name
	    FROM user
	    WHERE user_name = ?
	`, user.UserName).Scan(&user)

			// tao wallet cho user
			log.Println("asdada", user.ID)
			err = gc.DB.Exec("INSERT INTO `wallet` (`id`, `money`, `user_id`) VALUES ('', ?, ?)", 0, user.ID).Error
			if err != nil {
				log.Println("insert wallet", err)
			}
			c.JSON(200, user)
			return
		}
		c.JSON(500, err)
		return
	}

	c.JSON(200, user)
	return
}

func (gc *Controller) UserInfo(c *gin.Context) {
	UserID := c.Param("id")
	var userinfo model.UserInfo
	errGetUserInfo := gc.DB.Raw(`
		SELECT user.id , user.user_name , wallet.money as wallet
		
		FROM user JOIN wallet ON user.id = wallet.user_id
		WHERE user.id = ?
	`, UserID).Scan(&userinfo).Error
	if errGetUserInfo != nil {
		log.Println(errGetUserInfo)

		return
	}
	c.JSON(200, userinfo)

}

// list danh sach xe thue duoc
func (gc *Controller) ListVehicle(c *gin.Context) {
	var vehicle []model.Vehicle
	errGetVehicle := gc.DB.Raw(`
		SELECT id, name, type, daily_price ,status
		
		FROM vehicle
		WHERE status = "Available"
	`).Scan(&vehicle).Error
	if errGetVehicle != nil {
		log.Println(errGetVehicle)

		return
	}
	c.JSON(200, vehicle)
}

// lay thong tin chi tiet 1 xe
func (gc *Controller) DetailVehicle(c *gin.Context) {
	VehicleId := c.Param("id")

	var detailvehicle model.Vehicle
	//log.Println(VehicleId)
	errGetDetailVehicle := gc.DB.Raw(`
	   SELECT vehicle.id, vehicle.name , vehicle.status , vehicle.type, vehicle.description , 
	   vehicle.daily_price , GROUP_CONCAT(image_vehicle.image_url) as image
	   
	   FROM vehicle JOIN image_vehicle ON image_vehicle.id_vehicle = vehicle.id
	 
	   WHERE vehicle.id = ?
	   
	   GROUP BY vehicle.id, vehicle.name , vehicle.status , vehicle.type, vehicle.description , 
	   vehicle.daily_price
	`, VehicleId).Scan(&detailvehicle).Error
	if errGetDetailVehicle != nil {
		log.Println(errGetDetailVehicle)
		return
	}
	detailvehicle.ImageList = strings.Split(detailvehicle.Image, ",")
	//log.Println(detailvehicle)
	c.JSON(200, detailvehicle)

}

func (gc *Controller) Booking(c *gin.Context) {
	var infoBooking model.BookingJSON
	var user model.User
	var vehicle model.Vehicle
	var wallet model.Wallet

	c.BindJSON(&infoBooking)
	log.Println(infoBooking)

	err := gc.DB.Raw("SELECT id, user_name as name FROM user WHERE id=?", infoBooking.IdUser).Scan(&user).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, map[string]string{"message": "Lỗi server 0"})
		return
	}

	err = gc.DB.Raw("SELECT* FROM vehicle WHERE id=?", infoBooking.VehicleID).Scan(&vehicle).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, map[string]string{"message": "Lỗi server 1"})
		return
	}

	err = gc.DB.Raw("SELECT* FROM wallet WHERE user_id=?", user.ID).Scan(&wallet).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, map[string]string{"message": "Lỗi server 2"})
		return
	}

	session := (infoBooking.EndDate - infoBooking.StartDate) / DayTime
	totalAmount := session * vehicle.DailyPrice

	// kiem tra vi
	if wallet.Money < totalAmount {
		c.JSON(http.StatusOK, map[string]string{"message": "Ví không đủ tiền"})
		return
	}
	// không thì insert
	errPostInfoBooking := gc.DB.Exec(`
	      INSERT INTO booking (user_id,vehicle_id,total_cost,start_date,end_date) values (?,?,?,?,?)
		`, user.ID, infoBooking.VehicleID, totalAmount, infoBooking.StartDate, infoBooking.EndDate).Error
	if errPostInfoBooking != nil {
		log.Println(errPostInfoBooking)
		c.JSON(http.StatusInternalServerError, map[string]string{"message": "Lỗi server 3"})
		return
	}
	errUpdateStatus := gc.DB.Exec(`
	      UPDATE vehicle SET status = "Rented" WHERE ID = ?
		`, infoBooking.VehicleID).Error
	if errUpdateStatus != nil {
		log.Println(errUpdateStatus)
		c.JSON(http.StatusInternalServerError, map[string]string{"message": "Lỗi server 4"})
		return
	}

	c.JSON(200, map[string]interface{}{
		"username":     user.Name,
		"vehicle":      vehicle.Name,
		"startdate":    infoBooking.StartDate,
		"enddate":      infoBooking.EndDate,
		"total_amount": totalAmount,
		"status":       "Rented",
	})

}

// api top up wallet
func (gc *Controller) Booking(c *gin.Context)

