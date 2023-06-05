package library

import (
	"crypto/tls"
	"strconv"

	"gopkg.in/gomail.v2"
)

type Smtp struct {
	Attach      []string
	Body        string
	ContentType string
	Embed       []string
	From        []string
	Subject     []string
	To          []string
}

func (s Smtp) Send(r Redis) error {
	mail := gomail.NewMessage()
	mail.SetHeaders(map[string][]string{
		"From":    s.From,
		"Subject": s.Subject,
		"To":      s.To,
	})
	for _, v := range s.Attach {
		mail.Attach(v)
	}
	for _, v := range s.Embed {
		mail.Embed(v)
	}
	mail.SetBody(s.ContentType, s.Body)
	port, _ := strconv.Atoi(r.GetSecret("MAIL_PORT"))
	dialer := gomail.NewDialer(r.GetSecret("MAIL_HOST"), port, r.GetSecret("MAIL_USER"), r.GetSecret("MAIL_PASS"))
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return dialer.DialAndSend(mail)
}
