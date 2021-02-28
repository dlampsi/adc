// Package adc provides basic client library for Active Directory.
package adc

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
)

const (
	defaultClientTimeout = 10 * time.Second
)

// Active Direcotry client.
type Client struct {
	cfg    *Config
	ldapCl ldap.Client
}

// Creates new client and populate provided config and options.
func New(cfg *Config, opts ...Option) *Client {
	cl := &Client{
		cfg: &Config{
			Timeout: defaultClientTimeout,
			Users:   DefaultUsersConfigs(),
			Groups:  DefaultGroupsConfigs(),
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(cl)
	}

	// Populate optional config
	cl.popConfig(cfg)

	return cl
}

type Option func(*Client)

// Specifies ldap client for AD client.
func WithLdapClient(l ldap.Client) Option {
	return func(cl *Client) { cl.ldapCl = l }
}

func (cl *Client) Config() *Config {
	return cl.cfg
}

// Connects to AD server and store connection into client.
func (cl *Client) Connect() error {
	if cl.ldapCl != nil {
		return nil
	}
	conn, err := cl.connect(cl.cfg.Bind)
	if err != nil {
		return fmt.Errorf("can't connect: %s", err.Error())
	}
	cl.ldapCl = conn
	return nil
}

// Connects and bind to LDAP server by provided bind account.
func (cl *Client) connect(bind *BindAccount) (*ldap.Conn, error) {
	var opts []ldap.DialOpt
	if strings.HasPrefix("ldaps://", cl.cfg.URL) {
		opts = append(opts, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: cl.cfg.InsecureTLS}))
	}
	conn, err := ldap.DialURL(cl.cfg.URL, opts...)
	if err != nil {
		return nil, err
	}
	if bind != nil {
		if err := conn.Bind(bind.DN, bind.Password); err != nil {
			return nil, err
		}
	}
	return conn, nil
}

// Closes connection to AD.
func (cl *Client) Disconnect() {
	if cl.ldapCl != nil {
		cl.ldapCl.Close()
	}
}

// SearchEntry Perfrom search for single ldap entry.
// Returns nil if no entries found.
// Returns 'ErrTooManyEntriesFound' error if entries more that one.
func (cl *Client) searchEntry(req *ldap.SearchRequest) (*ldap.Entry, error) {
	result, err := cl.ldapCl.Search(req)
	if err != nil {
		return nil, err
	}
	if len(result.Entries) > 1 {
		return nil, errors.New("too many entries found")
	}
	if len(result.Entries) < 1 {
		return nil, nil
	}
	return result.Entries[0], nil
}

// SearchEntries Perfroms search for ldap entries.
func (cl *Client) searchEntries(req *ldap.SearchRequest) ([]*ldap.Entry, error) {
	result, err := cl.ldapCl.Search(req)
	if err != nil {
		return nil, err
	}
	return result.Entries, nil
}

// Performs update for provided entry attribure by entry DN.
func (cl *Client) updateAttribute(dn string, attribute string, values []string) error {
	mr := ldap.NewModifyRequest(dn, nil)
	mr.Replace(attribute, values)
	return cl.ldapCl.Modify(mr)
}

// Tries to authorise in AcitveDirecotry by provided DN and password and return error if failed.
func (cl *Client) CheckAuthByDN(dn, password string) error {
	bind := &BindAccount{DN: dn, Password: password}
	conn, err := cl.connect(bind)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}
