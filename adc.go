// Package adc provides basic client library for Active Directory.
package adc

import (
	"context"
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
	cfg     *Config
	ldapCl  ldap.Client
	logger  Logger
	useMock bool
}

// Creates new client and populate provided config and options.
func New(cfg *Config, opts ...Option) *Client {
	cl := &Client{
		cfg: &Config{
			Timeout: defaultClientTimeout,
			Users:   DefaultUsersConfigs(),
			Groups:  DefaultGroupsConfigs(),
		},
		logger: newNopLogger(),
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

// Specifies custom logger for client.
func WithLogger(l Logger) Option {
	return func(cl *Client) { cl.logger = l }
}

// Enables mock ldap interface for client
func withMock() Option {
	return func(cl *Client) { cl.useMock = true }
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
	conn, err := cl.dial()
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

// Dials ldap server provided in client configuration.
func (cl *Client) dial() (ldap.Client, error) {
	if cl.useMock {
		return &mockClient{}, nil
	}
	var opts []ldap.DialOpt
	if strings.HasPrefix(cl.cfg.URL, "ldaps://") {
		opts = append(opts, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: cl.cfg.InsecureTLS}))
	}
	return ldap.DialURL(cl.cfg.URL, opts...)
}

// Closes connection to AD.
func (cl *Client) Disconnect() error {
	if cl.ldapCl == nil {
		return nil
	}
	return cl.ldapCl.Close()
}

// Checks connections to AD and tries to reconnect if the connection is lost.
func (cl *Client) Reconnect(ctx context.Context, tickerDuration time.Duration, maxAttempts int) error {
	_, connErr := cl.searchEntry(&ldap.SearchRequest{
		BaseDN:       cl.cfg.SearchBase,
		Scope:        ldap.ScopeWholeSubtree,
		DerefAliases: ldap.NeverDerefAliases,
		TimeLimit:    int(cl.cfg.Timeout.Seconds()),
		Filter:       fmt.Sprintf(cl.cfg.Users.FilterByDn, ldap.EscapeFilter(cl.cfg.Bind.DN)),
		Attributes:   []string{cl.cfg.Users.IdAttribute},
	})
	if connErr == nil {
		return nil
	}

	if tickerDuration == 0 {
		tickerDuration = 5 * time.Second
	}
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	if maxAttempts == 0 {
		maxAttempts = 2
	}

	attempt := 0
	for {
		select {
		case <-ticker.C:
			if attempt >= maxAttempts {
				return fmt.Errorf("failed after '%d' attempts. error: %s", attempt, connErr)
			}
			attempt++
			cl.logger.Debugf("Reconnecting to AD server. Attempt: %d", attempt)

			if err := cl.Disconnect(); err != nil {
				return fmt.Errorf("failed to disconnect from the server: %w", err)
			}

			if err := cl.Connect(); err == nil {
				cl.logger.Debug("Successfully reconneted to AD server")
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
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
	conn, err := cl.dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.Bind(dn, password); err != nil {
		return err
	}
	return nil
}
