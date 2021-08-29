package dao

import "gorm.io/gorm"

//用户表
type Users struct {
	gorm.Model
	Account  string `gorm:"type:varchar(32);not null;unique" json:"account"`
	Pwd      string `gorm:"type:varchar(32);not null" json:"pwd"`
	NickName string `gorm:"type:varchar(32)" json:"nick_name"`
	Icon     string `json:"icon"`
}
