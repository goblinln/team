package auth

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/ldap.v3"
)

// LDAPProtocol definition.
type LDAPProtocol int

const (
	// LDAPUnencrypted protocol
	LDAPUnencrypted LDAPProtocol = iota
	// LDAPTLS protocol
	LDAPTLS
	// LDAPStartTLS protocol
	LDAPStartTLS
)

// LDAPProcessor implements auth using LDAP
type LDAPProcessor struct {
	Host       string
	Port       int
	Protocol   LDAPProtocol
	SkipVerify bool
}

// Login implements LoginProcessor interface.
func (l *LDAPProcessor) Login(account, password string) error {
	var (
		conn *ldap.Conn
		err  error
	)

	if l.Protocol == LDAPTLS {
		conn, err = ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", l.Host, l.Port), &tls.Config{
			ServerName:         l.Host,
			InsecureSkipVerify: l.SkipVerify,
		})
	} else {
		conn, err = ldap.Dial("tcp", fmt.Sprintf("%s:%d", l.Host, l.Port))
	}

	if err != nil {
		return fmt.Errorf("Failed to connect to LDAP server: %v", err)
	}

	defer conn.Close()

	if l.Protocol == LDAPStartTLS {
		if err = conn.StartTLS(&tls.Config{ServerName: l.Host, InsecureSkipVerify: l.SkipVerify}); err != nil {
			return fmt.Errorf("Failed to STARTTLS: %v", err)
		}
	}

	return conn.Bind(account, password)
}
