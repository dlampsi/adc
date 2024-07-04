package adctests

import (
	"testing"
	"time"

	"github.com/dlampsi/adc"
	"github.com/stretchr/testify/require"
)

func Test_Group_GetStringAttribute(t *testing.T) {
	t.Run("NonExists", func(t *testing.T) {
		g := &adc.Group{
			Attributes: map[string]interface{}{
				"one": "string",
			},
		}
		require.Empty(t, g.GetStringAttribute("nonexists"))
	})
	t.Run("NonString", func(t *testing.T) {
		g := &adc.Group{
			Attributes: map[string]interface{}{
				"two":   2,
				"three": []byte("bytedata"),
			},
		}
		require.Equal(t, "", g.GetStringAttribute("two"))
		require.Equal(t, "", g.GetStringAttribute("three"))
	})
	t.Run("Ok", func(t *testing.T) {
		g := &adc.Group{
			Attributes: map[string]interface{}{
				"one": "value",
			},
		}
		require.Equal(t, "value", g.GetStringAttribute("one"))
	})
}

func Test_GetGroupRequest_Validate(t *testing.T) {
	t.Run("ErrWithNil", func(t *testing.T) {
		var req adc.GetGroupArgs
		require.Error(t, req.Validate())
	})
	t.Run("ErrWithEmpty", func(t *testing.T) {
		req := adc.GetGroupArgs{}
		require.Error(t, req.Validate())
	})
	t.Run("OkWithId", func(t *testing.T) {
		req := adc.GetGroupArgs{Id: "fake"}
		require.NoError(t, req.Validate())
	})
	t.Run("OkWithFilter", func(t *testing.T) {
		req := adc.GetGroupArgs{Filter: "fake"}
		require.NoError(t, req.Validate())
	})
	t.Run("OkWithDn", func(t *testing.T) {
		req := adc.GetGroupArgs{Dn: "fake"}
		require.NoError(t, req.Validate())
	})
}

func Test_Group_MembersDn(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		g := &adc.Group{Members: []adc.GroupMember{}}
		require.Nil(t, g.MembersDn())
	})
	t.Run("Ok", func(t *testing.T) {
		g := &adc.Group{Members: []adc.GroupMember{{Id: "someId", DN: "someDn"}}}
		require.NotNil(t, g.MembersDn())
		require.Contains(t, g.MembersDn(), g.Members[0].DN)
	})
}

func Test_Group_MembersId(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		g := &adc.Group{Members: []adc.GroupMember{}}
		require.Nil(t, g.MembersId())
	})
	t.Run("Ok", func(t *testing.T) {
		g := &adc.Group{Members: []adc.GroupMember{{Id: "someId", DN: "someDn"}}}
		require.NotNil(t, g.MembersId())
		require.Contains(t, g.MembersId(), g.Members[0].Id)
	})
}

