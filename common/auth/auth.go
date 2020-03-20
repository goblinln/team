package auth

// Kind definition
type Kind int

const (
	// KindNone using build-in account
	KindNone Kind = iota
	// KindSMTP using SMTP auth
	KindSMTP
	// KindLDAP using LDAP auth
	KindLDAP
)

// LoginProcessor for custom login method.
type LoginProcessor interface {
	Login(account, password string) error
}
