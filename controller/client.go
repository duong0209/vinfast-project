package controller

import (

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
			
			err = gc.DB.Exec("INSERT INTO `wallet` (`id`, `money`, `user_id`) VALUES ('', ?, ?)", 0, user.ID).Error
			if err != nil {
				log.Println("insert wallet", err)
			}
			c.JSON(200, model.MakeRespond(user, 200, "tao wallet cho user thanh cong"))
			return
		}
		c.JSON(400, model.MakeRespond(nil, 400, "Loi ko dang nhap dc"))
		return
	}

	c.JSON(200, model.MakeRespond(user, 200, "dang nhap thanh cong"))
	return
}
// api lay thong tin user
func (gc *Controller) UserInfo(c *gin.Context) {
	UserID := c.Param("id")
	var userinfo model.UserInfo
	errGetUserInfo := gc.DB.Raw(`
		SELECT user.id , user.user_name 
		
		FROM user JOIN wallet ON user.id = wallet.user_id
		WHERE user.id = ?
	`, UserID).Scan(&userinfo).Error
	if errGetUserInfo != nil {
		log.Println(errGetUserInfo)
		c.JSON(400, model.MakeRespond(nil, 400, "Loi lay thong tin user"))
        return
	}
	// lay wallet
	errGetWallet := gc.DB.Raw(`
		SELECT wallet.id as id_wallet, wallet.user_id as id_user, wallet.money as money
		FROM wallet
		WHERE user_id = ?
	`, UserID).Scan(&userinfo.Wallet).Error
	if errGetWallet != nil {
		log.Println(errGetWallet)
		c.JSON(400, model.MakeRespond(nil, 400, "Loi lay thong tin vi"))

		return
	}
    // lay thong tin dat xe
	errGetBill := gc.DB.Raw(`
		SELECT booking.id as id , booking.user_id as id_user , booking.vehicle_id as id_vehicle ,
		       booking.total_cost as total_cost , booking.start_date as start_date , booking.end_date as end_date
		FROM booking 
		WHERE user_id = ?
	`, UserID).Scan(&userinfo.Bill).Error
	if errGetBill != nil {
		if errGetBill.Error() == "record not found"{
			c.JSON(200, model.MakeRespond(userinfo, 200, "Chua dat xe nao"))
			return
		}
		log.Println(errGetBill)
		c.JSON(400, model.MakeRespond(nil, 400, "Loi lay thong tin dat xe"))

		return
	}
	c.JSON(200, model.MakeRespond(userinfo, 200, "lay thong tin user thong tin thanh cong"))
}

// list danh sach xe thue duoc
func (gc *Controller) ListVehicle(c *gin.Context) {
	var vehicle []model.Vehicle
	errGetVehicle := gc.DB.Raw(`
	SELECT vehicle.id, vehicle.name , vehicle.status , vehicle.type , 
	vehicle.daily_price , GROUP_CONCAT(image_vehicle.image_url) as image
	
	FROM vehicle JOIN image_vehicle ON image_vehicle.id_vehicle = vehicle.id
  
	WHERE vehicle.status = "Available"
	
	GROUP BY vehicle.id, vehicle.name , vehicle.status , vehicle.type, 
	vehicle.daily_price
	`).Scan(&vehicle).Error
	if errGetVehicle != nil {
		log.Println(errGetVehicle)

		return
	}
	
	for i := 0; i < len(vehicle); i++ {
		
		vehicle[i].ImageList = strings.Split(vehicle[i].Image, ",")
	}
	
	c.JSON(200, model.MakeRespond(vehicle, 200, "Get danh sach xe thanh cong"))
}

