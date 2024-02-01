package main

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username    string    `gorm:"size:100;index;not null;comment:使用者名稱"`
	Password    string    `gorm:"size:100;not null;comment:使用者密碼"`
	RealName    string    `gorm:"size:100;not null;comment:真實姓名"`
	Card        string    `gorm:"size:11;not null;comment:身分證"`
	Birth       time.Time `gorm:"not null;comment:生日"`
	Arrived     time.Time `gorm:"not null;comment:到職日"`
	Leaved      time.Time `gorm:"comment:離職日"`
	Status      bool      `gorm:"default:true;not null;comment:狀態"`
	Remark      string    `gorm:"comment:備註"`
	Desc        string    `gorm:"comment:簡介"`
	Email       string    `gorm:"size:100;comment:信箱"`
	Mobile      string    `gorm:"size:11;comment:手機號"`
	Avatar      string    `gorm:"comment:頭像"`
	LastMover   uint      `gorm:"default:0;not null;comment:最後異動者"`
	CheckInCode string    `gorm:"size:50;comment:卡號"`
	Backend     bool      `gorm:"default:false;not null;comment:是否可登入後台"`
}

func (u *User) TableName() string {
	return "users"
}
