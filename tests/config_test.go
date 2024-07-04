package adctests

import (
	"testing"
	"time"

	"github.com/dlampsi/adc"
	"github.com/stretchr/testify/require"
)

func Test_AppendUsesAttributes(t *testing.T) {
	cfg := &adc.Config{
		Users: &adc.UsersConfigs{
			Attributes: []string{"one"},
		},
	}

	cfg.AppendUsesAttributes()
	require.Equal(t, []string{"one"}, cfg.Users.Attributes)

	cfg.AppendUsesAttributes("two")
	require.Equal(t, []string{"one", "two"}, cfg.Users.Attributes)
}

func Test_AppendGroupsAttributes(t *testing.T) {
	cfg := &adc.Config{
		Groups: &adc.GroupsConfigs{
			Attributes: []string{"one"},
		},
	}
	cfg.AppendGroupsAttributes()
	require.Equal(t, []string{"one"}, cfg.Groups.Attributes)

	cfg.AppendGroupsAttributes("two")
	require.Equal(t, []string{"one", "two"}, cfg.Groups.Attributes)
}

func Test_Config(t *testing.T) {
	t.Run("CustomConfigPartial", func(t *testing.T) {
		cfg := &adc.Config{
			URL: "ldaps://fakeurl:636",
			Bind: &adc.BindAccount{
				DN:       "some",
				Password: "fake",
			},
			SearchBase: "OU=some",
			Users: &adc.UsersConfigs{
				SearchBase: "OU=custom-users",
			},
			Groups: &adc.GroupsConfigs{
				SearchBase: "OU=custom-groups",
			},
		}

		cl := adc.New(cfg)
		require.NotNil(t, cl.Config)

		require.Equal(t, cfg.URL, cl.Config.URL)
		require.Equal(t, cfg.InsecureTLS, cl.Config.InsecureTLS)
		require.Equal(t, cfg.SearchBase, cl.Config.SearchBase)
		require.Equal(t, cfg.Bind, cl.Config.Bind)
		require.Equal(t, cfg.Users.SearchBase, cl.Config.Users.SearchBase)
		require.Equal(t, cfg.Groups.SearchBase, cl.Config.Groups.SearchBase)
	})

	t.Run("CustomConfigAll", func(t *testing.T) {
		cfg := &adc.Config{
			URL:         "ldaps://fakeurl:636",
			InsecureTLS: true,
			Timeout:     5 * time.Second,
			Bind: &adc.BindAccount{
				DN:       "some",
				Password: "fake",
			},
			SearchBase: "OU=some",
			Users: &adc.UsersConfigs{
				IdAttribute:      "custom-users-id-attr",
				Attributes:       []string{"dummy-user-attr"},
				SearchBase:       "OU=custom-users",
				FilterById:       "customFilterById",
				FilterByDn:       "customFilterByDn",
				FilterGroupsByDn: "customFilterGroupsByDn",
			},
			Groups: &adc.GroupsConfigs{
				IdAttribute:       "custom-groups-id-attr",
				Attributes:        []string{"dummy-group-attr"},
				SearchBase:        "OU=custom-groups",
				FilterById:        "customFilterById",
				FilterByDn:        "customFilterByDn",
				FilterMembersByDn: "customFilterMembersByDn",
			},
		}

		cl := adc.New(cfg)
		require.NotNil(t, cl.Config)

		require.Equal(t, cfg.Timeout, cl.Config.Timeout)
		require.Equal(t, cfg.URL, cl.Config.URL)
		require.Equal(t, cfg.InsecureTLS, cl.Config.InsecureTLS)
		require.Equal(t, cfg.SearchBase, cl.Config.SearchBase)
		require.Equal(t, cfg.Bind, cl.Config.Bind)

		require.Equal(t, cfg.Users.IdAttribute, cl.Config.Users.IdAttribute)
		require.Equal(t, cfg.Users.SearchBase, cl.Config.Users.SearchBase)
		require.Equal(t, cfg.Users.Attributes, cl.Config.Users.Attributes)
		require.Equal(t, cfg.Users.FilterById, cl.Config.Users.FilterById)
		require.Equal(t, cfg.Users.FilterByDn, cl.Config.Users.FilterByDn)
		require.Equal(t, cfg.Users.FilterGroupsByDn, cl.Config.Users.FilterGroupsByDn)

		require.Equal(t, cfg.Groups.IdAttribute, cl.Config.Groups.IdAttribute)
		require.Equal(t, cfg.Groups.SearchBase, cl.Config.Groups.SearchBase)
		require.Equal(t, cfg.Groups.Attributes, cl.Config.Groups.Attributes)
		require.Equal(t, cfg.Groups.FilterById, cl.Config.Groups.FilterById)
		require.Equal(t, cfg.Groups.FilterByDn, cl.Config.Groups.FilterByDn)
		require.Equal(t, cfg.Groups.FilterMembersByDn, cl.Config.Groups.FilterMembersByDn)
	})
}
