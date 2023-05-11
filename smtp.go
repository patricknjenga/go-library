package library

import (
	"fmt"
	"net/smtp"
)

type Mail struct {
	Message    []byte
	Recipients []string
}

func (m Mail) Send(r Redis) error {
	c, err := smtp.Dial(fmt.Sprintf("%s:%s", r.GetSecret("SMTP_HOST"), r.GetSecret("SMTP_PORT")))
	err = c.Mail(r.GetSecret("SMTP_USER"))
	if err != nil {
		return err
	}
	for _, v := range m.Recipients {
		err = c.Rcpt(v)
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
