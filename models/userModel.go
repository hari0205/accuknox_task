package models

type User struct {
	Name     string `json:"name" gorm:"name" binding:"required"`
	Email    string `json:"email" gorm:"email" binding:"required"`
	Password string `json:"password" gorm:"password" binding:"required"`
}

type Login struct {
	Email    string `json:"email" gorm:"email" binding:"required"`
	Password string `json:"password" gorm:"password" binding:"required"`
}
