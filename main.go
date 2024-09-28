package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"PayHandler/salebot"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type PostBack struct { //crypto cloud
	Status       string  `form:"status"`
	InvoiceID    string  `form:"invoice_id"`
	AmountCrypto float64 `form:"amount_crypto"`
	OrderID      string  `form:"order_id"`
	Currency     string  `form:"currency"`
	Token        string  `form:"token"`
}

const successStatus = "success"

func init() {
	f, err := os.OpenFile("server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.SetOutput(f)

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func main() {
	r := gin.Default()
	r.POST("/postback", PostBackHandler)
	err := r.Run()
	if err != nil {
		return
	}

}

func PostBackHandler(c *gin.Context) {
	var post PostBack

	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	log.Println(post.Token)

	ok, err := VerifyToken(post.Token)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	log.Println("Прошла провеврка статуса")

	if post.Status == successStatus {
		errChan := make(chan error, 1)
		salebot.SaleAsync(errChan, post.OrderID)

		err := <-errChan

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Неизвестный статус"})

}

func VerifyToken(tokenStr string) (bool, error) {
	secretKey := []byte(os.Getenv("CRYPTO_CLOUD_SECRET"))
	log.Println("Ключ")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		log.Println("Ключ2")
		return secretKey, nil
	})
	if err != nil {
		return false, err
	}

	log.Println("Ключ4")

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		log.Println("Ключ5")
		return claims.VerifyExpiresAt(time.Now().Unix(), true), nil

	}
	log.Println("Ключ6")

	return false, nil
}
