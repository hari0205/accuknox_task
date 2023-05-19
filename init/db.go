package ini

import (
	"fmt"
	"os"

	"github.com/hari0205/accu-task-crud/models"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var Redis *redis.Client

func ConnectToDB() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	//dburl := "postgres://vexndcof:8c1px5ukjsn26AGAuFzrYY_FegHQZFFs@tiny.db.elephantsql.com/vexndcof"
	//fmt.Println(os.Getenv("DATABASE_URL"))
	dsn := os.Getenv("DATABASE_URL")
	Db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	Db.AutoMigrate(&models.User{}, &models.Notes{}, models.Login{})
	fmt.Println("POSTGRES DB connection established")
	DB = Db

}

func ConnectToRedis() {
	//opt, _ := redis.ParseURL("redis://default:12452a1cc01a41eea9cf92653be6488e@apn1-finer-turtle-35435.upstash.io:35435")
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,
	})
	fmt.Println("REDIS connection established")
	Redis = client

}
