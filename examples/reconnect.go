package examples

import (
	"context"
	"time"

	"github.com/dlampsi/adc"
)

func mainReconnect(ctx context.Context) {
	// Init client
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

	// Recconect each 5 secconds with 24 retrie attempts
	if err := cl.Reconnect(ctx, 5*time.Second, 24); err != nil {
		panic(err)
	}
}
