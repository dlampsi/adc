package examples

import (
	"fmt"

	"github.com/dlampsi/adc"
)

func mainUsers() {
	cfg := &adc.Config{
		URL: "ldaps://my.ad.site:636",
		Bind: &adc.BindAccount{
			DN:       "CN=admin,DC=company,DC=com",
			Password: "***",
		},
		SearchBase: "OU=default,DC=company,DC=com",
	}

	cl := adc.New(cfg)

	if err := cl.Connect(); err != nil {
		panic(err)
	}

	/* -------------- Create -------------- */

	createReq := adc.CreateUserArgs{
		Id:       "exampleUserId",
		Password: "examplePassword",
		Attributes: map[string][]string{
			"sn": {"exampleUserSurname"},
		},
	}
	if err := cl.CreateUser(createReq); err != nil {
		panic(err)
	}

	/* -------------- Search -------------- */

	getReq := adc.GetUserArgs{
		Id: "exampleUserId",
	}
	user, err := cl.GetUser(getReq)
	if err != nil {
		panic(err)
	}
	if user == nil {
		panic("User not found")
	}
	fmt.Println(user)

	/* -------------- Delete -------------- */

	if err := cl.DeleteUser(user.DN); err != nil {
		panic(err)
	}
}
