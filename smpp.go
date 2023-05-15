package library

import (
	"fmt"
)

type Smpp struct {
	Messages   []string
	Recipients []string
}

func (s Smpp) Send(r Redis) error {
	ssh, err := Ssh{Address: r.GetSecret("SMPP_SSH_HOST"), Password: r.GetSecret("SMPP_SSH_PASS"), Port: r.GetSecret("SMPP_SSH_PORT"), User: r.GetSecret("SMPP_SSH_USER")}.New(r)
	if err != nil {
		return err
	}
	for _, v := range s.Recipients {
		for _, w := range s.Messages {
			_, err := ssh.Exec(fmt.Sprintf(`perl <<<"
				use Net::SMPP;
				Net::SMPP
				->new_transceiver('%s', port => '%s', system_id => '%s', password => '%s', system_type => 'smpp')
				->submit_sm(destination_addr => '%s', short_message => '%s');
			"`, r.GetSecret("SMPP_HOST"), r.GetSecret("SMPP_PORT"), r.GetSecret("SMPP_USER"), r.GetSecret("SMPP_PASS"), v, w))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
