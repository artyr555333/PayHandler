package main

import (
	"log"
	"net/http"

	"PayHandler/salebot"
	"github.com/gin-gonic/gin"
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

func main() {
	r := gin.Default()
	r.POST("/postback", Test)
	err := r.Run()
	if err != nil {
		return
	}

}

func Test(c *gin.Context) {
	var post PostBack

	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	log.Println(post.Token)

	//res := VerifyToken(post.Token)

	//if !res {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
	//	return
	//}

	if post.Status != successStatus {
		c.JSON(http.StatusBadRequest, gin.H{"error": post.Status})
	}

	log.Println("Прошла провеврка статуса")

	if post.Status == successStatus {
		log.Println("Статус успешный")
		log.Println("Начало выполнения запроса в салебот")
		errChan := make(chan error, 1)
		log.Println("Канал создан")
		salebot.SaleAsync(errChan, post.OrderID)

		log.Println("SaleAsync выполнена")
		err := <-errChan
		log.Println("Ошибка из канала записана в переменную err")

		log.Println("Окончание выполнения запроса в салебот")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

//func VerifyToken(tokenStr string) bool {
//	secretKey := []byte("gZashj4XJJJRrisI33QgDkMARZVTMPkkvooI")
//
//	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
//
//		_, ok := token.Method.(*jwt.SigningMethodHMAC)
//
//		if !ok {
//			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//		}
//
//		return secretKey, nil
//	})
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	claims, ok := token.Claims.(jwt.MapClaims)
//
//	if ok && token.Valid {
//		return claims.VerifyExpiresAt(time.Now().Unix(), true)
//
//	}
//
//	return false
//}
