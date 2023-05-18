package main

import (
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hari0205/accu-task-crud/controllers"
	ini "github.com/hari0205/accu-task-crud/init"

	log "github.com/sirupsen/logrus"
)

func Init() {

	// Initializing logger
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)

	// Initialize postgres connection
	ini.ConnectToDB()

	// Initializing Redis connection
	ini.ConnectToRedis()

}

func main() {
	// Initializes necesary connections
	Init()

	router := gin.Default()
	router.Use(cors.Default())

	apiGrp := router.Group("/api")
	{
		// TEST ROUTE
		apiGrp.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "PONG")
		})

		// USER ROUTES
		apiGrp.POST("/signup", controllers.SignUp)
		apiGrp.POST("/login", controllers.Login)

		// NOTE ROUTES
		apiGrp.POST("/notes", controllers.CreateNotes)
		apiGrp.GET("/notes", controllers.GetNotes)
		apiGrp.DELETE("/notes", controllers.DeleteNotes)
	}

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":  http.StatusNotFound,
			"error": "Please check the endpoint and try again.",
		})
	})

	router.NoMethod(func(ctx *gin.Context) {
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{
			"code":  http.StatusMethodNotAllowed,
			"error": "Please check the HTTP method and try again.",
		})
	})
	router.Run(":8080") // Firewall warning bypass

}
