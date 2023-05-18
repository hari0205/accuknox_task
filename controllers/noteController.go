package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ini "github.com/hari0205/accu-task-crud/init"
	"github.com/hari0205/accu-task-crud/models"
)

type NoteBody struct {
	Sid  string `json:"sid"`
	Note string `json:"note"`
}

type NoteBodyGet struct {
	Sid string `json:"sid"`
}

type NoteBodyResponse struct {
	Sid      string      `json:"sid"`
	NoteBody *[]NoteBody `json:"note"`
}

type NoteBodyDelete struct {
	Sid string `json:"sid"`
	Id  uint32 `json:"id"`
}

func CreateNotes(c *gin.Context) {
	var note NoteBody

	if err := c.ShouldBindJSON(&note); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Bad input. Please check your input and try again",
		})
		return
	}

	val, err := ini.Redis.Get(redisctx, note.Sid).Result()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusNotFound,
			"message": "Session not found. Please login and check your session ID",
		})
		return
	}

	insertNote := models.Notes{
		Text:   note.Note,
		Author: val,
	}
	ini.DB.Create(&insertNote).Where("email=?", val)
	// Must be status 204
	c.JSON(http.StatusOK, gin.H{
		"data": insertNote,
	})
}

func GetNotes(c *gin.Context) {
	var note NoteBodyGet
	if err := c.ShouldBindJSON(&note); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Bad input. Please check your input and try again",
		})
		return
	}
	val, err := ini.Redis.Get(redisctx, note.Sid).Result()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusNotFound,
			"message": "Session not found. Please login and check your session ID",
		})
		return
	}
	var userNotes []models.Notes
	res := ini.DB.Find(&userNotes).Where("email = ?", val)

	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No notes were found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": userNotes,
	})

}

func DeleteNotes(c *gin.Context) {

	var note NoteBodyDelete
	if err := c.ShouldBindJSON(&note); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Bad input. Please check your input and try again",
		})
		return
	}

	val, err := ini.Redis.Get(redisctx, note.Sid).Result()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Session not found. Please login and check your session ID",
		})
		return
	}

	//var notedel models.Notes
	res := ini.DB.Where(&models.Notes{Author: val, Id: note.Id}).Delete(&models.Notes{})
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No note of the ID was found.",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Note was successfully deleted.",
		})
	}

}
