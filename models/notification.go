package models

type MobileNotification struct {
	Id              int       `json:"id" gorm:"primary_key"`
	Message         string    `json:"message"`
	DataPengajuanId int       `json:"data_pengajuan_id"`
	DataPengajuan   Pengajuan `json:"data_pengajuan" gorm:"foreignKey:DataPengajuanId;references:Id"`
}

type MobileNotificationCreateRequest struct {
	Id              int    `json:"id" gorm:"primary_key"`
	Message         string `json:"message"`
	DataPengajuanId int    `json:"data_pengajuan_id"`
}

