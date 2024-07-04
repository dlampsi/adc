package examples

import "github.com/dlampsi/adc"

func mainCustomSearchBase() {
	cfg := &adc.Config{
		URL: "ldaps://my.ad.site:636",
		Bind: &adc.BindAccount{
			DN:       "CN=admin,DC=company,DC=com",
			Password: "***",
		},
		SearchBase: "OU=default,DC=company,DC=com",
		// Custom search base for users
		Users: &adc.UsersConfigs{
			SearchBase: "OU=users,DC=company,DC=com",
		},
		// Custom search base for groups
		Groups: &adc.GroupsConfigs{
			SearchBase: "OU=groups,DC=company,DC=com",
		},
	}

	cl := adc.New(cfg)

	if err := cl.Connect(); err != nil {
		panic(err)
	}

	// Do stuff...
}
