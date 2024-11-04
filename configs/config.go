package configs

import (
	"os"
	"strconv"
)

const (
	androidAppMinimumVersion = "ANDROID_APP_MINIMUM_VERSION"
	iosAppMinimumVersion     = "IOS_APP_MINIMUM_VERSION"

	externalURL         = "EXTERNAL_URL"
	externalToken       = "EXTERNAL_TOKEN"
	externalAccessToken = "EXTERNAL_ACCESS_TOKEN"

	appUrl             = "APP_URL"
	serverName         = "SERVER_NAME"
	dbConnectionString = "DB_CONNECTION_STRING"

	portApps = "PORT_APPS"

	redisAddr     = "REDIS_ADDR"
	redisDB       = "REDIS_DB"
	redisPassword = "REDIS_PASSWORD"
	redisTimeOut  = "REDIS_TIME_OUT"

	sendWhatsappAPI   = "SEND_WHATSAPP_API"
	sendWhatsappToken = "SEND_WHATSAPP_TOKEN"

	teleBotToken = "TELE_BOT_TOKEN"
	teleGroupID  = "TELE_GROUP_ID"

	firebaseServerKey        = "FIREBASE_SERVER_KEY"
	firebaseSenderID         = "FIREBASE_SENDER_ID"
	firebaseStorageBucketURL = "FIREBASE_BUCKET_URL"
	firebaseAuthFilePath     = "FIREBASE_AUTH_FILE_PATH"

	whitelistedIps = "WHITELISTED_IPS"
)

// TODO check mana yg masih dipakai
var (
	JwtActiveToken *string
)

// Config contains application configuration
type Config struct {
	// Minimum App versions
	AndroidAppMinimumVersion string
	IosAppMinimumVersion     string

	ExternalURL         string
	ExternalToken       string
	ExternalAccessToken string

	// DB
	DBConnectionString string

	// Misc
	PortApps       string
	ServerName     string
	WhitelistedIps string

	// Redis
	RedisAddr     string
	RedisDB       int
	RedisPassword string
	RedisTimeOut  int

	// WA
	SendWhatsappAPI   string
	SendWhatsappToken string

	// Telegram
	TeleBotToken string
	TeleGroupID  string

	// Firebase
	FirebaseServerKey        string
	FirebaseSenderID         string
	FirebaseStorageBucketURL string
	FirebaseAuthFilePath     string

	AppURL string
}

var config *Config

func getEnvOrDefault(env string, defaultVal string) string {
	e := os.Getenv(env)
	if e == "" {
		return defaultVal
	}

	return e
}

// GetConfiguration , get application configuration based on set environment
func GetConfiguration() (*Config, error) {
	if config != nil {
		return config, nil
	}

	config := &Config{}

	config.AndroidAppMinimumVersion = os.Getenv(androidAppMinimumVersion)
	config.IosAppMinimumVersion = os.Getenv(iosAppMinimumVersion)

	config.ExternalURL = os.Getenv(externalURL)
	config.ExternalToken = os.Getenv(externalToken)
	config.ExternalAccessToken = os.Getenv(externalAccessToken)

	config.DBConnectionString = os.Getenv(dbConnectionString)

	config.PortApps = os.Getenv(portApps)

	config.RedisAddr = os.Getenv(redisAddr)
	config.RedisDB, _ = strconv.Atoi(os.Getenv(redisDB))
	config.RedisPassword = os.Getenv(redisPassword)
	config.RedisTimeOut, _ = strconv.Atoi(os.Getenv(redisTimeOut))

	config.SendWhatsappAPI = os.Getenv(sendWhatsappAPI)
	config.SendWhatsappToken = os.Getenv(sendWhatsappToken)

	// TELEGRAM
	config.TeleBotToken = os.Getenv(teleBotToken)
	config.TeleGroupID = os.Getenv(teleGroupID)

	config.FirebaseServerKey = os.Getenv(firebaseServerKey)
	config.FirebaseAuthFilePath = os.Getenv(firebaseAuthFilePath)
	config.FirebaseStorageBucketURL = os.Getenv(firebaseStorageBucketURL)
	config.FirebaseSenderID = os.Getenv(firebaseSenderID)

	config.AppURL = os.Getenv(appUrl)
	config.WhitelistedIps = os.Getenv(whitelistedIps)
	config.ServerName = os.Getenv(serverName)

	return config, nil
}
