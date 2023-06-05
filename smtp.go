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
	client, err := smtp.Dial(fmt.Sprintf("%s:%s", r.GetSecret("SMTP_HOST"), r.GetSecret("SMTP_PORT")))
	err = client.Mail(r.GetSecret("SMTP_USER"))
	if err != nil {
		return err
	}
	for _, v := range m.Recipients {
		err = client.Rcpt(v)
		if err != nil {
			return err
		}
	}
	w, err := client.Data()
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
	err = client.Close()
	if err != nil {
		return err
	}
	return nil
}
