package frontend

import (
	"errors"
	"encoding/binary"
	"bytes"
)

func PatchNwcConfig(originalConfig []byte) (patchedConfig []byte, err error) {
	if len(originalConfig) != 1024 {
		return nil, errors.New("invalid config size")
	}

	var config ConfigFormat
	configReadingBuf := bytes.NewBuffer(originalConfig)
	err = binary.Read(configReadingBuf, binary.BigEndian, &config)
	if err != nil {
		return nil, err
	}

	if bytes.Compare(config.Magic[:], ConfigMagic) != 0 {
		return nil, errors.New("invalid magic")
	}

	copy(config.MailDomain[:], []byte("@mail.disconnect24.xyz"))

	// The following is very redundantly written. TODO: fix that?
	copy(config.AccountURL[:], []byte("http://mail.disconnect24.xyz/cgi-bin/account.cgi"))
	copy(config.CheckURL[:], []byte("http://mail.disconnect24.xyz/cgi-bin/account.cgi"))
	copy(config.ReceiveURL[:], []byte("http://mail.disconnect24.xyz/cgi-bin/account.cgi"))
	copy(config.DeleteURL[:], []byte("http://mail.disconnect24.xyz/cgi-bin/account.cgi"))
	copy(config.SendURL[:], []byte("http://mail.disconnect24.xyz/cgi-bin/account.cgi"))

	return nil, nil
}
