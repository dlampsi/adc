package adc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_User_GetStringAttribute(t *testing.T) {
	u := &User{
		Attributes: map[string]interface{}{
			"one":   "string",
			"two":   2,
			"three": []byte("bytedata"),
		},
	}
	require.NotEmpty(t, u.GetStringAttribute("one"))
	require.Equal(t, "string", u.GetStringAttribute("one"))
	require.Empty(t, u.GetStringAttribute("two"))
	require.Empty(t, u.GetStringAttribute("three"))
	require.Empty(t, u.GetStringAttribute("nonexists"))
}

func Test_GetUserArgs_Validate(t *testing.T) {
	var req GetUserArgs
	err := req.Validate()
	require.Error(t, err)

	req = GetUserArgs{}
	err1 := req.Validate()
	require.Error(t, err1)

	req = GetUserArgs{Id: "fake"}
	errOk := req.Validate()
	require.NoError(t, errOk)

	req = GetUserArgs{Filter: "fake"}
	errOk = req.Validate()
	require.NoError(t, errOk)

	req = GetUserArgs{Dn: "fake"}
	errOk = req.Validate()
	require.NoError(t, errOk)
}

func Test_Client_GetUser(t *testing.T) {
	cl := New(&Config{}, withMock())
	err := cl.Connect()
	require.NoError(t, err)

	var badArgs GetUserArgs
	_, badReqErr := cl.GetUser(badArgs)
	require.Error(t, badReqErr)

	args := GetUserArgs{Id: "entryForErr", SkipGroupsSearch: true}
	_, err = cl.GetUser(args)
	require.Error(t, err)

	args = GetUserArgs{Id: "userFake", SkipGroupsSearch: true}
	user, err := cl.GetUser(args)
	require.NoError(t, err)
	require.Nil(t, user)

	// Too many entries error
	user, err = cl.GetUser(GetUserArgs{
		Id:               "notUniq",
		SkipGroupsSearch: true,
		Attributes:       []string{"sAMAccountName"},
	})
	require.Error(t, err)
	require.Nil(t, user)

	dnReq := GetUserArgs{Dn: "OU=user1,DC=company,DC=com", SkipGroupsSearch: true}
	userByDn, err := cl.GetUser(dnReq)
	require.NoError(t, err)
	require.NotNil(t, userByDn)
	require.Equal(t, dnReq.Dn, userByDn.DN)

	filterReq := GetUserArgs{Filter: "customFilterToSearchUser", SkipGroupsSearch: true}
	userByFilter, err := cl.GetUser(filterReq)
	require.NoError(t, err)
	require.NotNil(t, userByFilter)
	require.Equal(t, userByFilter.Id, "user1")

	args = GetUserArgs{Id: "user1", SkipGroupsSearch: true}
	user, err = cl.GetUser(args)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, args.Id, user.Id)
	require.Nil(t, user.Groups)

	args.Attributes = []string{"something"}
	user, err = cl.GetUser(args)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, args.Id, user.Id)
	require.Nil(t, user.Groups)

	args.SkipGroupsSearch = false
	user, err = cl.GetUser(args)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, args.Id, user.Id)
	require.NotNil(t, user.Groups)
	require.Len(t, user.Groups, 1)
}

func Test_User_IsGroupMember(t *testing.T) {
	u := &User{}
	require.Equal(t, false, u.IsGroupMember("group1"))

	u.Groups = []UserGroup{
		{Id: "group1"},
		{Id: "group2"},
	}
	require.Equal(t, false, u.IsGroupMember("group3"))
	require.Equal(t, true, u.IsGroupMember("group1"))
	require.Equal(t, true, u.IsGroupMember("group2"))
}

func Test_User_GroupsDn(t *testing.T) {
	u := &User{
		Groups: []UserGroup{},
	}
	require.Nil(t, u.GroupsDn())

	newGroup := UserGroup{Id: "someId", DN: "someDn"}
	u.Groups = append(u.Groups, newGroup)
	require.NotNil(t, u.GroupsDn())
	require.Contains(t, u.GroupsDn(), newGroup.DN)
}

func Test_User_GroupsId(t *testing.T) {
	u := &User{
		Groups: []UserGroup{},
	}

	require.Nil(t, u.GroupsId())
	newGroup := UserGroup{Id: "someId", DN: "someDn"}
	u.Groups = append(u.Groups, newGroup)
	require.NotNil(t, u.GroupsId())
	require.Contains(t, u.GroupsId(), newGroup.Id)
}
