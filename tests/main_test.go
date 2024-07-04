package adctests

import (
	"fmt"
	"os"
	"testing"

	"github.com/dlampsi/adc"
)

var (
	tClient *adc.Client
)

type logger struct {
	t *testing.T
}

func (l *logger) Debug(args ...interface{}) {
	msg := l.t.Name() + ": " + fmt.Sprint(args...)
	l.t.Log(msg)
}

func (l *logger) Debugf(template string, args ...interface{}) {
	msg := l.t.Name() + ": " + fmt.Sprintf(template, args...)
	l.t.Log(msg)
}

func getClientConfig() adc.Config {
	return adc.Config{
		URL:         "ldaps://127.0.0.1:636",
		SearchBase:  "CN=Users,DC=adc,DC=dev",
		InsecureTLS: true,
		Bind: &adc.BindAccount{
			DN:       "CN=Administrator,CN=Users,DC=adc,DC=dev",
			Password: os.Getenv("TESTS_AD_USER_PWD"),
		},
	}
}

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}

func testMainWrapper(m *testing.M) int {
	cfg := getClientConfig()
	tClient = adc.New(&cfg)

	if err := tClient.Connect(); err != nil {
		fmt.Println(fmt.Errorf("Failed to connect: %v", err))
		return 1
	}

	return m.Run()
}
