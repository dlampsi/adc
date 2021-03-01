package adc

import (
	"fmt"
	"time"
)

type Config struct {
	// LDAP server URL. Examle 'ldaps://cl.local:636'
	URL string `json:"url"`
	// Use insecure SSL connection.
	InsecureTLS bool `json:"insecure_tls"`
	// Time limit for requests.
	Timeout time.Duration
	// Base OU for search requests.
	SearchBase string `json:"search_base"`

	// Bind account info.
	Bind *BindAccount `json:"bind"`

	// Requests filters vars.
	Users *UsersConfigs `json:"users"`
	// Requests filters vars.
	Groups *GroupsConfigs `json:"groups"`
}

// Account attributes to authentificate in AD.
type BindAccount struct {
	DN       string `json:"dn"`
	Password string `json:"password"`
}

type UsersConfigs struct {
	// The ID attribute name for group.
	IdAttribute string `json:"id_attribute"`
	// User attributes for fetch from AD.
	Attributes []string `json:"attributes"`
	// Base OU to search users requests. Sets to Config.SearchBase if not provided.
	SearchBase string `json:"search_base"`
	// LDAP filter to get user by ID.
	FilterById string `json:"filter_by_id"`
	// LDAP filter to get user by DN.
	FilterByDn string `json:"filter_by_dn"`
	// LDAP filter to get user groups membership.
	FilterGroupsByDn string `json:"filter_groups_by_dn"`
}

func DefaultUsersConfigs() *UsersConfigs {
	const defaultIdAttribute = "sAMAccountName"
	return &UsersConfigs{
		IdAttribute:      defaultIdAttribute,
		Attributes:       []string{"givenName", "sn", "mail"},
		FilterById:       fmt.Sprintf("(&(objectClass=person)(%s=%%v))", defaultIdAttribute),
		FilterByDn:       "(&(objectClass=person)(distinguishedName=%v))",
		FilterGroupsByDn: "(&(objectClass=group)(member=%v))",
	}
}

type GroupsConfigs struct {
	// The ID attribute name for group.
	IdAttribute string `json:"id_attribute"`
	// Group attributes for fetch from AD.
	Attributes []string `json:"attributes"`
	// Base OU to search groups requests. Sets to Config.SearchBase if not provided.
	SearchBase string `json:"search_base"`
	// LDAP filter to get group by ID.
	FilterById string `json:"filter_by_id"`
	// LDAP filter to get group by DN.
	FilterByDn string `json:"filter_by_dn"`
	// LDAP filter to get group members.
	FilterMembersByDn string `json:"filter_members_by_dn"`
}

func DefaultGroupsConfigs() *GroupsConfigs {
	const defaultIdAttribute = "sAMAccountName"
	return &GroupsConfigs{
		IdAttribute:       defaultIdAttribute,
		Attributes:        []string{"cn", "description"},
		FilterById:        fmt.Sprintf("(&(objectClass=group)(%s=%%v))", defaultIdAttribute),
		FilterByDn:        "(&(objectClass=group)(distinguishedName=%v))",
		FilterMembersByDn: "(&(objectCategory=person)(memberOf=%v))",
	}
}

// Appends attributes to params in client config file.
func (cfg *Config) AppendUsesAttributes(attrs ...string) {
	cfg.Users.Attributes = append(cfg.Users.Attributes, attrs...)
}

// Appends attributes to params in client config file.
func (cfg *Config) AppendGroupsAttributes(attrs ...string) {
	cfg.Groups.Attributes = append(cfg.Groups.Attributes, attrs...)
}

// Populates client config by provided config struct.
func (cl *Client) popConfig(cfg *Config) {
	if cfg == nil {
		return
	}
	if cfg.URL != "" {
		cl.cfg.URL = cfg.URL
	}
	cl.cfg.InsecureTLS = cfg.InsecureTLS
	if cfg.Timeout != 0 {
		cl.cfg.Timeout = cfg.Timeout
	}
	if cfg.SearchBase != "" {
		cl.cfg.SearchBase = cfg.SearchBase
	}
	if cfg.Bind != nil {
		cl.cfg.Bind = cfg.Bind
	}

	cl.cfg.Users.SearchBase = cl.cfg.SearchBase
	cl.cfg.Groups.SearchBase = cl.cfg.SearchBase

	if cfg.Users != nil {
		if cfg.Users.Attributes != nil {
			cl.cfg.Users.Attributes = cfg.Users.Attributes
		}
		if cfg.Users.SearchBase != "" {
			cl.cfg.Users.SearchBase = cfg.Users.SearchBase
		}
		if cfg.Users.FilterById != "" {
			cl.cfg.Users.FilterById = cfg.Users.FilterById
		}
		if cfg.Users.FilterByDn != "" {
			cl.cfg.Users.FilterByDn = cfg.Users.FilterByDn
		}
		if cfg.Users.FilterGroupsByDn != "" {
			cl.cfg.Users.FilterGroupsByDn = cfg.Users.FilterGroupsByDn
		}
	}
	if cfg.Groups != nil {
		if cfg.Groups.Attributes != nil {
			cl.cfg.Groups.Attributes = cfg.Groups.Attributes
		}
		if cfg.Groups.SearchBase != "" {
			cl.cfg.Groups.SearchBase = cfg.Groups.SearchBase
		}
		if cfg.Groups.FilterById != "" {
			cl.cfg.Groups.FilterById = cfg.Groups.FilterById
		}
		if cfg.Groups.FilterByDn != "" {
			cl.cfg.Groups.FilterByDn = cfg.Groups.FilterByDn
		}
		if cfg.Groups.FilterMembersByDn != "" {
			cl.cfg.Groups.FilterMembersByDn = cfg.Groups.FilterMembersByDn
		}
	}
}
