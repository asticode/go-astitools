package astissh

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// AuthMethodPublicKey creates a public key auth method
func AuthMethodPublicKey(paths ...string) (m ssh.AuthMethod, err error) {
	// Loop through paths
	var ss []ssh.Signer
	for _, p := range paths {
		// Read private key
		var b []byte
		if b, err = ioutil.ReadFile(p); err != nil {
			err = errors.Wrapf(err, "main: reading private key %s failed", p)
			return
		}

		// Parse private key
		var s ssh.Signer
		if s, err = ssh.ParsePrivateKey(b); err != nil {
			err = errors.Wrapf(err, "main: parsing private key %s failed", p)
			return
		}

		// Append
		ss = append(ss, s)
	}

	// Create auth method
	m = ssh.PublicKeys(ss...)
	return
}
