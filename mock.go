package adc

import (
	"crypto/tls"
	"errors"
	"time"

	"github.com/dlampsi/generigo"
	"github.com/go-ldap/ldap/v3"
)

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

func (cl *mockClient) SetTimeout(time.Duration) {}

func (cl *mockClient) Bind(username, password string) error {
	if username == "mrError" {
		return errors.New("error for tests")
	}
	if username == "validUser" {
		if password == "validPass" {
			return nil
		}
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

func (cl *mockClient) Add(*ldap.AddRequest) error {
	return nil
}

func (cl *mockClient) Del(*ldap.DelRequest) error {
	return nil
}

func (cl *mockClient) Modify(req *ldap.ModifyRequest) error {
	if req.DN == "OU=group2,DC=company,DC=com" {
		return errors.New("error for tests")
	}
	return nil
}

func (cl *mockClient) ModifyDN(*ldap.ModifyDNRequest) error {
	return nil
}

func (cl *mockClient) Compare(dn, attribute, value string) (bool, error) {
	return true, nil
}

func (cl *mockClient) PasswordModify(*ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error) {
	return nil, nil
}

type mockDataEntry struct {
	entry   *ldap.Entry
	filters []string
}

var mockData = map[string]mockDataEntry{
	"exists_user": {
		entry: &ldap.Entry{
			DN: "OU=user1,DC=company,DC=com",
			Attributes: []*ldap.EntryAttribute{
				{Name: "sAMAccountName", Values: []string{"user1"}, ByteValues: [][]byte{}},
			},
		},
		filters: []string{
			"(&(objectClass=person)(sAMAccountName=user1))",
			"(&(objectClass=person)(distinguishedName=OU=user1,DC=company,DC=com))",
			"(&(objectCategory=person)(memberOf=OU=group1,DC=company,DC=com))",
		},
	},
	"toadd_user": {
		entry: &ldap.Entry{
			DN: "OU=user3,DC=company,DC=com",
			Attributes: []*ldap.EntryAttribute{
				{Name: "sAMAccountName", Values: []string{"user3"}, ByteValues: [][]byte{}},
			},
		},
		filters: []string{
			"(&(objectClass=person)(sAMAccountName=user3))",
			"(&(objectClass=person)(distinguishedName=OU=user3,DC=company,DC=com))",
			"(&(objectCategory=person)(memberOf=OU=group3,DC=company,DC=com))",
		},
	},
	"exists_group": {
		entry: &ldap.Entry{
			DN: "OU=group1,DC=company,DC=com",
			Attributes: []*ldap.EntryAttribute{
				{Name: "sAMAccountName", Values: []string{"group1"}, ByteValues: [][]byte{}},
			},
		},
		filters: []string{
			"(&(objectClass=group)(sAMAccountName=group1))",
			"(&(objectClass=group)(distinguishedName=OU=group1,DC=company,DC=com))",
			"(&(objectClass=group)(member=OU=user1,DC=company,DC=com))",
		},
	},
	"for_errors": {
		filters: []string{
			"(&(objectClass=person)(sAMAccountName=user2))",
			"(&(objectClass=person)(distinguishedName=OU=user2,DC=company,DC=com))",
			"(&(objectClass=group)(sAMAccountName=group2))",
			"(&(objectClass=group)(distinguishedName=OU=group2,DC=company,DC=com))",
		},
	},
}

func (cl *mockClient) Search(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
	if generigo.StringInSlice(req.Filter, mockData["exists_user"].filters) {
		result := &ldap.SearchResult{}
		result.Entries = append(result.Entries, mockData["exists_user"].entry)
		return result, nil
	}

	if generigo.StringInSlice(req.Filter, mockData["exists_group"].filters) {
		result := &ldap.SearchResult{}
		result.Entries = append(result.Entries, mockData["exists_group"].entry)
		return result, nil
	}

	if generigo.StringInSlice(req.Filter, mockData["toadd_user"].filters) {
		result := &ldap.SearchResult{}
		result.Entries = append(result.Entries, mockData["toadd_user"].entry)
		return result, nil
	}

	if generigo.StringInSlice(req.Filter, mockData["for_errors"].filters) {
		return nil, errors.New("error for tests")
	}

	return &ldap.SearchResult{Entries: nil}, nil
}

func (cl *mockClient) SearchWithPaging(searchRequest *ldap.SearchRequest, pagingSize uint32) (*ldap.SearchResult, error) {
	return nil, nil
}
