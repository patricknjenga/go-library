package library

import (
	"fmt"
	"github.com/pkg/sftp"
	"io"
	"os"
)

type Sftp struct {
	*sftp.Client `gorm:"-"`
	Directory    string
	Ssh
}

func (s Sftp) New(r Redis) (Sftp, error) {
	var err error
	s.Ssh, err = s.Ssh.New(r)
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

func (s Sftp) Ls() ([]os.FileInfo, error) {
	return s.Client.ReadDir(s.Directory)
}

func (s Sftp) Mkdir() error {
	return s.Client.MkdirAll(s.Directory)
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

func (s Sftp) GetAddress() string   { return s.Address }
func (s Sftp) GetDirectory() string { return s.Directory }
func (s Sftp) GetPort() string      { return s.Port }
func (s Sftp) GetUser() string      { return s.User }
