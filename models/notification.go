package models

import "time"

type MobileNotification struct {
	Id              int       `json:"id" gorm:"primary_key"`
	Message         string    `json:"message"`
	UserId          int       `json:"user_id"`
	DataPengajuanId int       `json:"data_pengajuan_id"`
	DataPengajuan   Pengajuan `json:"data_pengajuan" gorm:"foreignKey:DataPengajuanId;references:Id"`
	CreatedAt       time.Time
}

type MobileNotificationCreateRequest struct {
	Id              int    `json:"id" gorm:"primary_key"`
	Message         string `json:"message"`
	UserId          int    `json:"user_id"`
	DataPengajuanId int    `json:"data_pengajuan_id"`
	StatusPengajuan string `json:"status_pengajuan"`
}
