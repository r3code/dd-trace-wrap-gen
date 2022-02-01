package tests

import "golang.org/x/net/context"

// TestInterface is used to test templates
type TestInterface interface {
	F(ctx context.Context, a1 string, a2 ...string) (result1, result2 string, err error)
	NoError(string) string
	NoParamsOrResults()
	Channels(chA chan bool, chB chan<- bool, chanC <-chan bool)
}
