package library

import (
	"fmt"
)

type Sms struct {
	Message string
	Phone   string
}

func (s Sms) Send(r Redis) ([]byte, error) {
	ssh, err := Ssh{Address: r.GetSecret("SMPP_SSH_HOST"), Password: r.GetSecret("SMPP_SSH_PASS"), Port: r.GetSecret("SMPP_SSH_PORT"), User: r.GetSecret("SMPP_SSH_USER")}.New(r)
	if err != nil {
		return nil, err
	}
	return ssh.Exec(fmt.Sprintf(`perl <<<"
		use Net::SMPP;
		Net::SMPP
		->new_transceiver('%s', port => '%s', system_id => '%s', password => '%s', system_type => 'smpp')
		->submit_sm(destination_addr => '%s', short_message => '%s');
	"`, r.GetSecret("SMPP_HOST"), r.GetSecret("SMPP_PORT"), r.GetSecret("SMPP_USER"), r.GetSecret("SMPP_PASS"), s.Phone, s.Message))
}
