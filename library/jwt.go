package library

import (
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fritz-immanuel/eral-promo-library-go/configs"
	"github.com/fritz-immanuel/eral-promo-library-go/library/appcontext"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Credential struct {
	ID          string `json:"ID"`
	Name        string `json:"Name"`
	Username    string `json:"Username"`
	Email       string `json:"Email"`
	PhoneNumber string `json:"PhoneNumber"`
	Type        string `json:"Type"`
	BusinessID  string    `json:"BusinessID"`

	FsId         string `json:"fsid"`
	ClientId     string `json:"clientid"`
	ClientSecret string `json:"clientsecret"`
	RefreshToken string `json:"refreshtoken"`
}

type CredentialMobile struct {
	ID          string `json:"ID"`
	Name        string `json:"Name"`
	Username    string `json:"Username"`
	Email       string `json:"Email"`
	PhoneNumber string `json:"PhoneNumber"`
	Type        string `json:"Type"`

	FsId         string `json:"fsid"`
	ClientId     string `json:"clientid"`
	ClientSecret string `json:"clientsecret"`
	RefreshToken string `json:"refreshtoken"`
}

const JwtSalt = "secret"

func JwtSignString(c Credential) (string, error) {
	sign := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := sign.Claims.(jwt.MapClaims)

	claims["ID"] = c.ID
	claims["Name"] = c.Name
	claims["Email"] = c.Email
	claims["PhoneNumber"] = c.PhoneNumber
	claims["LoginTime"] = time.Now()
	claims["Exp"] = time.Now().Add(time.Hour * 72)
	claims["Type"] = c.Type
	claims["BusinessID"] = c.BusinessID

	config, _ := configs.GetConfiguration()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	token, err := sign.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	if errRedis := redisClient.Set(
		token,
		fmt.Sprintf("{\"id\":%s}", c.ID),
		time.Second*time.Duration(config.RedisTimeOut),
	).Err(); errRedis != nil {
		log.Printf(`
		======================================================================
		Error Storing Caching in "Auth":
		Error: %v,
		======================================================================
		`, errRedis)
		return "", errRedis
	}
	return token, nil
}

func JwtSignMobileString(c Credential) (string, error) {
	sign := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := sign.Claims.(jwt.MapClaims)

	claims["ID"] = c.ID
	claims["Name"] = c.Name
	claims["Email"] = c.Email
	claims["PhoneNumber"] = c.PhoneNumber
	claims["LoginTime"] = time.Now()
	claims["Exp"] = time.Now().Add(time.Hour * 72)
	claims["Type"] = c.Type
	claims["BusinessID"] = c.BusinessID

	config, _ := configs.GetConfiguration()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	token, err := sign.SignedString([]byte("secretmobile"))
	if err != nil {
		return "", err
	}

	if errRedis := redisClient.Set(
		token,
		fmt.Sprintf("{\"id\":%s}", c.ID),
		time.Second*time.Duration(config.RedisTimeOut),
	).Err(); errRedis != nil {
		log.Printf(`
		======================================================================
		Error Storing Caching in "Auth":
		Error: %v,
		======================================================================
		`, errRedis)
		return "", errRedis
	}
	return token, nil
}

func extractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecretString := "secret" // Value
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}

func extractMobileClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecretString := "secretmobile" // Value
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}

func GetJWTClaims(ctx *gin.Context, token string) (jwt.MapClaims, bool) {
	var claims jwt.MapClaims
	var ok bool
	if token == "" {
		JwtActiveToken := appcontext.SessionID(ctx)
		claims, ok = extractClaims(*JwtActiveToken)
	} else {
		JwtActiveToken := token
		claims, ok = extractClaims(JwtActiveToken)

	}
	return claims, ok
}

func GetJWTMobileClaims(ctx *gin.Context, token string) (jwt.MapClaims, bool) {
	var claims jwt.MapClaims
	var ok bool
	if token == "" {
		JwtActiveToken := appcontext.SessionID(ctx)
		claims, ok = extractMobileClaims(*JwtActiveToken)
	} else {
		JwtActiveToken := token
		claims, ok = extractMobileClaims(JwtActiveToken)

	}
	return claims, ok
}

func GetJWTClaimsMock() jwt.MapClaims {
	var ctx *gin.Context
	claims, _ := GetJWTClaims(ctx, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJCdXNpbmVzc0lEIjoxLCJFbWFpbCI6ImppbW15QHNlcGFyaW5kby5jb20iLCJFeHAiOjE1NzYyMzg1NTAsIklEIjoyLCJSb2xlSUQiOjEsIlRyaXBJRCI6MCwiVXNlcm5hbWUiOiJqaW1teSJ9.RvdZ6I7VTSspCnsvQflBgwrCVwUtENGu846CQqgcSh4")
	return claims
}

func SetJwtClaimsMock() {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJCdXNpbmVzc0lEIjoxLCJFbWFpbCI6ImppbW15QHNlcGFyaW5kby5jb20iLCJFeHAiOjE1NzYyMzg1NTAsIklEIjoyLCJSb2xlSUQiOjEsIlRyaXBJRCI6MCwiVXNlcm5hbWUiOiJqaW1teSJ9.RvdZ6I7VTSspCnsvQflBgwrCVwUtENGu846CQqgcSh4"
	configs.JwtActiveToken = &token
}
