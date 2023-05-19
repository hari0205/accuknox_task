package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ini "github.com/hari0205/accu-task-crud/init"
	"github.com/hari0205/accu-task-crud/models"
	"golang.org/x/crypto/bcrypt"
)

var redisctx = context.Background()

func SignUp(c *gin.Context) {

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println(fmt.Errorf("unable to bind JSON: %s", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Bad input. Please check your input and try again",
		})
		return

	}

	// Check to see if the user already exists; if not, create a new user else send an error user already exists

	res := ini.DB.Where("email = ?", user.Email).First(&user)
	if res.RowsAffected != 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "User already exists",
		})
		return
	}

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Unable to process . Please try again",
		})
		return
	}

	// Return success
	newuser := models.User{Name: user.Name, Email: user.Email, Password: string(hashedpassword)}
	res = ini.DB.Create(&newuser)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating user ",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "User created successfully",
		})
	}

}

func Login(c *gin.Context) {

	var LoginUser *models.Login
	if err := c.ShouldBindJSON(&LoginUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Bad input. Please check your input and try again",
		})
	}

	// Check if Email exist

	var user models.User
	ini.DB.First(&user, "email= ? ", LoginUser.Email)
	if user.Email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "User does not exist or Invalid Credentials",
		})
		return
	}
	// Password check
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(LoginUser.Password))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Credentials. Check you credentials and try again",
		})
		return
	}

	// If session for an user exists, clear session ( One user must only have one session)
	keys, _ := ini.Redis.Keys(redisctx, "*").Result()
	for _, k := range keys {
		fmt.Println(k)
		val, _ := ini.Redis.Get(redisctx, k).Result()
		if val == LoginUser.Email {
			ini.Redis.Del(redisctx, k)
		}
	}
	// Set session ID in redis
	id := uuid.New().String()
	err = ini.Redis.Set(redisctx, id, user.Email, 0).Err()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":       http.StatusInternalServerError,
			"errorMessage": "Unable to generate session",
		})

	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Successfully logged in",
		"session-id": id,
	})

}
