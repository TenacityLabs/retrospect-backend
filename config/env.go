package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	GoEnv                  string
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	JWTExpirationInSeconds int64
	JWTSecret              string
	GCSBucketName          string
	GmailAppPassword       string
	AdminAPIKey            string
	TwilioAccountSID 	   string
	TwilioAuthToken		   string
	TwilioVerifyServiceSID string
}

// create global variable so that env isn't reinitialized every time it's called
var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		GoEnv:                  getEnv("GO_ENV", "development"),
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", "8080"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnv("DB_PASSWORD", "Password1"),
		DBAddress:              fmt.Sprintf("%s:%s", getEnv("DB_ADDRESS", "localhost"), getEnv("DB_PORT", "3306")),
		DBName:                 getEnv("DB_NAME", "retrospect"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXP", 3600*24*7),
		JWTSecret:              getEnv("JWT_SECRET", "sneakysneaky"),
		GCSBucketName:          getEnv("BUCKET_NAME", "retrospect_file_bucket"),
		GmailAppPassword:       getEnv("GMAIL_APP_PASSWORD", ""),
		AdminAPIKey:            getEnv("ADMIN_API_KEY", "spartan"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		integerValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return integerValue
	}
	return fallback
}
