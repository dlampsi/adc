// Package adc provides basic client library for Active Directory.
package adc

import (
	"crypto/tls"
	"errors"
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
	logger Logger
}

// Creates new client and populate provided config and options.
func New(cfg *Config, opts ...Option) *Client {
	cl := &Client{
		cfg: &Config{
			Timeout: defaultClientTimeout,
			Users:   DefaultUsersConfigs(),
			Groups:  DefaultGroupsConfigs(),
		},
		logger: &nopLogger{},
	}

	// Apply options
	for _, opt := range opts {
		opt(cl)
	}

	// Populate optional config
	cl.popConfig(cfg)

	return cl
}

// Client logger interface.
type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
}

type Option func(*Client)

// Specifies ldap client for AD client.
func WithLdapClient(l ldap.Client) Option {
	return func(cl *Client) { cl.ldapCl = l }
}

// Specifies custom logger for client.
func WithLogger(l Logger) Option {
	return func(cl *Client) { cl.logger = l }
}

func (cl *Client) Config() *Config {
	return cl.cfg
}

// Connects to AD server and store connection into client.
func (cl *Client) Connect() error {
	conn, err := cl.connect(cl.cfg.Bind)
	if err != nil {
		return err
	}
	cl.ldapCl = conn
	return nil
}

// Connects and bind to LDAP server by provided bind account.
func (cl *Client) connect(bind *BindAccount) (ldap.Client, error) {
	ldapCl := cl.ldapCl

	// Use default ldap module connection if no ldap client provided in client
	if ldapCl == nil {
		conn, err := cl.dialLdap()
		if err != nil {
			return nil, err
		}
		ldapCl = conn
	}

	if bind != nil {
		if err := ldapCl.Bind(bind.DN, bind.Password); err != nil {
			return nil, err
		}
	}

	return ldapCl, nil
}

// Dials ldap server provided in client configuration.
func (cl *Client) dialLdap() (ldap.Client, error) {
	var opts []ldap.DialOpt
	if strings.HasPrefix("ldaps://", cl.cfg.URL) {
		opts = append(opts, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: cl.cfg.InsecureTLS}))
	}
	return ldap.DialURL(cl.cfg.URL, opts...)
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
	conn, err := cl.dialLdap()
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.Bind(dn, password); err != nil {
		return err
	}
	return nil
}
