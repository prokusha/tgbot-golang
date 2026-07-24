package ENV

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type TelegramEnv struct {
	Token   string
	User_ID int64
}

type DBEnv struct {
	User     string
	Password string
	Database string
}

var Telegram TelegramEnv
var DB DBEnv

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	Telegram.Token = os.Getenv("TOKEN")
	if Telegram.Token == "" {
		panic("TOKEN environment variable is empty")
	}

	user_id_env := os.Getenv("USER_ID")
	user_id_int, _ := strconv.Atoi(user_id_env)
	Telegram.User_ID = int64(user_id_int)

	DB = DBEnv{User: os.Getenv("P_USER"), Password: os.Getenv("P_PASSWORD"), Database: os.Getenv("P_DB")}
	if DB.User == "" || DB.Password == "" || DB.Database == "" {
		panic("DB environment variable is empty")
	}
}
