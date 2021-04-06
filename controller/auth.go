package controller

import (
	"apinews/config"
	"apinews/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var db = config.ConnectDB()

type acc struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
	FullName string `json:"fullname"`
}

func Login(c *gin.Context) {
	c.JSON(http.StatusFound, "Welcome, u must login before access main api, if u don't have account, please create first")

}
func Register(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.JSON(http.StatusFound, "Hey, u want make account in this api?")
	} else if c.Request.Method == "POST" {
		var account acc
		if err := c.ShouldBindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}
		pass, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}
		str := string(pass)
		p := models.User{
			Username: account.Username,
			Password: str,
			FullName: account.FullName,
			Email:    account.Email,
		}
		db.Create(&p)
		c.JSON(http.StatusOK, gin.H{
			"Message": "Account Created",
			"Status": gin.H{
				"code": http.StatusOK,
			},
		})

	}
}
