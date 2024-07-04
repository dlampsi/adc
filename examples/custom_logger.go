package examples

import (
	"github.com/dlampsi/adc"
)

// Implementing a custom logger
type myLogger struct{}

func (l *myLogger) Debug(args ...interface{}) {
	// Your custom logger implementation
}

func (l *myLogger) Debugf(template string, args ...interface{}) {
	// Your custom logger implementation
}

func mainCustomLogger() {
	cfg := &adc.Config{
		URL: "ldaps://my.ad.site:636",
		Bind: &adc.BindAccount{
			DN:       "CN=admin,DC=company,DC=com",
			Password: "***",
		},
		SearchBase: "OU=default,DC=company,DC=com",
	}

	logger := &myLogger{}

	cl := adc.New(cfg, adc.WithLogger(logger))

	if err := cl.Connect(); err != nil {
		panic(err)
	}

	// Do stuff...
}
