package frontend

import (
	"errors"
	"encoding/binary"
	"bytes"
	"io/ioutil"
	"log"
)

// PatchNwcConfig takes an original config, applies needed patches to the URL and such,
// updates the checksum and returns either nil, error or a patched config w/o error.
func PatchNwcConfig(originalConfig []byte) ([]byte, error) {
	if len(originalConfig) != 1024 {
		return nil, errors.New("invalid config size")
	}

	var config ConfigFormat
	configReadingBuf := bytes.NewBuffer(originalConfig)
	err := binary.Read(configReadingBuf, binary.BigEndian, &config)
	if err != nil {
		return nil, err
	}

	if bytes.Compare(config.Magic[:], ConfigMagic) != 0 {
		return nil, errors.New("invalid magic")
	}

	copy(config.MailDomain[:], []byte("@mail.disconnect24.xyz"))

	// The following is very redundantly written. TODO: fix that?
	copy(config.AccountURL[:128], []byte("http://mail.disconnect24.xyz/cgi-bin/account.cgi"))
	copy(config.CheckURL[:128], []byte("http://mail.disconnect24.xyz/cgi-bin/check.cgi"))
	copy(config.ReceiveURL[:128], []byte("http://mail.disconnect24.xyz/cgi-bin/receive.cgi"))
	copy(config.DeleteURL[:128], []byte("http://mail.disconnect24.xyz/cgi-bin/delete.cgi"))
	copy(config.SendURL[:128], []byte("http://mail.disconnect24.xyz/cgi-bin/send.cgi"))

	// Read from struct to buffer
	fileBuf := new(bytes.Buffer)
	err = binary.Write(fileBuf, binary.BigEndian, config)
	if err != nil {
		return nil, err
	}
	patchedConfig, err := ioutil.ReadAll(fileBuf)
	if err != nil {
		return nil, err
	}

	var checksumInt uint32

	// Checksum.
	// We loop from 1020 to avoid current checksum.
	// Take every 4 bytes, add 'er up!
	for i := 0; i < 1020; i += 4 {
		log.Printf("hello, %s", i)
		addition := binary.BigEndian.Uint32(patchedConfig[i:i+4])
		checksumInt += addition
	}

	// Grab lower 32 bits of int
	var finalChecksum uint32
	finalChecksum = checksumInt & 0xFFFFFFFF
	binaryChecksum := make([]byte, 4)
	binary.BigEndian.PutUint32(binaryChecksum, finalChecksum)

	// Update patched config checksum
	copy(patchedConfig[1020:1024], binaryChecksum)
	return patchedConfig, nil
}
