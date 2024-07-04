package examples

import (
	"fmt"

	"github.com/dlampsi/adc"
)

func mainCustomSearchFilters() {
	cfg := &adc.Config{
		URL: "ldaps://my.ad.site:636",
		Bind: &adc.BindAccount{
			DN:       "CN=admin,DC=company,DC=com",
			Password: "***",
		},
		SearchBase: "OU=default,DC=company,DC=com",
		// Custom filter for all users search requests
		Users: &adc.UsersConfigs{
			FilterById: "(&(objectClass=person)(cn=%v))",
		},
		// Custom filter for all groups search requests
		Groups: &adc.GroupsConfigs{
			FilterById: "(&(objectClass=group)(cn=%v))",
		},
	}

	cl := adc.New(cfg)

	if err := cl.Connect(); err != nil {
		panic(err)
	}

	// Do stuff...

	// Custom search filter per search
	// Note that provided `Filter` argument int `GetUserArgs` overwrites `Id` and `Dn` arguments usage.

	user, err := cl.GetUser(adc.GetUserArgs{Filter: "(&(objectClass=person)(sAMAccountName=someID))"})
	if err != nil {
		panic(err)
	}
	if user == nil {
		panic("User not found")
	}
	fmt.Println(user)
}
