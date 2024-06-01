package security

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/utils"
)

var (
	privateKey   = []byte(os.Getenv("TOKEN_PRIVATE_KEY"))
	BEARER_TOKEN = "Bearer"
)

func InitJWT(path string) {

	if path != "" {
		fmt.Println("Load file in: ", path)
		err := godotenv.Load(path)
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	} else {
		fmt.Println("Load file in: ./pkg/security/.env")
		err := godotenv.Load("./pkg/security/.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func GenerateJWTToken(userId string, roles []string) ([]string, error) {

	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	refresh_tokenTTL, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TTL"))
	refreshToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"roles":  roles,
		"exp":    time.Now().Add(7 * time.Second * time.Duration(refresh_tokenTTL)).Unix(),
	}).SignedString(privateKey)
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"roles":  roles,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	}).SignedString(privateKey)

	return []string{token, refreshToken}, nil
}

func CheckRole(req *gin.Context, roles []string) error {
	mapClaims, err := parseClaims(req)
	if err != nil {
		return err
	}
	var claimRoles []string
	if roles, ok := mapClaims["roles"].([]interface{}); ok {
		for _, role := range roles {
			if roleStr, ok := role.(string); ok {
				claimRoles = append(claimRoles, roleStr)
			} else {
				log.Println("Role is not a string")
			}
		}
	} else {
		log.Println("Roles claim is not a []interface{}")
	}

	for _, role := range claimRoles {
		for _, userRole := range roles {
			if role == userRole {
				v, ok := mapClaims["exp"]
				if err = IsTimeExpired(req, v.(float64)); !ok && err != nil {
					return errors.New("need to renew token")
				}
				return nil
			}
		}
	}
	utils.WrapperResponse(req, http.StatusForbidden, "403 Forbidden")
	return errors.New("Forbidden")
}

func GetRole(req *gin.Context) ([]string, error) {
	mapClaims, err := parseClaims(req)
	if err != nil {
		return make([]string, 0), err
	}
	var claimRoles []string
	if roles, ok := mapClaims["roles"].([]interface{}); ok {
		for _, role := range roles {
			if roleStr, ok := role.(string); ok {
				claimRoles = append(claimRoles, roleStr)
			} else {
				log.Println("Role is not a string")
			}
		}
	} else {
		log.Println("Roles claim is not a []interface{}")
	}
	return claimRoles, nil
}

func GetUserId(req *gin.Context) (string, error) {
	mapClaims, err := parseClaims(req)
	if err != nil {
		return "", err
	}
	return mapClaims["userId"].(string), nil
}

func RenewToken(req *gin.Context) ([]string, error) {

	mapClaims, err := parseClaims(req)
	if err != nil {
		return []string{}, err
	}
	var claimRoles []string
	if roles, ok := mapClaims["roles"].([]interface{}); ok {
		for _, role := range roles {
			if roleStr, ok := role.(string); ok {
				claimRoles = append(claimRoles, roleStr)
			} else {
				log.Println("Role is not a string")
			}
		}
	} else {
		log.Println("Roles claim is not a []interface{}")
	}
	return GenerateJWTToken(mapClaims["userId"].(string), claimRoles)
}

func IsTimeExpired(req *gin.Context, expire float64) error {

	if expire != 0 {
		exp := time.Unix(int64(expire), 0)
		if exp.After(time.Now()) {
			return nil
		}
	} else {
		claims, err := parseClaims(req)
		if err != nil {
			return errors.New("parse token fail")
		}
		exp := time.Unix(claims["exp"].(int64), 0)
		if exp.After(time.Now()) {
			return nil
		}
	}
	return errors.New("token is expired")
}

func parseClaims(req *gin.Context) (jwt.MapClaims, error) {

	tokenString, ok := getTokenFromRequestHeader(req)
	if ok != nil {
		return nil, ok
	}
	token, err := parseToken(tokenString)
	if token != nil {
		return token.Claims.(jwt.MapClaims), err
	}
	return jwt.MapClaims{}, errors.New("cannot parse token")
}

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(tokenString *jwt.Token) (interface{}, error) {
		if _, ok := tokenString.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tokenString.Header["alg"])
		}
		// return signature
		return privateKey, nil
	})
}

func getTokenFromRequestHeader(context *gin.Context) (string, error) {
	bearerToken := context.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 && splitToken[0] == BEARER_TOKEN {
		return splitToken[1], nil
	}
	return "", errors.New("request missed header 'Authorization' Bearer")
}
