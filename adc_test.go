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
		require.NotNil(t, cl.cfg)
	})
	t.Run("WithLogger", func(t *testing.T) {
		tlogger := &testLogger{}
		cl := New(nil, WithLogger(tlogger))
		require.NotNil(t, cl)
		require.Equal(t, tlogger, cl.logger)
	})
}

func Test_Client_Config(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		cl := New(nil)
		require.NotNil(t, cl.Config())
	})
}

func Test_Client_Connect(t *testing.T) {
	badCl := New(nil)
	expErr := badCl.Connect()
	require.Error(t, expErr)

	cl := New(nil, withMock())

	err := cl.Connect()
	require.NoError(t, err)
	require.NoError(t, cl.Disconnect())

	cfg := &Config{
		Bind: &BindAccount{DN: "fakeone", Password: "fake"},
	}
	cl = New(cfg, withMock())
	err = cl.Connect()
	require.Error(t, err)

	cfg.Bind.DN = "mrError"
	cl = New(cfg)
	err = cl.Connect()
	require.Error(t, err)

	cfg.Bind.DN = "validUser"
	cfg.Bind.Password = "badPass"
	cl = New(cfg, withMock())
	err = cl.Connect()
	require.Error(t, err)

	cfg.Bind.DN = "validUser"
	cfg.Bind.Password = "validPass"
	cl = New(cfg, withMock())
	err = cl.Connect()
	require.NoError(t, err)
}

func Test_Client_Reconnect(t *testing.T) {
	cfg := &Config{
		Bind: validMockBind,
	}
	cl := New(cfg, withMock())
	err := cl.Connect()
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = cl.Reconnect(ctx, 2*time.Second, 2)
	require.NoError(t, err)

	cl.cfg.Bind = &BindAccount{DN: mockEntriesData["entryForErr"].DN}
	err = cl.Reconnect(ctx, 2*time.Second, 2)
	require.Error(t, err)

	cl.cfg.Bind = reconnectMockBind
	err = cl.Reconnect(ctx, 0, 1)
	require.Error(t, err)
	err = cl.Reconnect(ctx, 30*time.Millisecond, 0)
	require.Error(t, err)
	err = cl.Reconnect(ctx, 1*time.Second, 1)
	require.Error(t, err)

	nctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	err = cl.Reconnect(nctx, 5*time.Second, 1)
	require.Error(t, err)

	cl.cfg.Bind = validMockBind
	err = cl.Reconnect(ctx, 30*time.Millisecond, 1)
	require.NoError(t, err)
}
