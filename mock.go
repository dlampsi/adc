package adc

import (
	"crypto/tls"
	"errors"
	"time"

	"github.com/dlampsi/generigo"
	"github.com/go-ldap/ldap/v3"
)

// Entry attribute name, that helps match entry to provided request.
const mockFiltersAttribute = "filtersToFind"

// Data for mock ldap provider.
var mockEntriesData = mockEntries{
	"user1": &ldap.Entry{
		DN: "OU=user1,DC=company,DC=com",
		Attributes: []*ldap.EntryAttribute{
			{Name: "sAMAccountName", Values: []string{"user1"}},
			{Name: mockFiltersAttribute, Values: []string{
				"(&(objectClass=person)(sAMAccountName=user1))",
				"(&(objectClass=person)(distinguishedName=OU=user1,DC=company,DC=com))",
				"(&(objectCategory=person)(memberOf=OU=group1,DC=company,DC=com))",
				"customFilterToSearchUser",
			}},
		},
	},
	"user2": &ldap.Entry{
		DN: "OU=user2,DC=company,DC=com",
		Attributes: []*ldap.EntryAttribute{
			{Name: "sAMAccountName", Values: []string{"user2"}},
			{Name: mockFiltersAttribute, Values: []string{
				"(&(objectClass=person)(sAMAccountName=user2))",
				"(&(objectClass=person)(distinguishedName=OU=user2,DC=company,DC=com))",
				"(&(objectCategory=person)(memberOf=OU=group2,DC=company,DC=com))",
			}},
		},
	},
	"userToAdd": &ldap.Entry{
		DN: "OU=userToAdd,DC=company,DC=com",
		Attributes: []*ldap.EntryAttribute{
			{Name: "sAMAccountName", Values: []string{"userToAdd"}},
			{Name: mockFiltersAttribute, Values: []string{
				"(&(objectClass=person)(sAMAccountName=userToAdd))",
				"(&(objectClass=person)(distinguishedName=OU=userToAdd,DC=company,DC=com))",
				"(&(objectCategory=person)(memberOf=OU=group2,DC=company,DC=com))",
			}},
		},
	},
	"group1": &ldap.Entry{
		DN: "OU=group1,DC=company,DC=com",
		Attributes: []*ldap.EntryAttribute{
			{Name: "sAMAccountName", Values: []string{"group1"}},
			{Name: mockFiltersAttribute, Values: []string{
				"(&(objectClass=group)(sAMAccountName=group1))",
				"(&(objectClass=group)(distinguishedName=OU=group1,DC=company,DC=com))",
				"(&(objectClass=group)(member=OU=user1,DC=company,DC=com))",
				"customFilterToSearchGroup",
			}},
		},
	},
	"group2": &ldap.Entry{
		DN: "OU=group2,DC=company,DC=com",
		Attributes: []*ldap.EntryAttribute{
			{Name: "sAMAccountName", Values: []string{"group2"}},
			{Name: mockFiltersAttribute, Values: []string{
				"(&(objectClass=group)(sAMAccountName=group2))",
				"(&(objectClass=group)(distinguishedName=OU=group2,DC=company,DC=com))",
				"(&(objectClass=group)(member=OU=user2,DC=company,DC=com))",
			}},
		},
	},
	"entryForErr": &ldap.Entry{
		DN: "OU=entryForErr,DC=company,DC=com",
		Attributes: []*ldap.EntryAttribute{
			{Name: "sAMAccountName", Values: []string{"entryForErr"}},
			{Name: mockFiltersAttribute, Values: []string{
				"(&(objectClass=person)(sAMAccountName=entryForErr))",
				"(&(objectClass=person)(distinguishedName=OU=entryForErr,DC=company,DC=com))",
				"(&(objectClass=group)(sAMAccountName=entryForErr))",
				"(&(objectClass=group)(distinguishedName=OU=entryForErr,DC=company,DC=com))",
				"(&(objectCategory=person)(memberOf=OU=groupWithErrMember,DC=company,DC=com))",
			}},
		},
	},
	"groupWithErrMember": &ldap.Entry{
		DN: "OU=groupWithErrMember,DC=company,DC=com",
		Attributes: []*ldap.EntryAttribute{
			{Name: "sAMAccountName", Values: []string{"groupWithErrMember"}},
			{Name: mockFiltersAttribute, Values: []string{
				"(&(objectClass=group)(sAMAccountName=groupWithErrMember))",
				"(&(objectClass=group)(distinguishedName=OU=groupWithErrMember,DC=company,DC=com))",
				"(&(objectClass=group)(member=OU=entryForErr,DC=company,DC=com))",
			}},
		},
	},
	"userToReconnect": &ldap.Entry{
		DN: "OU=userToReconnect,DC=company,DC=com",
		Attributes: []*ldap.EntryAttribute{
			{Name: "sAMAccountName", Values: []string{"userToReconnect"}},
			{Name: mockFiltersAttribute, Values: []string{
				"(&(objectClass=person)(sAMAccountName=userToReconnect))",
				"(&(objectClass=person)(distinguishedName=OU=userToReconnect,DC=company,DC=com))",
			}},
		},
	},
	"notUniq1": &ldap.Entry{
		DN: "OU=notUniq,DC=company,DC=com",
		Attributes: []*ldap.EntryAttribute{
			{Name: "sAMAccountName", Values: []string{"notUniq"}},
			{Name: mockFiltersAttribute, Values: []string{
				"(&(objectClass=person)(sAMAccountName=notUniq))",
				"(&(objectClass=person)(distinguishedName=OU=notUniq,DC=company,DC=com))",
				"(&(objectClass=group)(sAMAccountName=notUniq))",
				"(&(objectClass=group)(distinguishedName=OU=notUniq,DC=company,DC=com))",
			}},
		},
	},
	"notUniq2": &ldap.Entry{
		DN: "OU=notUniq,DC=company,DC=com",
		Attributes: []*ldap.EntryAttribute{
			{Name: "sAMAccountName", Values: []string{"notUniq"}},
			{Name: mockFiltersAttribute, Values: []string{
				"(&(objectClass=person)(sAMAccountName=notUniq))",
				"(&(objectClass=person)(distinguishedName=OU=notUniq,DC=company,DC=com))",
				"(&(objectClass=group)(sAMAccountName=notUniq))",
				"(&(objectClass=group)(distinguishedName=OU=notUniq,DC=company,DC=com))",
			}},
		},
	},
}

