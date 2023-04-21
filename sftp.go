package library

import (
	"fmt"
	"github.com/pkg/sftp"
	"io"
)

type Sftp struct {
	*sftp.Client
	Directory string
	Ssh
}

func (s Sftp) New() (Sftp, error) {
	var err error
	s.Ssh, err = s.Ssh.New()
	if err != nil {
		return s, err
	}
	s.Client, err = sftp.NewClient(s.Ssh.Client)
	if err != nil {
		return s, err
	}
	return s, err
}

func (s Sftp) Get(p string) ([]byte, error) {
	reader, err := s.Client.Open(fmt.Sprintf("%s/%s", s.Directory, p))
	if err != nil {
		return []byte{}, err
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func (s Sftp) Put(p string, data []byte) error {
	writer, err := s.Client.Create(fmt.Sprintf("%s/%s", s.Directory, p))
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	return nil
}
