package library

import (
	"fmt"
	"github.com/hirochachacha/go-smb2"
	"net"
)

type SmbConnector struct {
	*smb2.Session
	*smb2.Share
	Address   string
	Directory string
	Mount     string
	Password  string
	Pattern   string
	User      string
	paths     []string
}

func (s *SmbConnector) init() error {
	connection, err := net.Dial("tcp", s.Address)
	if err != nil {
		return err
	}
	dialer := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     s.User,
			Password: s.Password,
		},
	}
	session, err := dialer.Dial(connection)
	if err != nil {
		return err
	}
	share, err := session.Mount(s.Mount)
	if err != nil {
		return err
	}
	s.Share = share
	s.Session = session
	return nil
}

func (s *SmbConnector) mkdir() error {
	err := s.Share.MkdirAll(s.Directory, 0777)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmbConnector) upload(path string, data []byte) (int, error) {
	_path := fmt.Sprintf("%s/%s", s.Directory, path)
	writer, err := s.Share.Create(_path)
	if err != nil {
		return 0, err
	}
	return writer.Write(data)
}

func (s *SmbConnector) download(p string) ([]byte, error) {
	return []byte{}, nil
}

func (s *SmbConnector) list() ([]string, error) {
	return []string{}, nil
}
