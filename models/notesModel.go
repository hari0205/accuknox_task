package models

type Notes struct {
	Id     uint32 `json:"id" gorm:"primaryKey"`
	Text   string `json:"text" binding:"required" gorm:"NOT NULL"`
	Author string `json:"-" gorm:"-"`
}
