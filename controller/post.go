package controller

import (
	"apinews/models"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func IndexPost(c *gin.Context) {
	reqtoken := c.Request.Header.Get("Authorization")
	s := strings.Split(reqtoken, " ")
	tkn := s[1]
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tkn, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("Access_Secret")), nil
	})
	if err == nil && token.Valid {
	}
	for i, v := range claims {
		if i == "user_id" {
			fmt.Println(v)
			user := models.User{}
			db.Where("id = ?", v).Find(&user)
			c.JSON(http.StatusOK, user)
			break
		}
	}
}
