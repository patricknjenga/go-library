package library

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

type Ssh struct {
	*ssh.Client `gorm:"-"`
	Address     string
	Password    string
	Port        string
	User        string
}

func (s Ssh) New(r Redis) (Ssh, error) {
	signer, err := ssh.ParsePrivateKey([]byte(r.GetSecret("PRIVATE_KEY")))
	if err != nil {
		return s, err
	}
	s.Client, err = ssh.Dial("tcp", fmt.Sprintf(`%s:%s`, s.Address, s.Port), &ssh.ClientConfig{
		Auth: []ssh.AuthMethod{
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
				answers := make([]string, len(questions))
				for i := range answers {
					answers[i] = r.GetSecret(s.Password)
				}
				return answers, nil
			}),
			ssh.Password(r.GetSecret(s.Password)),
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		User:            s.User,
	})
	if err != nil {
		return s, err
	}
	return s, nil
}

func (s Ssh) Exec(command string) ([]byte, error) {
	session, err := s.Client.NewSession()
	if err != nil {
		return []byte{}, nil
	}
	return session.Output(command)
}