func Test_Client_GetGroup(t *testing.T) {
	cfg := getClientConfig()
	cl := adc.New(&cfg, adc.WithLogger(&logger{t: t}))
	require.NoError(t, cl.Connect())

	t.Run("ErrWithNil", func(t *testing.T) {
		var req adc.GetGroupArgs
		group, err := cl.GetGroup(req)
		require.Error(t, err)
		require.Nil(t, group, "Group should be nil on error")
	})

	t.Run("Non exists", func(t *testing.T) {
		req := adc.GetGroupArgs{
			Id: "nonexists",
		}
		group, err := cl.GetGroup(req)
		require.NoError(t, err, "Expected no error on non exists group")
		require.Nil(t, group, "Non exists group error should return nil")
	})
	t.Run("TooManyEntries", func(t *testing.T) {
		req := adc.GetGroupArgs{
			Filter: "(&(objectClass=group))",
		}
		group, err := cl.GetGroup(req)
		require.Error(t, err, "Expected error on too many entries")
		require.Nil(t, group, "Group should be nil on error")
	})
	t.Run("OkById", func(t *testing.T) {
		req := adc.GetGroupArgs{
			Id: "testgroup2",
		}
		group, err := cl.GetGroup(req)
		require.NoError(t, err)
		require.NotNil(t, group)
		require.Equal(t, req.Id, group.Id)
		require.NotEmpty(t, group.Members)
	})
	t.Run("OkBySkipMembersSearch", func(t *testing.T) {
		req := adc.GetGroupArgs{
			Id:                "testgroup2",
			SkipMembersSearch: true,
		}
		group, err := cl.GetGroup(req)
		require.NoError(t, err)
		require.NotNil(t, group)
		require.Equal(t, req.Id, group.Id)
		require.Empty(t, group.Members)
	})
	t.Run("OkByDn", func(t *testing.T) {
		req := adc.GetGroupArgs{
			Dn: "CN=testgroup2,CN=Users,DC=adc,DC=dev",
		}
		group, err := cl.GetGroup(req)
		require.NoError(t, err)
		require.NotNil(t, group)
		require.Equal(t, req.Dn, group.DN)
	})
	t.Run("OkByFilter", func(t *testing.T) {
		req := adc.GetGroupArgs{
			Filter: "(&(objectClass=group)(cn=testgroup2))",
		}
		group, err := cl.GetGroup(req)
		require.NoError(t, err)
		require.NotNil(t, group)
		require.Equal(t, "testgroup2", group.Id)
	})
	t.Run("OkWithAttributes", func(t *testing.T) {
		req := adc.GetGroupArgs{
			Id:         "testgroup2",
			Attributes: []string{"sAMAccountName"},
		}
		group, err := cl.GetGroup(req)
		require.NoError(t, err)
		require.NotNil(t, group)
		require.Equal(t, req.Id, group.Id)
		require.NotEmpty(t, group.Attributes)
		require.Len(t, group.Attributes, 1)
	})
}

func Test_Client_AddGroupMembers(t *testing.T) {
	cfg := getClientConfig()
	cl := adc.New(&cfg, adc.WithLogger(&logger{t: t}))
	require.NoError(t, cl.Connect())

	t.Run("ErrBadGroupId", func(t *testing.T) {
		cnt, err := cl.AddGroupMembers("", "someDn")
		require.Error(t, err, "Expected error on non exists group")
		require.Zero(t, cnt, "Added members count should be zero on error")
	})
	t.Run("ErrNonExistsGroup", func(t *testing.T) {
		cnt, err := cl.AddGroupMembers("nonexists", "someDn")
		require.Error(t, err, "Expected error on bad group id")
		require.Zero(t, cnt, "Added members count should be zero on error")
	})
	t.Run("AlreadyAMember", func(t *testing.T) {
		cnt, err := cl.AddGroupMembers("testgroup1", "testuser1")
		require.NoError(t, err, "Expected no error on already a member")
		require.Zero(t, cnt, "Added members count should be zero on already a member")
	})
	t.Run("BadMember", func(t *testing.T) {
		cnt, err := cl.AddGroupMembers("testgroup1", "")
		require.Error(t, err, "Expected error on bad member id")
		require.Equal(t, 0, cnt, "Added members count should be zero on error")
	})
	t.Run("NonExistsMember", func(t *testing.T) {
		cnt, err := cl.AddGroupMembers("testgroup1", "nonexists")
		require.NoError(t, err, "No error on non exists member")
		require.Equal(t, 0, cnt, "Added members count should be zero on error")
	})
	t.Run("Ok", func(t *testing.T) {
		const (
			groupId  = "testgroup1"
			memberId = "testuser2"
		)

		cnt, err := cl.AddGroupMembers(groupId, memberId)
		require.NoError(t, err)
		require.Equal(t, 1, cnt)

		req := adc.GetGroupArgs{
			Id: groupId,
		}
		group, err := cl.GetGroup(req)
		require.NoError(t, err)
		require.NotNil(t, group)
		require.Len(t, group.Members, 2)

		dc, err := cl.DeleteGroupMembers(groupId, memberId)
		require.NoError(t, err)
		require.Equal(t, 1, dc)
	})
}

