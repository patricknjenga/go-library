package library

import (
	"fmt"
	"net"
	"net/smtp"
)

type Recipient struct {
	Email string
}

type Mail struct {
	Message    []byte
	Recipients []Recipient
}

type Smtp struct {
	Host     string
	Password string
	Port     string
	User     string
}

func (s Smtp) Send(m Mail) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", s.Host, s.Port))
	if err != nil {
		return err
	}
	c, err := smtp.NewClient(conn, s.Host)
	if err != nil {
		return err
	}
	err = c.Mail(s.User)
	if err != nil {
		return err
	}
	for _, v := range m.Recipients {
		err = c.Rcpt(v.Email)
		if err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(m.Message)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	err = c.Close()
	if err != nil {
		return err
	}
	return nil
}
