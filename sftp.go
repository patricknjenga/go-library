package library

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/pkg/sftp"
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
	return io.ReadAll(reader)
}

func (s Sftp) Ls(regex string) ([]os.FileInfo, error) {
	var result []os.FileInfo
	regexp, err := regexp.Compile(regex)
	if err != nil {
		return result, err
	}
	fileInfos, err := s.Client.ReadDir(s.Directory)
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

func (s Sftp) Mkdir() error {
	return s.Client.MkdirAll(s.Directory)
}

func (s Sftp) Put(p string, data []byte) (int, error) {
	writer, err := s.Client.Create(fmt.Sprintf("%s/%s", s.Directory, p))
	if err != nil {
		return 0, err
	}
	return writer.Write(data)
}

func (s Sftp) GetAddress() string { return s.Address }

func (s Sftp) GetDirectory() string { return s.Directory }

func (s Sftp) GetPort() string { return s.Port }

func (s Sftp) GetUser() string { return s.User }
