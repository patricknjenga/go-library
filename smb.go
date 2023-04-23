package library

import (
	"fmt"
	"github.com/hirochachacha/go-smb2"
	"io"
	"net"
	"os"
)

type Smb struct {
	*smb2.Session `gorm:"-"`
	*smb2.Share   `gorm:"-"`
	Address       string
	Directory     string
	Mount         string
	Password      string
	Port          string
	User          string
}

func (s Smb) New(r Redis) (Smb, error) {
	connection, err := net.Dial("tcp", fmt.Sprintf("%s:%s", s.Address, s.Port))
	if err != nil {
		return s, err
	}
	dialer := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     s.User,
			Password: r.GetSecret(s.Password),
		},
	}
	session, err := dialer.Dial(connection)
	if err != nil {
		return s, err
	}
	share, err := session.Mount(s.Mount)
	if err != nil {
		return s, err
	}
	s.Share = share
	s.Session = session
	return s, err
}

func (s Smb) Get(p string) ([]byte, error) {
	reader, err := s.Share.Open(fmt.Sprintf("%s/%s", s.Directory, p))
	if err != nil {
		return []byte{}, err
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func (s Smb) Ls() ([]os.FileInfo, error) {
	return s.Share.ReadDir(s.Directory)
}

func (s Smb) Mkdir() error {
	return s.Share.MkdirAll(s.Directory, 0777)
}

func (s Smb) Put(path string, data []byte) (int, error) {
	writer, err := s.Share.Create(fmt.Sprintf("%s/%s", s.Directory, path))
	if err != nil {
		return 0, err
	}
	return writer.Write(data)
}

func (s Smb) GetAddress() string   { return s.Address }
func (s Smb) GetDirectory() string { return s.Directory }
func (s Smb) GetPort() string      { return s.Port }
func (s Smb) GetUser() string      { return s.User }
