package auth

// LoginProcessor for custom login method.
type LoginProcessor interface {
	Login(account, password string) error
}