// lay thong tin chi tiet 1 xe
func (gc *Controller) DetailVehicle(c *gin.Context) {
	VehicleId := c.Param("id")

	var detailvehicle model.Vehicle
	
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
	c.JSON(200, model.MakeRespond(detailvehicle, 200, "Get thong tin 1 xe success"))

}
// book xe
func (gc *Controller) Booking(c *gin.Context) {
	var infoBooking model.BookingJSON
	var user model.User
	var vehicle model.Vehicle
	var wallet model.Wallet

	c.BindJSON(&infoBooking)
	log.Println(infoBooking)

	err := gc.DB.Raw("SELECT id, user_name as name FROM user WHERE id=?", infoBooking.UserID).Scan(&user).Error
	if err != nil {
		log.Println(err)
		c.JSON(400, model.MakeRespond(nil, 400, "Loi ko lay dc id user"))
		return
	}

	err = gc.DB.Raw("SELECT * FROM vehicle WHERE id=?", infoBooking.VehicleID).Scan(&vehicle).Error
	if err != nil {
		log.Println(err)
		c.JSON(400, model.MakeRespond(nil, 400, "ko lay dc id vehicle"))
		return
	}

	err = gc.DB.Raw("SELECT * FROM wallet WHERE user_id=?", user.ID).Scan(&wallet).Error
	if err != nil {
		log.Println(err)
		c.JSON(400, model.MakeRespond(nil, 400, "Loi ko lay dc wallet"))
		return
	}

	session := (infoBooking.EndDate - infoBooking.StartDate) / DayTime
	totalAmount := session * vehicle.DailyPrice

	// kiem tra vi
	if wallet.Money < totalAmount {
		c.JSON(http.StatusBadRequest, map[string]string{"message": "Ví không đủ tiền"})
		return
	}
	// không thì insert
	errPostInfoBooking := gc.DB.Exec(`
	      INSERT INTO booking (user_id,vehicle_id,total_cost,start_date,end_date) values (?,?,?,?,?)
		`, user.ID, infoBooking.VehicleID, totalAmount, infoBooking.StartDate, infoBooking.EndDate).Error
	if errPostInfoBooking != nil {
		log.Println(errPostInfoBooking)
		c.JSON(400, model.MakeRespond(nil, 400, "Loi ko insert dc vao DB"))
		return
	}

	// update so tien vi sau khi dat xe
	errUpdateWallet := gc.DB.Exec(`
		  UPDATE wallet SET money = money - ? WHERE user_id = ?
		  
		`,totalAmount, infoBooking.UserID).Error
	if errUpdateWallet != nil {
		log.Println(errUpdateWallet)
		c.JSON(400, model.MakeRespond(nil, 400, "Loi ko update đc tien trong vi sau khi thue xe"))
		return
	}

	// update status xe rented
	errUpdateStatus := gc.DB.Exec(`
	      UPDATE vehicle SET status = "Rented" WHERE ID = ?
		`, infoBooking.VehicleID).Error
	if errUpdateStatus != nil {
		log.Println(errUpdateStatus)
		c.JSON(400, model.MakeRespond(nil, 400, "Loi ko chuyen status sang Rented"))
		return
	}
	

	c.JSON(200, model.MakeRespond(map[string]interface{}{
		"user_id":      user.ID,
		"vehicle_id":   vehicle.ID,
		"startdate":    infoBooking.StartDate,
		"enddate":      infoBooking.EndDate,
		"total_amount": totalAmount,
		"money":        wallet.Money - totalAmount,
		"status":       "Rented",
	},200, "dat xe thanh cong"))

}

// api top up wallet
func (gc *Controller)TopUpWallet(c *gin.Context){
   var topupwallet model.TopUpWalletJSON
   var wallet model.Wallet
     c.BindJSON(&topupwallet)
     log.Println(topupwallet)

      errUpdateWallet := gc.DB.Exec(`
		  UPDATE wallet SET money = money + ? WHERE user_id = ?
		  
		`,topupwallet.Amount, topupwallet.UserID ).Error
	if errUpdateWallet != nil {
		log.Println(errUpdateWallet)
		c.JSON(400, model.MakeRespond(nil, 400, "Loi ko lay dc so tien nap vao"))
		return
	}

	err := gc.DB.Raw("SELECT * FROM wallet WHERE user_id=?", topupwallet.UserID).Scan(&wallet).Error
	if err != nil {
		log.Println(err)
		c.JSON(400, model.MakeRespond(nil, 400, "Loi ko update đc tien trong vi sau khi nap tien"))
		return
	}
	c.JSON(200, model.MakeRespond(map[string]interface{}{
		"user_id" : topupwallet.UserID,
		"amount ": topupwallet.Amount,
		"balance"  : wallet.Money ,
    }, 200 , "nap tien thanh cong"))
}




    

