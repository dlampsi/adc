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

func Test_populateConfig(t *testing.T) {
	defCfg := getDefaultConfig()

	t.Run("NilConfig", func(t *testing.T) {
		require.Equal(t, defCfg, populateConfig(nil))
	})

	t.Run("EmptyConfig", func(t *testing.T) {
		require.Equal(t, defCfg, populateConfig(&Config{}))
	})

	t.Run("CustomConfigPartial", func(t *testing.T) {
		customCfg := &Config{
			URL: "ldaps://fakeurl:636",
			Bind: &BindAccount{
				DN:       "some",
				Password: "fake",
			},
			SearchBase: "OU=some",
			Users: &UsersConfigs{
				SearchBase: "OU=custom-users",
			},
			Groups: &GroupsConfigs{
				SearchBase: "OU=custom-groups",
			},
		}

		cfg := populateConfig(customCfg)
		require.NotNil(t, cfg)

		require.Equal(t, defCfg.Timeout, cfg.Timeout)
		require.Equal(t, customCfg.URL, cfg.URL)
		require.Equal(t, defCfg.InsecureTLS, cfg.InsecureTLS)
		require.Equal(t, customCfg.SearchBase, cfg.SearchBase)
		require.Equal(t, customCfg.Bind, cfg.Bind)

		require.Equal(t, defCfg.Users.IdAttribute, cfg.Users.IdAttribute)
		require.Equal(t, customCfg.Users.SearchBase, cfg.Users.SearchBase)
		require.Equal(t, defCfg.Users.Attributes, cfg.Users.Attributes)
		require.Equal(t, defCfg.Users.FilterById, cfg.Users.FilterById)
		require.Equal(t, defCfg.Users.FilterByDn, cfg.Users.FilterByDn)
		require.Equal(t, defCfg.Users.FilterGroupsByDn, cfg.Users.FilterGroupsByDn)

		require.Equal(t, defCfg.Groups.IdAttribute, cfg.Groups.IdAttribute)
		require.Equal(t, customCfg.Groups.SearchBase, cfg.Groups.SearchBase)
		require.Equal(t, defCfg.Groups.Attributes, cfg.Groups.Attributes)
		require.Equal(t, defCfg.Groups.FilterById, cfg.Groups.FilterById)
		require.Equal(t, defCfg.Groups.FilterByDn, cfg.Groups.FilterByDn)
		require.Equal(t, defCfg.Groups.FilterMembersByDn, cfg.Groups.FilterMembersByDn)
	})

	t.Run("CustomConfigAll", func(t *testing.T) {
		customCfg := &Config{
			URL:         "ldaps://fakeurl:636",
			InsecureTLS: true,
			Timeout:     5 * time.Second,
			Bind: &BindAccount{
				DN:       "some",
				Password: "fake",
			},
			SearchBase: "OU=some",
			Users: &UsersConfigs{
				IdAttribute:      "custom-users-id-attr",
				Attributes:       []string{"dummy-user-attr"},
				SearchBase:       "OU=custom-users",
				FilterById:       "customFilterById",
				FilterByDn:       "customFilterByDn",
				FilterGroupsByDn: "customFilterGroupsByDn",
			},
			Groups: &GroupsConfigs{
				IdAttribute:       "custom-groups-id-attr",
				Attributes:        []string{"dummy-group-attr"},
				SearchBase:        "OU=custom-groups",
				FilterById:        "customFilterById",
				FilterByDn:        "customFilterByDn",
				FilterMembersByDn: "customFilterMembersByDn",
			},
		}

		cfg := populateConfig(customCfg)
		require.NotNil(t, cfg)

		require.Equal(t, customCfg.Timeout, cfg.Timeout)
		require.Equal(t, customCfg.URL, cfg.URL)
		require.Equal(t, customCfg.InsecureTLS, cfg.InsecureTLS)
		require.Equal(t, customCfg.SearchBase, cfg.SearchBase)
		require.Equal(t, customCfg.Bind, cfg.Bind)

		require.Equal(t, customCfg.Users.IdAttribute, cfg.Users.IdAttribute)
		require.Equal(t, customCfg.Users.SearchBase, cfg.Users.SearchBase)
		require.Equal(t, customCfg.Users.Attributes, cfg.Users.Attributes)
		require.Equal(t, customCfg.Users.FilterById, cfg.Users.FilterById)
		require.Equal(t, customCfg.Users.FilterByDn, cfg.Users.FilterByDn)
		require.Equal(t, customCfg.Users.FilterGroupsByDn, cfg.Users.FilterGroupsByDn)

		require.Equal(t, customCfg.Groups.IdAttribute, cfg.Groups.IdAttribute)
		require.Equal(t, customCfg.Groups.SearchBase, cfg.Groups.SearchBase)
		require.Equal(t, customCfg.Groups.Attributes, cfg.Groups.Attributes)
		require.Equal(t, customCfg.Groups.FilterById, cfg.Groups.FilterById)
		require.Equal(t, customCfg.Groups.FilterByDn, cfg.Groups.FilterByDn)
		require.Equal(t, customCfg.Groups.FilterMembersByDn, cfg.Groups.FilterMembersByDn)
	})
}
