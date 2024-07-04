package adctests

import (
	"fmt"
	"testing"

	"github.com/dlampsi/adc"
	"github.com/stretchr/testify/require"
)

func Test_User_GetStringAttribute(t *testing.T) {
	t.Run("NonExists", func(t *testing.T) {
		u := &adc.User{
			Attributes: map[string]interface{}{
				"one": "string",
			},
		}
		require.Empty(t, u.GetStringAttribute("nonexists"))
	})
	t.Run("NonString", func(t *testing.T) {
		user := &adc.User{
			Attributes: map[string]interface{}{
				"two":   2,
				"three": []byte("bytedata"),
			},
		}
		require.Equal(t, "", user.GetStringAttribute("two"))
		require.Equal(t, "", user.GetStringAttribute("three"))
	})
	t.Run("Ok", func(t *testing.T) {
		user := &adc.User{
			Attributes: map[string]interface{}{
				"one": "value",
			},
		}
		require.Equal(t, "value", user.GetStringAttribute("one"))
	})
}

func Test_GetUserArgs_Validate(t *testing.T) {
	t.Run("ErrWithNil", func(t *testing.T) {
		var req adc.GetUserArgs
		require.Error(t, req.Validate())
	})
	t.Run("ErrWithEmpty", func(t *testing.T) {
		req := adc.GetUserArgs{}
		require.Error(t, req.Validate())
	})
	t.Run("OkWithId", func(t *testing.T) {
		req := adc.GetUserArgs{Id: "fake"}
		require.NoError(t, req.Validate())
	})
	t.Run("OkWithFilter", func(t *testing.T) {
		req := adc.GetUserArgs{Filter: "fake"}
		require.NoError(t, req.Validate())
	})
	t.Run("OkWithDn", func(t *testing.T) {
		req := adc.GetUserArgs{Dn: "fake"}
		require.NoError(t, req.Validate())
	})
}

func Test_User_IsGroupMember(t *testing.T) {
	t.Run("EmptyGroups", func(t *testing.T) {
		u := &adc.User{}
		require.Equal(t, false, u.IsGroupMember("group1"))
	})
	t.Run("NotAMember", func(t *testing.T) {
		u := &adc.User{
			Groups: []adc.UserGroup{{Id: "group1"}},
		}
		require.Equal(t, false, u.IsGroupMember("group3"))
	})
	t.Run("IsMember", func(t *testing.T) {
		u := &adc.User{
			Groups: []adc.UserGroup{{Id: "group1"}, {Id: "group2"}},
		}
		require.Equal(t, true, u.IsGroupMember("group1"))
	})
}

func Test_User_GroupsDn(t *testing.T) {
	t.Run("EmptyGroups", func(t *testing.T) {
		u := &adc.User{
			Groups: []adc.UserGroup{},
		}
		require.Nil(t, u.GroupsDn())
	})
	t.Run("Ok", func(t *testing.T) {
		u := &adc.User{
			Groups: []adc.UserGroup{{Id: "someId", DN: "someDn"}},
		}
		require.NotNil(t, u.GroupsDn())
		require.Contains(t, u.GroupsDn(), "someDn")
	})
}

func Test_User_GroupsId(t *testing.T) {
	t.Run("EmptyGroups", func(t *testing.T) {
		u := &adc.User{
			Groups: []adc.UserGroup{},
		}
		require.Nil(t, u.GroupsId())
	})
	t.Run("Ok", func(t *testing.T) {
		u := &adc.User{
			Groups: []adc.UserGroup{{Id: "someId", DN: "someDn"}},
		}
		require.NotNil(t, u.GroupsId())
		require.Contains(t, u.GroupsId(), "someId")
	})
}

func Test_Client_GetUser(t *testing.T) {
	cfg := getClientConfig()
	cl := adc.New(&cfg, adc.WithLogger(&logger{t: t}))
	require.NoError(t, cl.Connect())

	t.Run("BadArgs", func(t *testing.T) {
		var req adc.GetUserArgs
		user, err := cl.GetUser(req)
		require.Error(t, err, "Expected error on empty request")
		require.Nil(t, user, "User should be nil on error")
	})
	t.Run("Non exists", func(t *testing.T) {
		getUserReq := adc.GetUserArgs{
			Id: "nonexists",
		}
		user, err := cl.GetUser(getUserReq)
		require.NoError(t, err, "Expected no error on non exists user")
		require.Nil(t, user, "Non exists user error should return nil")
	})
	t.Run("TooManyEntries", func(t *testing.T) {
		getUserReq := adc.GetUserArgs{
			Filter: "(&(objectClass=user))",
		}
		user, err := cl.GetUser(getUserReq)
		require.Error(t, err, "Expected error on too many entries")
		require.Nil(t, user, "User should be nil on error")
	})
	t.Run("OkById", func(t *testing.T) {
		getUserReq := adc.GetUserArgs{
			Id: "testuser1",
		}
		user, err := cl.GetUser(getUserReq)
		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, getUserReq.Id, user.Id)
		require.NotEmpty(t, user.Groups)
	})
	t.Run("OkByIdWithoutGroups", func(t *testing.T) {
		getUserReq := adc.GetUserArgs{
			Id:               "testuser1",
			SkipGroupsSearch: true,
		}
		user, err := cl.GetUser(getUserReq)
		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, getUserReq.Id, user.Id)
		require.Empty(t, user.Groups)
		fmt.Println(user)
	})
	t.Run("OkByDN", func(t *testing.T) {
		getUserReq := adc.GetUserArgs{
			Dn: "CN=testuser1,CN=Users,DC=adc,DC=dev",
		}
		user, err := cl.GetUser(getUserReq)
		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, getUserReq.Dn, user.DN)
	})
	t.Run("OkByFilter", func(t *testing.T) {
		getUserReq := adc.GetUserArgs{
			Filter: "(&(objectClass=user)(sAMAccountName=testuser1))",
		}
		user, err := cl.GetUser(getUserReq)
		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, "testuser1", user.Id)
	})
	t.Run("OkWithAttributes", func(t *testing.T) {
		req := adc.GetUserArgs{
			Id:         "testuser1",
			Attributes: []string{"sAMAccountName"},
		}
		user, err := cl.GetUser(req)
		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, req.Id, user.Id)
		require.Len(t, user.Attributes, 1)
	})
}
