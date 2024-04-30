package adc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testLogger struct{}

func (l *testLogger) Debug(args ...interface{})                   {}
func (l *testLogger) Debugf(template string, args ...interface{}) {}

func Test_New(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		cl := New(nil)
		require.NotNil(t, cl)
		require.NotNil(t, cl.Config)
	})
	t.Run("WithLogger", func(t *testing.T) {
		tlogger := &testLogger{}
		cl := New(nil, WithLogger(tlogger))
		require.NotNil(t, cl)
		require.Equal(t, tlogger, cl.logger)
	})
}

func Test_Client_Connect(t *testing.T) {
	t.Run("BadClient", func(t *testing.T) {
		cl := New(nil)
		require.Error(t, cl.Connect(), "Connection should fail without any configuration")
	})
	t.Run("Bindless", func(t *testing.T) {
		cl := newMockClient(nil)
		require.NoError(t, cl.Connect(), "Bindless connection should be successful")
		require.NoError(t, cl.Disconnect())
	})
	t.Run("BindWithErr", func(t *testing.T) {
		cl := newMockClient(&Config{
			Bind: &BindAccount{DN: "mrError", Password: "fake"},
		})
		require.Error(t, cl.Connect(), "Connection should fail with bad bind credentials")
	})
	t.Run("BindWithSuccess", func(t *testing.T) {
		cl := newMockClient(&Config{
			Bind: validMockBind,
		})
		require.NoError(t, cl.Connect(), "Connection should be successful with valid bind credentials")
	})
	t.Run("WithSecureDialOpts", func(t *testing.T) {
		cl := New(&Config{
			URL: "ldaps://fake:636",
		})
		require.Error(t, cl.Connect())
	})
	t.Run("WithSecureDialOpts", func(t *testing.T) {
		cl := New(&Config{
			URL:         "ldaps://fake:636",
			InsecureTLS: true,
		})
		require.Error(t, cl.Connect())
	})
}

func Test_Client_Disconnect(t *testing.T) {
	t.Run("WithNilLdap", func(t *testing.T) {
		cl := newMockClient(&Config{Bind: validMockBind})
		require.NoError(t, cl.Disconnect())
	})
	t.Run("Ok", func(t *testing.T) {
		cl := newMockClient(&Config{Bind: validMockBind})
		require.NoError(t, cl.Connect())
		require.NoError(t, cl.Disconnect())
	})
}

func Test_Client_Reconnect(t *testing.T) {
	t.Run("WithErr", func(t *testing.T) {
		cl := newMockClient(&Config{Bind: validMockBind})
		require.NoError(t, cl.Connect())

		cl.Config.Bind = &BindAccount{DN: "OU=entryForErr,DC=company,DC=com"}
		require.Error(t, cl.Reconnect(context.TODO(), 2*time.Second, 2), "Reconnect for a bad user should fail with error")
	})
	t.Run("WithReconnectErr", func(t *testing.T) {
		ctx := context.TODO()

		cl := newMockClient(&Config{Bind: validMockBind})
		require.NoError(t, cl.Connect())

		cl.Config.Bind = reconnectMockBind

		require.Error(t, cl.Reconnect(ctx, 0, 1))
		require.Error(t, cl.Reconnect(ctx, 30*time.Millisecond, 0))
		require.Error(t, cl.Reconnect(ctx, 1*time.Second, 1))
	})
	t.Run("WithContextCancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		cl := newMockClient(&Config{Bind: validMockBind})
		require.NoError(t, cl.Connect())

		cl.Config.Bind = reconnectMockBind

		require.Error(t, cl.Reconnect(ctx, 5*time.Second, 1))

		cl.Config.Bind = validMockBind
		require.NoError(t, cl.Reconnect(ctx, 30*time.Millisecond, 1))
	})
	t.Run("Ok", func(t *testing.T) {
		ctx := context.TODO()

		cl := newMockClient(&Config{Bind: validMockBind})

		require.NoError(t, cl.Connect())
		require.NoError(t, cl.Reconnect(ctx, 2*time.Second, 2), "Reconnect should be successful")
	})
}

func Test_Client_CheckAuthByDN(t *testing.T) {
	t.Run("WithCheckErr", func(t *testing.T) {
		cl := newMockClient(&Config{Bind: validMockBind})
		require.NoError(t, cl.Connect())
		require.Error(t, cl.CheckAuthByDN("OU=entryForErr,DC=company,DC=com", "fake"))
	})
	t.Run("Ok", func(t *testing.T) {
		cl := newMockClient(&Config{Bind: validMockBind})
		require.NoError(t, cl.Connect())
		require.NoError(t, cl.CheckAuthByDN(validMockBind.DN, validMockBind.Password))
	})
}
