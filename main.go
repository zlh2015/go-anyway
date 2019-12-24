package main

import (
	"crypto/tls"
	"fmt"
	imaplib "github.com/emersion/go-imap"
	imapcli "github.com/emersion/go-imap/client"
	"go-anyway/email/pop3"
	smtplibext "go-anyway/email/smtp"
	smtplib "net/smtp"
	// "github.com/gin-gonic/gin"
)

func test() (err error) {
	conn, err := tls.Dial("tcp", "pop-mail.outlook.com:995", nil)
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
	addr := "pop-mail.outlook.com:995"
	client, err := pop3.DialTLS(addr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = client.User("jack.hg2018@outlook.com")
	err = client.Pass("Gzhg2018")
	count, size, err := client.Stat()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(count, size)
	size, err = client.List(1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(size)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return

}

func imap() (err error) {
	addr := "imap-mail.outlook.com:993"
	client, err := imapcli.DialTLS(addr, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer client.LoggedOut()

	// Login
	if err := client.Login("jack.hg2018@outlook.com", "***"); err != nil {
		fmt.Println(err)
	}

	fmt.Println("login success!")

	// List mailboxes
	mailboxes := make(chan *imaplib.MailboxInfo, 20)
	done := make(chan error, 1)
	go func() {
		done <- client.List("", "*", mailboxes)
	}()

	fmt.Println("Mailboxes:")
	for m := range mailboxes {
		fmt.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		fmt.Println(err)
	}
	return
}

func smtp() (err error) {
	host := "smtp.office365.com"
	// host := "52.98.77.98"
	// au := smtplib.PlainAuth("", "jack.hg208@outlook.com", "Gzhg2018", host)
	au := smtplibext.LoginAuth("jack.hg2018@outlook.com", "Gzhg2018", host)
	fmt.Println(au)
	client, err := smtplib.Dial(host + ":25")
	if err != nil {
		return err
	}
	if err = client.Hello("LAPTOP-CML0ECA3"); err != nil {
		return err
	}
	if ok, _ := client.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: host, InsecureSkipVerify: false}
		if err = client.StartTLS(config); err != nil {
			return err
		}
		if err = client.Auth(au); err != nil {
			return err
		}
	}
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return nil
	}
}

func main() {
	fmt.Println("hello")
	// test()
	// pop()
	// imap()
	smtp()
	// r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	// r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
