package library

import (
	"fmt"
	"net/smtp"
)

type Smtp struct {
	*smtp.Client
	Host     string
	Password string
	Port     string
	User     string
}

func (s Smtp) New() (Smtp, error) {
	var err error
	s.Client, err = smtp.Dial(fmt.Sprintf("%s:%s", s.Host, s.Port))
	if err != nil {
		return s, err
	}
	return s, err
}

func (s Smtp) Send(message []byte, recipients []string) error {
	err := s.Client.Mail(s.User)
	if err != nil {
		return err
	}
	for _, v := range recipients {
		err = s.Client.Rcpt(v)
		if err != nil {
			return err
		}
	}
	w, err := s.Client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(message)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	err = s.Client.Close()
	if err != nil {
		return err
	}
	return nil
}
