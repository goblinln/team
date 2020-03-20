package auth

// LoginedAccount returns by LoginProcessor.
type LoginedAccount struct {
	Account string
	Name    string
}

// LoginProcessor for custom login method.
type LoginProcessor interface {
	Login(account, password string) error
}
