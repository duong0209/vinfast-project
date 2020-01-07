package controller

import (
	"vinfast-project/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Controller struct {
	// DB instance
	DB *gorm.DB

	// Cấu hình config
	Config model.Config
}

func NewController() *Controller {
	var c Controller
	return &c
}
