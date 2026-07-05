package TGBOT_ENV

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Token   string
	User_ID int64
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	Token = os.Getenv("TOKEN")
	if Token == "" {
		panic("TOKEN environment variable is empty")
	}

	user_id_env := os.Getenv("USER_ID")
	user_id_int, _ := strconv.Atoi(user_id_env)
	User_ID = int64(user_id_int)
}
