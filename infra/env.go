package infra

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env        string
	MaxRetries int
	DbHost     string
	DbPort     string
	DbName     string
	DbUser     string
	DbPassword string
}

var config *Config

func GetConfig() *Config {
	goDotEnvVariable(".env")

	config = &Config{
		Env:        GetString("Environment", "local"),
		MaxRetries: GetInt("MaxRetries", 3),
		DbHost:     GetString("DbHost", ""),
		DbPort:     GetString("DbPort", ""),
		DbName:     GetString("DbName", ""),
		DbUser:     GetString("DbUser", ""),
		DbPassword: GetString("DbPassword", ""),
	}

	log.Println(`Config loaded with db host: ` + config.DbHost)

	return config
}

func GetString(key string, defaultVal string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return defaultVal
	}
	return val
}

func GetInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if len(val) == 0 {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
