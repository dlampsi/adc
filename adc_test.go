package adc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	cl := New(nil)
	require.NotNil(t, cl)
	require.Equal(t, cl.cfg.Timeout, defaultClientTimeout)
	require.NotNil(t, cl.cfg.Users)
	require.Equal(t, cl.cfg.Users, DefaultUsersConfigs())
	require.NotNil(t, cl.cfg.Groups)
	require.Equal(t, cl.cfg.Groups, DefaultGroupsConfigs())

	cfg := &Config{Timeout: 60 * time.Second}
	cl = New(cfg)
	require.NotNil(t, cl)
	require.Equal(t, cfg.Timeout, cl.Config().Timeout)
}

func Test_Config(t *testing.T) {
	cfg := &Config{
		Users:  DefaultUsersConfigs(),
		Groups: DefaultGroupsConfigs(),
	}
	cl := New(cfg)
	cl.Config().Users.FilterById = "(&(objectClass=group)(cn=%v))"
	require.Equal(t, "(&(objectClass=group)(cn=%v))", cl.Config().Users.FilterById)
}

func Test_Client_Connect(t *testing.T) {
	badCl := New(nil)
	expErr := badCl.Connect()
	require.Error(t, expErr)

	mock := &mockClient{}
	cl := New(nil, WithLdapClient(mock))

	err := cl.Connect()
	require.NoError(t, err)
	cl.Disconnect()

	cfg := &Config{
		Bind: &BindAccount{DN: "fakeone", Password: "fake"},
	}
	cl = New(cfg, WithLdapClient(mock))
	err = cl.Connect()
	require.Error(t, err)

	cfg.Bind.DN = "mrError"
	cl = New(cfg, WithLdapClient(mock))
	err = cl.Connect()
	require.Error(t, err)

	cfg.Bind.DN = "validUser"
	cfg.Bind.Password = "badPass"
	cl = New(cfg, WithLdapClient(mock))
	err = cl.Connect()
	require.Error(t, err)

	cfg.Bind.DN = "validUser"
	cfg.Bind.Password = "validPass"
	cl = New(cfg, WithLdapClient(mock))
	err = cl.Connect()
	require.NoError(t, err)
}

func Test_Client_Reconnect(t *testing.T) {
	mock := &mockClient{}
	cfg := &Config{
		Bind: validMockBind,
	}
	cl := New(cfg, WithLdapClient(mock))
	err := cl.Connect()
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = cl.Reconnect(ctx, time.NewTicker(2*time.Second), 2)
	require.NoError(t, err)

	cl.cfg.Bind = &BindAccount{DN: mockEntriesData["entryForErr"].DN}
	err = cl.Reconnect(ctx, time.NewTicker(2*time.Second), 2)
	require.Error(t, err)

	cl.cfg.Bind = reconnectMockBind
	err = cl.Reconnect(ctx, nil, 1)
	require.Error(t, err)
	err = cl.Reconnect(ctx, time.NewTicker(30*time.Millisecond), 0)
	require.Error(t, err)
	err = cl.Reconnect(ctx, time.NewTicker(1*time.Second), 1)
	require.Error(t, err)

	nctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	err = cl.Reconnect(nctx, time.NewTicker(5*time.Second), 1)
	require.Error(t, err)

	cl.cfg.Bind = validMockBind
	err = cl.Reconnect(ctx, time.NewTicker(30*time.Millisecond), 1)
	require.NoError(t, err)
}