type mockEntries map[string]*ldap.Entry

func (me mockEntries) getEntryByDn(dn string) *ldap.Entry {
	for _, entry := range me {
		if entry.DN == dn {
			return entry
		}
	}
	return nil
}

func (me mockEntries) getEntriesByFilter(filter string) ([]*ldap.Entry, error) {
	var result []*ldap.Entry
	for id, entry := range me {
		filters := entry.GetAttributeValues(mockFiltersAttribute)
		if generigo.StringInSlice(filter, filters) {
			if id == "entryForErr" {
				return nil, errors.New("error for tests")
			}
			if id == "userToReconnect" {
				return nil, ldap.NewError(200, errors.New("connection error"))
			}
			result = append(result, entry)
		}
	}
	return result, nil
}

// Dummy not operational logger.
type nopLogger struct{}

func (l *nopLogger) Debug(args ...interface{})                   {}
func (l *nopLogger) Debugf(template string, args ...interface{}) {}

// Mock client. Implements ldap client interface.
type mockClient struct {
}

// Validates interfaces compliance.
var _ ldap.Client = (*mockClient)(nil)

func (cl *mockClient) Start() {
}

func (cl *mockClient) StartTLS(*tls.Config) error {
	return nil
}

func (cl *mockClient) Close() {}

func (cl *mockClient) IsClosing() bool {
	return false
}

func (cl *mockClient) SetTimeout(time.Duration) {}

func (cl *mockClient) TLSConnectionState() (tls.ConnectionState, bool) {
	return tls.ConnectionState{}, true
}

var (
	validMockBind     = &BindAccount{DN: "validUser", Password: "validPass"}
	invalidMockBind   = &BindAccount{DN: "mrError", Password: "mrErrorPass"}
	reconnectMockBind = &BindAccount{DN: "OU=userToReconnect,DC=company,DC=com", Password: "validPass"}
)

func (cl *mockClient) Bind(username, password string) error {
	if username == invalidMockBind.DN {
		return errors.New("error for tests")
	}
	if username == validMockBind.DN && password == validMockBind.Password {
		return nil
	}
	return errors.New("unauthorised")
}

func (cl *mockClient) UnauthenticatedBind(username string) error {
	return nil
}

func (cl *mockClient) SimpleBind(*ldap.SimpleBindRequest) (*ldap.SimpleBindResult, error) {
	return nil, nil
}

func (cl *mockClient) ExternalBind() error {
	return nil
}

func (cl *mockClient) NTLMUnauthenticatedBind(domain, username string) error {
	return nil
}

func (cl *mockClient) Unbind() error {
	return nil
}

func (cl *mockClient) Add(*ldap.AddRequest) error {
	return nil
}

func (cl *mockClient) Del(*ldap.DelRequest) error {
	return nil
}

func (cl *mockClient) Modify(req *ldap.ModifyRequest) error {
	entry := mockEntriesData.getEntryByDn(req.DN)
	if entry == nil {
		return errors.New("entry not found")
	}
	if entry.DN == mockEntriesData["entryForErr"].DN {
		return errors.New("error for tests")
	}
	return nil
}

func (cl *mockClient) ModifyDN(*ldap.ModifyDNRequest) error {
	return nil
}

func (cl *mockClient) ModifyWithResult(*ldap.ModifyRequest) (*ldap.ModifyResult, error) {
	return nil, nil
}

func (cl *mockClient) Compare(dn, attribute, value string) (bool, error) {
	return true, nil
}

func (cl *mockClient) PasswordModify(*ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error) {
	return nil, nil
}

func (cl *mockClient) Search(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
	entries, err := mockEntriesData.getEntriesByFilter(req.Filter)
	if err != nil {
		return nil, err
	}
	return &ldap.SearchResult{Entries: entries}, nil
}

func (cl *mockClient) SearchWithPaging(searchRequest *ldap.SearchRequest, pagingSize uint32) (*ldap.SearchResult, error) {
	return nil, nil
}
