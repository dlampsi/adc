package examples

import (
	"fmt"

	"github.com/dlampsi/adc"
)

func mainCustomAttributes() {
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

	// Append custom user attributes to all searches and get it.

	cl.Config.AppendUsersAttributes("manager")

	user, err := cl.GetUser(adc.GetUserArgs{
		Id: "exampleUserId",
	})
	if err != nil {
		panic(err)
	}
	if user == nil {
		panic("User not found")
	}

	fmt.Println(user.GetStringAttribute("manager"))

	// Append custom user attributes per-search.

	user2, err := cl.GetUser(adc.GetUserArgs{
		Id:         "exampleUserId",
		Attributes: []string{"manager"},
	})
	if err != nil {
		panic(err)
	}
	if user2 == nil {
		panic("User not found")
	}
	fmt.Println(user2.GetStringAttribute("manager"))
}
