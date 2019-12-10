package email

import "fmt"

type OutputCustomer struct {
}

func (obj *OutputCustomer) Say(in string) {
	fmt.Println(fmt.Sprintf(":)-> %v", in))
}
