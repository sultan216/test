package routes

import (
	"apinews/config"
	"apinews/controller"
	"apinews/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var r = gin.Default()
var api = r.Group("/api/v1")
var db = config.ConnectDB()
var tkn string

type acc struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
	FullName string `json:"fullname"`
}

// Login Jwt Token
func Login(c *gin.Context) {
	user := &models.User{}
	var account acc
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
	}
	if err := db.Where("username = ?", account.Username).First(user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, "Not Found")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(account.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		c.JSON(http.StatusUnauthorized, "Invalid Login Credentials, please try again")
	}
	db.Where("username = ?", account.Username).Find(user)
	token, err := CreateToken(uint64(user.ID))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"Token": token,
	})
}

func CreateToken(userid uint64) (string, error) {
	var err error
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userid
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(os.Getenv("Access_Secret")))
	if err != nil {
		return "", err
	}
	tkn = token
	return tkn, nil
}

// Route & Middleware
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// manual Set bearer token in Authorization header
		reqtoken := c.Request.Header.Get("Authorization")

		splitoken := strings.Split(reqtoken, "Bearer")
		errx := len(splitoken)
		if errx != 2 {
			c.JSON(http.StatusForbidden, "Please input A token in here")
		}
		reqtoken = strings.TrimSpace(splitoken[1])
		token, err := jwt.Parse(reqtoken, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("Access_Secret")), nil
		})
		if err == nil && token.Valid {
			c.Set("token_Jwt", reqtoken)
			c.Next()
		} else {
			c.JSON(http.StatusBadRequest, "Token Invalid")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Auto Set Bearer Token in AuAuthorization Header,
		// var p = tkn
		// if tkn == "" {
		// 	c.JSON(http.StatusBadRequest, "You Need Bearer Token to access this API")
		// 	c.AbortWithStatus(http.StatusTemporaryRedirect)
		// 	return
		// }
		// token, err := jwt.Parse(tkn, func(t *jwt.Token) (interface{}, error) {
		// 	return []byte(os.Getenv("Access_Secret")), nil
		// })
		// if err == nil && token.Valid {
		// 	barer := "Bearer " + p
		// 	c.Request.Header.Set("Authorization", barer)
		// 	c.Request.Header.Add("Accept", "application/json")
		// 	c.Set("token_Jwt", p)
		// 	c.Next()
		// } else {
		// 	c.JSON(http.StatusBadRequest, "Please Input a Valid Token in Here")
		// 	return
		// }

	}
}
func MainRoute() {
	AuthRoute()
	PostRoute()
	r.Run()
}
func PostRoute() {
	p := api.Group("post")
	p.Use(Middleware())
	p.GET("/", controller.IndexPost)

}

func AuthRoute() {
	auth := api.Group("/auth")
	auth.GET("/login", controller.Login)
	auth.POST("/login", Login)
	auth.GET("/register", controller.Register)
	auth.POST("/register", controller.Register)
	auth.GET("/logout", func(c *gin.Context) {
		tkn = ""
	})
}
