package adc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_AppendUsesAttributes(t *testing.T) {
	cfg := &Config{
		Users: &UsersConfigs{
			Attributes: []string{"one"},
		},
	}

	cfg.AppendUsesAttributes()
	require.Equal(t, []string{"one"}, cfg.Users.Attributes)

	cfg.AppendUsesAttributes("two")
	require.Equal(t, []string{"one", "two"}, cfg.Users.Attributes)
}

func Test_AppendGroupsAttributes(t *testing.T) {
	cfg := &Config{
		Groups: &GroupsConfigs{
			Attributes: []string{"one"},
		},
	}
	cfg.AppendGroupsAttributes()
	require.Equal(t, []string{"one"}, cfg.Groups.Attributes)

	cfg.AppendGroupsAttributes("two")
	require.Equal(t, []string{"one", "two"}, cfg.Groups.Attributes)
}

func Test_popConfig(t *testing.T) {
	cfg := &Config{
		Timeout: 5 * time.Second,
		Users:   DefaultUsersConfigs(),
		Groups:  DefaultGroupsConfigs(),
	}

	cl := &Client{cfg: cfg}

	cl.popConfig(nil)
	require.Equal(t, DefaultUsersConfigs(), cl.cfg.Users)
	require.Equal(t, DefaultGroupsConfigs(), cl.cfg.Groups)
	require.Nil(t, cl.cfg.Bind)

	cfg.URL = "ldaps://fakeurl:636"
	cfg.InsecureTLS = true
	cfg.Timeout = 42 * time.Second
	cfg.SearchBase = "OU=some"
	cl.popConfig(cfg)
	require.Equal(t, cfg.URL, cl.cfg.URL)
	require.Equal(t, cfg.InsecureTLS, cl.cfg.InsecureTLS)
	require.Equal(t, cfg.Timeout, cl.cfg.Timeout)
	require.Equal(t, cfg.SearchBase, cl.cfg.SearchBase)
	require.Nil(t, cl.cfg.Bind)
	require.Equal(t, cfg.SearchBase, cl.cfg.Users.SearchBase)
	require.Equal(t, cfg.SearchBase, cl.cfg.Groups.SearchBase)

	cfg.Bind = &BindAccount{DN: "some", Password: "fake"}
	cl.popConfig(cfg)
	require.NotNil(t, cl.cfg.Bind)
	require.Equal(t, cfg.Bind, cl.cfg.Bind)

	cfg.Users.SearchBase = "OU=custom-users"
	cfg.Groups.SearchBase = "OU=custom-groups"
	cfg.AppendUsesAttributes("dummy-user-attr")
	cfg.AppendGroupsAttributes("dummy-group-attr")
	cl.popConfig(cfg)
	require.Equal(t, cfg.Users.SearchBase, cl.cfg.Users.SearchBase)
	require.Equal(t, cfg.Groups.SearchBase, cl.cfg.Groups.SearchBase)
	require.Contains(t, cl.cfg.Users.Attributes, "dummy-user-attr")
	require.Contains(t, cl.cfg.Groups.Attributes, "dummy-group-attr")
}
