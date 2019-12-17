package main

import (
	"fmt"
	"github.com/zlh2015/go-anyway/email"
	"net"
	// "github.com/gin-gonic/gin"
)

func test() (err error) {
	conn, err := net.Dial("tcp", "pop-mail.outlook.com:995")
	if err != nil {
		fmt.Println(err)
	}
	data := make([]byte, 1000)
	conn.Write([]byte(""))
	_, err = conn.Read(data)
	fmt.Println(string(data))
	conn.Write([]byte("USER jack.hg2018@outlook.com\r\n"))
	_, err = conn.Read(data)
	fmt.Println(string(data))
	return err
}

func pop() (err error) {
	var addr = "pop-mail.outlook.com:955"
	client, err := email.Dial(addr)
	fmt.Print(11)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = client.User("jack.hg2018@outlook.com")
	if err != nil {
		fmt.Println(err)
		return err
	}
	return

}

func main() {
	fmt.Println("hello")
	// test()
	pop()
	// r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	// r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
