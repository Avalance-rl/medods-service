package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"medods-service/internal/cache"
	"medods-service/internal/config"
	"medods-service/internal/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)


const (
	AccessLife  = 3600 * 24
	RefreshLife = 3600 * 24
)

type request struct {
	Guid string `json:"guid"`
}


func CreateUserTokens(c *gin.Context) {
	body := request{}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	if _, err := uuid.Parse(body.Guid); err != nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{
            "error": "GUID entered incorrectly",
        })
        return
    }
	accessToken, refreshToken, rawRefreshToken, err := newTokens(body, c.ClientIP())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid to create token",
		})
		return
	}

	err = cache.RDB.Set(body.Guid+":access_token", accessToken, AccessLife * time.Second).Err()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to save access token in cache",
        })
        return
    }

	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		 c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to bcrypt refresh token",
       })
        return	
	}

	err = database.InsertToken(string(hashedToken), rawRefreshToken, body.Guid)
	if err != nil {
        c.JSON(http.StatusConflict, gin.H{
            "error": "Conflict, you have already registared",
        })
        return
    }
	jwtCookie := http.Cookie{
		Name:     "refreshToken",
		Value: rawRefreshToken,
		MaxAge: RefreshLife,
		Path:     "/api/auth",
		HttpOnly: true,
	}
	c.SetCookie(
		jwtCookie.Name,
		jwtCookie.Value,
		jwtCookie.MaxAge,
		jwtCookie.Path,
		jwtCookie.Domain,
		jwtCookie.Secure,
		jwtCookie.HttpOnly,
	)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})

}

func newTokens(req request, clientIP string) (string, string, string, error){
	jwtSecretKey := []byte(config.CFG.JWTSecretKey)
	rawRefreshToken, err := generateRandomString(32)
	if err != nil {
		return "", "", "", err
	}
	accessPayload := jwt.MapClaims{
		"guid": req.Guid,
		"exp": AccessLife,
		"client_ip": clientIP,
		"refresh_token": rawRefreshToken,
		"time": time.Now(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessPayload)
	signedAccessToken, err := accessToken.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", "", err
	}

	
	refreshToken, err := generateRandomString(10)

	if err != nil {
		return "", "", "", err
	}

	return signedAccessToken, refreshToken, rawRefreshToken, nil
}

func generateRandomString(length int) (string, error) {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}



