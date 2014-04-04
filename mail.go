package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/smtp"
	"os"
)

type Mail struct {
	From string
	To   string
	Body []byte
}

type Server struct {
	Address       string
	SkipSSLVerify bool
	User          string
	Password      string
}

func main() {
	var (
		from      = flag.String("from", "root@localhost", "From e-mail")
		to        = flag.String("to", "", "To e-mail")
		subject   = flag.String("subject", "", "E-mail Subject")
		server    = flag.String("smtp", "", "SMTP Server")
		user      = flag.String("user", "", "SMTP User")
		pass      = flag.String("pass", "", "SMTP Pass")
		sslverify = flag.Bool("skip-ssl-verify", false, "Skip Verifcation of SSL Cert")
	)
	flag.Parse()

	stdin, _ := ioutil.ReadAll(os.Stdin)
	body := "From:" + *from + "\r\n" + "Subject:" + *subject + "\r\n" + string(stdin)

	mail := Mail{
		From: *from,
		To:   *to,
		Body: []byte(body),
	}

	srv := Server{
		Address:       *server,
		SkipSSLVerify: *sslverify,
		User:          *user,
		Password:      *pass,
	}

	if err := sendMail(mail, srv); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Sent Mail")
		os.Exit(0)
	}
}

func sendMail(mail Mail, server Server) error {
	host, _, _ := net.SplitHostPort(server.Address)
	auth := smtp.PlainAuth(
		"",
		server.User,
		server.Password,
		host,
	)

	client, err := smtp.Dial(server.Address)
	if err != nil {
		fmt.Errorf("Error clientecting", err)
	}
	defer client.Quit()

	tlc := &tls.Config{
		InsecureSkipVerify: server.SkipSSLVerify,
		ServerName:         host,
	}

	if err = client.StartTLS(tlc); err != nil {
		fmt.Errorf("SSL error:", err)
	}

	if ok, _ := client.Extension("AUTH"); ok {
		if err = client.Auth(auth); err != nil {
			fmt.Errorf("Auth Error:", err)
		}
	}
	if err = client.Mail(mail.From); err != nil {
		fmt.Errorf("Error Setting From:", err)
	}

	if err = client.Rcpt(mail.To); err != nil {
		fmt.Errorf("Error Setting To:", err)
	}

	w, err := client.Data()
	if err != nil {
		fmt.Errorf("Error connecting to server:", err)
	}

	_, err = w.Write(mail.Body)
	if err != nil {
		fmt.Errorf("Error writing to server:", err)
	}

	err = w.Close()
	if err != nil {
		fmt.Errorf("Error closing stream:", err)
	}
	return nil
}
