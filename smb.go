package library

import (
	"fmt"
	"io"
	"net"
	"os"
	"regexp"

	"github.com/hirochachacha/go-smb2"
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
	s.Session, err = (&smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     s.User,
			Password: r.GetSecret(s.Password),
		},
	}).Dial(connection)
	if err != nil {
		return s, err
	}
	s.Share, err = s.Session.Mount(s.Mount)
	if err != nil {
		return s, err
	}
	return s, err
}

func (s Smb) Get(p string) ([]byte, error) {
	reader, err := s.Share.Open(fmt.Sprintf("%s/%s", s.Directory, p))
	if err != nil {
		return []byte{}, err
	}
	return io.ReadAll(reader)
}

func (s Smb) Ls(regex string) ([]os.FileInfo, error) {
	var result []os.FileInfo
	regexp, err := regexp.Compile(regex)
	if err != nil {
		return result, err
	}
	fileInfos, err := s.Share.ReadDir(s.Directory)
	if err != nil {
		return result, err
	}
	for _, v := range fileInfos {
		if regexp.Match([]byte(v.Name())) {
			result = append(result, v)
		}
	}
	return result, err
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

func (s Smb) GetAddress() string { return s.Address }

func (s Smb) GetDirectory() string { return s.Directory }

func (s Smb) GetPort() string { return s.Port }

func (s Smb) GetUser() string { return s.User }