func Test_Client_DeleteGroupMembers(t *testing.T) {
	cfg := getClientConfig()
	cl := adc.New(&cfg, adc.WithLogger(&logger{t: t}))
	require.NoError(t, cl.Connect())

	t.Run("ErrBadGroupId", func(t *testing.T) {
		cnt, err := cl.DeleteGroupMembers("", "someDn")
		require.Error(t, err, "Expected error on non exists group")
		require.Zero(t, cnt, "Deleted members count should be zero on error")
	})
	t.Run("ErrNonExistsGroup", func(t *testing.T) {
		cnt, err := cl.DeleteGroupMembers("nonexists", "someDn")
		require.Error(t, err, "Expected error on bad group id")
		require.Zero(t, cnt, "Added groups count should be zero on error")
	})
	t.Run("BadMember", func(t *testing.T) {
		cnt, err := cl.DeleteGroupMembers("testgroup1", "")
		require.Error(t, err, "Expected error on bad member id")
		require.Equal(t, 0, cnt, "Deleted members count should be zero on error")
	})
	t.Run("NonExistsMember", func(t *testing.T) {
		cnt, err := cl.DeleteGroupMembers("testgroup1", "nonexists")
		require.NoError(t, err, "No error on non exists member")
		require.Equal(t, 0, cnt, "Deleted members count should be zero on error")
	})
	t.Run("AlreadyNotAMember", func(t *testing.T) {
		cnt, err := cl.DeleteGroupMembers("testgroup1", "testuser2")
		require.NoError(t, err, "Expected no error on already a member")
		require.Zero(t, cnt, "Added groups count should be zero on already a member")
	})
	t.Run("Ok", func(t *testing.T) {
		const (
			groupId  = "testgroup2"
			memberId = "testuser2"
		)

		cnt, err := cl.DeleteGroupMembers(groupId, memberId)
		require.NoError(t, err)
		require.Equal(t, 1, cnt)

		req := adc.GetGroupArgs{
			Id: groupId,
		}
		group, err := cl.GetGroup(req)
		require.NoError(t, err)
		require.NotNil(t, group)
		require.Len(t, group.Members, 0)

		dc, err := cl.AddGroupMembers(groupId, memberId)
		require.NoError(t, err)
		require.Equal(t, 1, dc)
	})
}

func Test_Client_CreateGroup(t *testing.T) {
	cfg := getClientConfig()
	cl := adc.New(&cfg, adc.WithLogger(&logger{t: t}))
	require.NoError(t, cl.Connect())

	t.Run("BadArgs", func(t *testing.T) {
		req := adc.CreateGroupArgs{}
		require.Error(t, cl.CreateGroup(req))
	})

	t.Run("Ok", func(t *testing.T) {
		req := adc.CreateGroupArgs{
			Id: "createdGroup" + time.Now().Format("20060102150405"),
			Attributes: map[string][]string{
				"description": {"Test group 3"},
			},
		}
		require.NoError(t, cl.CreateGroup(req))

		group, err := cl.GetGroup(adc.GetGroupArgs{Id: req.Id})
		require.NoError(t, err)
		require.NotNil(t, group, "Created group should be found")
		require.Equal(t, req.Id, group.Id)
		require.Equal(t, req.Attributes["description"][0], group.GetStringAttribute("description"))
		require.Len(t, group.Members, 0, "Group should have no members")
	})
}

func Test_Client_DeleteGroup(t *testing.T) {
	cfg := getClientConfig()
	cl := adc.New(&cfg, adc.WithLogger(&logger{t: t}))
	require.NoError(t, cl.Connect())

	t.Run("BadGroupId", func(t *testing.T) {
		require.Error(t, cl.DeleteGroup(""))
	})
	t.Run("NonExistsGroup", func(t *testing.T) {
		require.NoError(t, cl.DeleteGroup("nonexists"), "No error on non exists (maybe already deleted) group")
	})

	t.Run("Ok", func(t *testing.T) {
		req := adc.CreateGroupArgs{
			Id: "groupForDelete" + time.Now().Format("20060102150405"),
		}
		require.NoError(t, cl.CreateGroup(req))

		created, err := cl.GetGroup(adc.GetGroupArgs{Id: req.Id})
		require.NoError(t, err)
		require.NotNil(t, created)

		require.NoError(t, cl.DeleteGroup(req.Id))

		deleted, err := cl.GetGroup(adc.GetGroupArgs{Id: req.Id})
		require.NoError(t, err)
		require.Nil(t, deleted, "Deleted group should not be found")
	})
}
