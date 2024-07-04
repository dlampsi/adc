package adctests

import (
	"context"
	"testing"
	"time"

	"github.com/dlampsi/adc"
	"github.com/stretchr/testify/require"
)

type testLogger struct{}

func (l *testLogger) Debug(args ...interface{})                   {}
func (l *testLogger) Debugf(template string, args ...interface{}) {}

func Test_New(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		cl := adc.New(nil)
		require.NotNil(t, cl)
		require.NotNil(t, cl.Config)
	})
	t.Run("WithLogger", func(t *testing.T) {
		tlogger := &testLogger{}
		cl := adc.New(nil, adc.WithLogger(tlogger))
		require.NotNil(t, cl)
	})
}

func Test_Client_Connect(t *testing.T) {
	t.Run("BadClient", func(t *testing.T) {
		cl := adc.New(nil)
		require.Error(t, cl.Connect(), "Connection should fail without any configuration")
	})

	t.Run("BadBind", func(t *testing.T) {
		cfg := getClientConfig()
		cfg.Bind.Password = ""
		cl := adc.New(&cfg)
		require.NotNil(t, cl)
		require.Error(t, cl.Connect(), "Connection should fail without password")
	})

	t.Run("Connect_OK", func(t *testing.T) {
		require.NoError(t, tClient.Connect())
	})
}
func Test_Client_Reconnect(t *testing.T) {
	ctx := context.TODO()

	t.Run("Ok", func(t *testing.T) {
		cfg := getClientConfig()
		cl := adc.New(&cfg)
		require.NotNil(t, cl)

		require.NoError(t, cl.Connect())
		require.NoError(t, cl.Reconnect(ctx, 2*time.Second, 1), "No error should be returned if reconnect is not needed.")

		require.NoError(t, cl.Disconnect())
		require.NoError(t, cl.Reconnect(ctx, 2*time.Second, 1), "Reconnect should be successful.")
	})

	t.Run("OkWithDefaults", func(t *testing.T) {
		cfg := getClientConfig()
		cl := adc.New(&cfg)
		require.NotNil(t, cl)

		require.NoError(t, cl.Connect())
		require.NoError(t, cl.Disconnect())

		cl.Config.Bind.Password = "bad_password"

		require.Error(t, cl.Reconnect(ctx, 0, 0), "Here we're testing default attempts and retries, but reconnect should fail cause of bad password.")
	})

	t.Run("WithError", func(t *testing.T) {
		cfg := getClientConfig()
		cl := adc.New(&cfg)
		require.NotNil(t, cl)

		require.NoError(t, cl.Connect())
		require.NoError(t, cl.Disconnect())

		cl.Config.Bind.Password = "bad_password"
		require.Error(t, cl.Reconnect(ctx, 2*time.Second, 1))
	})

	t.Run("WithContextCancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		cfg := getClientConfig()
		cl := adc.New(&cfg)
		require.NotNil(t, cl)
		require.NoError(t, cl.Connect())
		require.NoError(t, cl.Disconnect())

		require.Error(t, cl.Reconnect(ctx, 30*time.Millisecond, 1))
	})
}

func Test_Client_Disconnect(t *testing.T) {
	cfg := getClientConfig()
	cl := adc.New(&cfg)
	require.NoError(t, cl.Disconnect())
}

func Test_Client_CheckAuthByDN(t *testing.T) {
	t.Run("Bad", func(t *testing.T) {
		require.Error(t, tClient.CheckAuthByDN(tClient.Config.Bind.DN, "bad_password"))
	})
	t.Run("Ok", func(t *testing.T) {
		require.NoError(t, tClient.CheckAuthByDN(tClient.Config.Bind.DN, tClient.Config.Bind.Password))
	})
	t.Run("WithConnectErr", func(t *testing.T) {
		cfg := getClientConfig()
		cl := adc.New(&cfg)
		require.NotNil(t, cl)
		require.NoError(t, cl.Disconnect())

		cl.Config.URL = "ldaps://fakeurl:636"

		require.Error(t, cl.CheckAuthByDN(tClient.Config.Bind.DN, tClient.Config.Bind.Password))
	})
}
