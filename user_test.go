package adc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_User_GetStringAttribute(t *testing.T) {
	t.Run("NonExists", func(t *testing.T) {
		u := &User{
			Attributes: map[string]interface{}{
				"one": "string",
			},
		}
		require.Empty(t, u.GetStringAttribute("nonexists"))
	})
	t.Run("NonString", func(t *testing.T) {
		user := &User{
			Attributes: map[string]interface{}{
				"two":   2,
				"three": []byte("bytedata"),
			},
		}
		require.Equal(t, "", user.GetStringAttribute("two"))
		require.Equal(t, "", user.GetStringAttribute("three"))
	})
	t.Run("Ok", func(t *testing.T) {
		user := &User{
			Attributes: map[string]interface{}{
				"one": "value",
			},
		}
		require.Equal(t, "value", user.GetStringAttribute("one"))
	})
}

func Test_GetUserArgs_Validate(t *testing.T) {
	t.Run("ErrWithNil", func(t *testing.T) {
		var req GetUserArgs
		require.Error(t, req.Validate())
	})
	t.Run("ErrWithEmpty", func(t *testing.T) {
		req := GetUserArgs{}
		require.Error(t, req.Validate())
	})
	t.Run("OkWithId", func(t *testing.T) {
		req := GetUserArgs{Id: "fake"}
		require.NoError(t, req.Validate())
	})
	t.Run("OkWithFilter", func(t *testing.T) {
		req := GetUserArgs{Filter: "fake"}
		require.NoError(t, req.Validate())
	})
	t.Run("OkWithDn", func(t *testing.T) {
		req := GetUserArgs{Dn: "fake"}
		require.NoError(t, req.Validate())
	})
}

func Test_Client_GetUser(t *testing.T) {
	cl := newMockClient(&Config{})
	require.NoError(t, cl.Connect())

	t.Run("BadArgs", func(t *testing.T) {
		var badArgs GetUserArgs
		_, badReqErr := cl.GetUser(badArgs)
		require.Error(t, badReqErr)
	})
	t.Run("EntryForErr", func(t *testing.T) {
		args := GetUserArgs{Id: "entryForErr", SkipGroupsSearch: true}
		u, err := cl.GetUser(args)
		require.Error(t, err)
		require.Nil(t, u)
	})
	t.Run("NonExistsUser", func(t *testing.T) {
		args := GetUserArgs{Id: "userFake", SkipGroupsSearch: true}
		user, err := cl.GetUser(args)
		require.NoError(t, err)
		require.Nil(t, user)
	})
	t.Run("TooManyEntries", func(t *testing.T) {
		user, err := cl.GetUser(GetUserArgs{
			Id:               "notUniq",
			SkipGroupsSearch: true,
			Attributes:       []string{"sAMAccountName"},
		})
		require.Error(t, err)
		require.Nil(t, user)
	})
	t.Run("ByDn", func(t *testing.T) {
		dnReq := GetUserArgs{Dn: "OU=user1,DC=company,DC=com", SkipGroupsSearch: true}
		userByDn, err := cl.GetUser(dnReq)
		require.NoError(t, err)
		require.NotNil(t, userByDn)
		require.Equal(t, dnReq.Dn, userByDn.DN)
	})
	t.Run("ByFilter", func(t *testing.T) {
		filterReq := GetUserArgs{Filter: "customFilterToSearchUser", SkipGroupsSearch: true}
		userByFilter, err := cl.GetUser(filterReq)
		require.NoError(t, err)
		require.NotNil(t, userByFilter)
		require.Equal(t, userByFilter.Id, "user1")
	})
	t.Run("Ok", func(t *testing.T) {
		args := GetUserArgs{Id: "user1", SkipGroupsSearch: true}
		user, err := cl.GetUser(args)
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
	})
}

func Test_User_IsGroupMember(t *testing.T) {
	t.Run("EmptyGroups", func(t *testing.T) {
		u := &User{}
		require.Equal(t, false, u.IsGroupMember("group1"))
	})
	t.Run("NotAMember", func(t *testing.T) {
		u := &User{
			Groups: []UserGroup{{Id: "group1"}},
		}
		require.Equal(t, false, u.IsGroupMember("group3"))
	})
	t.Run("IsMember", func(t *testing.T) {
		u := &User{
			Groups: []UserGroup{{Id: "group1"}, {Id: "group2"}},
		}
		require.Equal(t, true, u.IsGroupMember("group1"))
	})
}

func Test_User_GroupsDn(t *testing.T) {
	t.Run("EmptyGroups", func(t *testing.T) {
		u := &User{
			Groups: []UserGroup{},
		}
		require.Nil(t, u.GroupsDn())
	})
	t.Run("Ok", func(t *testing.T) {
		u := &User{
			Groups: []UserGroup{{Id: "someId", DN: "someDn"}},
		}
		require.NotNil(t, u.GroupsDn())
		require.Contains(t, u.GroupsDn(), "someDn")
	})
}

func Test_User_GroupsId(t *testing.T) {
	t.Run("EmptyGroups", func(t *testing.T) {
		u := &User{
			Groups: []UserGroup{},
		}
		require.Nil(t, u.GroupsId())
	})
	t.Run("Ok", func(t *testing.T) {
		u := &User{
			Groups: []UserGroup{{Id: "someId", DN: "someDn"}},
		}
		require.NotNil(t, u.GroupsId())
		require.Contains(t, u.GroupsId(), "someId")
	})
}
