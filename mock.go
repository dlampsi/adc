package adc

import (
	"crypto/tls"
	"errors"
	"time"

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
	return nil
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

func (cl *mockClient) Modify(*ldap.ModifyRequest) error {
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

const (
	existsUserFilter         = "(&(objectClass=person)(sAMAccountName=user1))"
	errorUserFilter          = "(&(objectClass=person)(sAMAccountName=user2))"
	existsUserGroupsFilter   = "(&(objectClass=group)(member=OU=user1,DC=company,DC=com))"
	existsGroupFilter        = "(&(objectClass=group)(sAMAccountName=group1))"
	errorGroupFilter         = "(&(objectClass=group)(sAMAccountName=group2))"
	existsGroupMembersFilter = "(&(objectCategory=person)(memberOf=OU=group1,DC=company,DC=com))"
)

var (
	existsUserEntry = &ldap.Entry{
		DN: "OU=user1,DC=company,DC=com",
		Attributes: []*ldap.EntryAttribute{
			{Name: "sAMAccountName", Values: []string{"user1"}, ByteValues: [][]byte{}},
		},
	}
	existsGroupEntry = &ldap.Entry{
		DN: "OU=group1,DC=company,DC=com",
		Attributes: []*ldap.EntryAttribute{
			{Name: "sAMAccountName", Values: []string{"group1"}, ByteValues: [][]byte{}},
		},
	}
)

func (cl *mockClient) Search(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
	if req.Filter == existsUserFilter || req.Filter == existsGroupMembersFilter {
		result := &ldap.SearchResult{}
		result.Entries = append(result.Entries, existsUserEntry)
		return result, nil
	}
	if req.Filter == errorUserFilter {
		return nil, errors.New("error for tests")
	}

	if req.Filter == existsGroupFilter || req.Filter == existsUserGroupsFilter {
		result := &ldap.SearchResult{}
		result.Entries = append(result.Entries, existsGroupEntry)
		return result, nil
	}
	if req.Filter == errorGroupFilter {
		return nil, errors.New("error for tests")
	}

	return nil, nil
}

func (cl *mockClient) SearchWithPaging(searchRequest *ldap.SearchRequest, pagingSize uint32) (*ldap.SearchResult, error) {
	return nil, nil
}
