package auth

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
)

// SMTPLoginProcessor for login using SMTP
type SMTPLoginProcessor struct {
	Host       string
	Port       int
	Plain      bool
	TLS        bool
	SkipVerfiy bool
}

// Login implements interface of LoginProcessor.
func (s *SMTPLoginProcessor) Login(account, password string) error {
	c, err := smtp.Dial(fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return err
	}

	defer c.Close()

	if c.Hello("team"); err != nil {
		return err
	}

	if s.TLS {
		if ok, _ := c.Extension("STARTTLS"); ok {
			err = c.StartTLS(&tls.Config{
				InsecureSkipVerify: s.SkipVerfiy,
				ServerName:         s.Host,
			})

			if err != nil {
				return err
			}
		} else {
			return errors.New("SMTP server unsupports TLS")
		}
	}

	if ok, _ := c.Extension("AUTH"); !ok {
		return errors.New("SMTP server unsupport AUTH")
	}

	var auth smtp.Auth
	if s.Plain {
		auth = smtp.PlainAuth("", account, password, s.Host)
	} else {
		auth = &customSMTPAuth{account: account, password: password}
	}

	return c.Auth(auth)
}

type customSMTPAuth struct {
	account  string
	password string
}

func (s *customSMTPAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(s.account), nil
}

func (s *customSMTPAuth) Next(from []byte, more bool) ([]byte, error) {
	if more {
		switch string(from) {
		case "Username:":
			return []byte(s.account), nil
		case "Password:":
			return []byte(s.password), nil
		}
	}
	return nil, nil
}
