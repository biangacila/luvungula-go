package email

import "testing"

func TestOutputCustomer_Say(t *testing.T) {
	o := OutputCustomer{}
	o.Say("lubungula")
}
