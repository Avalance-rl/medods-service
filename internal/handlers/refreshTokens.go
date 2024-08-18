package handlers

import (
	"fmt"
	"log/slog"
	"medods-service/internal/cache"
	"medods-service/internal/config"
	"medods-service/internal/database"
	"net/http"
	"net/smtp"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	smtpServer   = "smtp.example.com"
	smtpPort     = "587"
	smtpUser     = "some-email@example.com"
	smtpPassword = "some-password"
)

func RefreshTokens(c *gin.Context) {
	rawRefreshToken, _ := c.Cookie("refreshToken")
	userEmail := c.GetHeader("Email")
	if !isEmailValid(userEmail)  {
		slog.Error("Email is not valid")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email is not valid",
		})
		return	
	}
	token, err := database.GetToken(rawRefreshToken)
	if err != nil {
		slog.Error(fmt.Sprintf("RefreshToken() error = %v, can't find value from DB", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server",
		})
		return	
	}
	cacheKey := token.UserGUID+":access_token"
 	tokenInCache, err := cache.RDB.Get(cacheKey).Result()
	if err != nil {
		slog.Error(fmt.Sprintf("RefreshToken() error = %v, invalid cookies redis", err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": token.UserGUID,
		})
		return
	}
	claimsMap, err := decodeAndVerifyJWT(tokenInCache, []byte(config.CFG.JWTSecretKey))
	if err != nil {
		slog.Error(fmt.Sprintf("Decode token error = %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}
	
	if guid, ok := claimsMap["guid"].(string); ok {
        if guid != token.UserGUID {
			slog.Error(fmt.Sprintf("Token is not valid = %v", err))
			c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Your token is not valid",
			})
			return	
        }
    }

	if clientIP, ok := claimsMap["client_ip"].(string); ok {
        if clientIP != c.ClientIP() {
			subject := "Warning: IP Address Changed"
			body := fmt.Sprintf("Your token was used from a different IP address: %s. If this wasn't you, please check your account for security issues.", c.ClientIP())
			err = sendEmail(userEmail, subject, body)
			slog.Error(fmt.Sprintf("Token is not valid: IP address changed to %s", err))
			c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Your token is not valid",
			})
			return	
        }
    }
	
	if refresh_token, ok := claimsMap["refresh_token"].(string); ok {
        if rawRefreshToken != refresh_token {
			slog.Error(fmt.Sprintf("Token is not valid = %v", err))
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "Your token is not valid",
			})
			return	
        }
    }	
	
	if time.Since(token.CreatedAt) > time.Duration(RefreshLife) * time.Second {
		err = database.DeleteToken(rawRefreshToken)
		if err != nil {
        	c.JSON(http.StatusInternalServerError, gin.H{
         	   "error": "Intrenal server error",
        	})
        	return
    	}	
		err = database.InsertToken(token.TokenHash, token.RawRefreshToken, token.UserGUID)
		if err != nil {
        	c.JSON(http.StatusInternalServerError, gin.H{
         	   "error": "Intrenal server error",
        	})
        	return
    	}
		err = cache.RDB.Set(token.UserGUID+":access_token", tokenInCache, AccessLife * time.Second).Err()
		if err != nil {
        	c.JSON(http.StatusInternalServerError, gin.H{
         	   "error": "Intrenal server error",
        	})
        	return
    	}
	}
	c.JSON(http.StatusOK, gin.H{
        "message": "Your refresh token and access token update",
    })
	
}


func decodeAndVerifyJWT(tokenString string, jwtSecretKey []byte) (map[string]interface{}, error) {

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtSecretKey, nil
    })
	
	claimsMap := make(map[string]interface{})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
        for key, value := range claims {
			claimsMap[key] = value
        }
    } else {
        slog.Error("Invalid token claims")
    }


	return claimsMap, nil
}


func sendEmail(to, subject, body string) error {
	from := smtpUser
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpServer)

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, from, []string{to}, msg)
	return err
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}