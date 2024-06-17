package types

type TestCustomer struct {
	GivenName string `json:"given_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}
