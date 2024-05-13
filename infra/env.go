package infra

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type ctxKey string

const CtxServiceKey ctxKey = "serviceKey"
const CtxPostInteractorKey ctxKey = "postInteractorKey"
const CtxPostIdKey ctxKey = "postIdKey"
const CtxPostOwnerIdKey ctxKey = "ownerIdKey"

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
	//goDotEnvVariable(".env")

	curDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	var envLoc string
	if strings.HasSuffix(curDir, "/cmd") {
		trimLength := len(curDir) - len("/cmd")
		envLoc = string(curDir[0:trimLength])
	} else {
		envLoc = curDir
	}
	loadErr := godotenv.Load(envLoc + "/.env")
	if loadErr != nil {
		log.Fatalln("can't load env file from current directory: " + curDir + "/.env")
	}
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

func RespondJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	b, _ := json.Marshal(data)

	w.Write(b)
}
