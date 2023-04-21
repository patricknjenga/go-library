package library

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

type Ssh struct {
	*ssh.Client
	Address    string
	Password   string
	Port       string
	PrivateKey string
	User       string
}

func (s Ssh) New() (Ssh, error) {
	signer, err := ssh.ParsePrivateKey([]byte(s.PrivateKey))
	if err != nil {
		return s, err
	}
	s.Client, err = ssh.Dial("tcp", fmt.Sprintf(`%s:%s`, s.Address, s.Port), &ssh.ClientConfig{
		Auth: []ssh.AuthMethod{
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
				answers := make([]string, len(questions))
				for i := range answers {
					answers[i] = s.Password
				}
				return answers, nil
			}),
			ssh.Password(s.Password),
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
	defer session.Close()
	return session.Output(command)
}
