package adc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Group_GetStringAttribute(t *testing.T) {
	t.Run("NonExists", func(t *testing.T) {
		g := &Group{
			Attributes: map[string]interface{}{
				"one": "string",
			},
		}
		require.Empty(t, g.GetStringAttribute("nonexists"))
	})
	t.Run("NonString", func(t *testing.T) {
		g := &Group{
			Attributes: map[string]interface{}{
				"two":   2,
				"three": []byte("bytedata"),
			},
		}
		require.Equal(t, "", g.GetStringAttribute("two"))
		require.Equal(t, "", g.GetStringAttribute("three"))
	})
	t.Run("Ok", func(t *testing.T) {
		g := &Group{
			Attributes: map[string]interface{}{
				"one": "value",
			},
		}
		require.Equal(t, "value", g.GetStringAttribute("one"))
	})
}

func Test_GetGroupRequest_Validate(t *testing.T) {
	t.Run("ErrWithNil", func(t *testing.T) {
		var req GetGroupArgs
		require.Error(t, req.Validate())
	})
	t.Run("ErrWithEmpty", func(t *testing.T) {
		req := GetGroupArgs{}
		require.Error(t, req.Validate())
	})
	t.Run("OkWithId", func(t *testing.T) {
		req := GetGroupArgs{Id: "fake"}
		require.NoError(t, req.Validate())
	})
	t.Run("OkWithFilter", func(t *testing.T) {
		req := GetGroupArgs{Filter: "fake"}
		require.NoError(t, req.Validate())
	})
	t.Run("OkWithDn", func(t *testing.T) {
		req := GetGroupArgs{Dn: "fake"}
		require.NoError(t, req.Validate())
	})
}

func Test_Client_GetGroup(t *testing.T) {
	cl := newMockClient(&Config{})
	require.NoError(t, cl.Connect())

	t.Run("BadArgs", func(t *testing.T) {
		var badArgs GetGroupArgs
		g, badReqErr := cl.GetGroup(badArgs)
		require.Error(t, badReqErr)
		require.Nil(t, g, "Group should be nil on error")
	})
	t.Run("Err", func(t *testing.T) {
		args := GetGroupArgs{Id: "entryForErr", SkipMembersSearch: true}
		g, err := cl.GetGroup(args)
		require.Error(t, err)
		require.Nil(t, g, "Group should be nil on error")
	})
	t.Run("NonExistsGroup", func(t *testing.T) {
		args := GetGroupArgs{Id: "groupFake", SkipMembersSearch: true}
		group, err := cl.GetGroup(args)
		require.NoError(t, err)
		require.Nil(t, group)
	})
	t.Run("TooManyEntries", func(t *testing.T) {
		// Too many entries error
		group, err := cl.GetGroup(GetGroupArgs{
			Id:                "notUniq",
			SkipMembersSearch: true,
			Attributes:        []string{"sAMAccountName"},
		})
		require.Error(t, err)
		require.Nil(t, group)
	})
	t.Run("GroupWithErrMember", func(t *testing.T) {
		group, err := cl.GetGroup(GetGroupArgs{
			Id:                "groupWithErrMember",
			SkipMembersSearch: false,
		})
		require.Error(t, err)
		require.Nil(t, group)
	})
	t.Run("ByDn", func(t *testing.T) {
		dnArgs := GetGroupArgs{Dn: "OU=group1,DC=company,DC=com", SkipMembersSearch: true}
		groupByDn, err := cl.GetGroup(dnArgs)
		require.NoError(t, err)
		require.NotNil(t, groupByDn)
		require.Equal(t, dnArgs.Dn, groupByDn.DN)
	})
	t.Run("ByFilter", func(t *testing.T) {
		filterReq := GetGroupArgs{Filter: "customFilterToSearchGroup", SkipMembersSearch: true}
		userByFilter, err := cl.GetGroup(filterReq)
		require.NoError(t, err)
		require.NotNil(t, userByFilter)
		require.Equal(t, userByFilter.Id, "group1")
	})
	t.Run("Ok", func(t *testing.T) {
		args := GetGroupArgs{Id: "group1", SkipMembersSearch: true}
		group, err := cl.GetGroup(args)
		require.NoError(t, err)
		require.NotNil(t, group)
		require.Equal(t, args.Id, group.Id)

		args.Attributes = []string{"something"}
		group, err = cl.GetGroup(args)
		require.NoError(t, err)
		require.NotNil(t, group)
		require.Equal(t, args.Id, group.Id)

		args.SkipMembersSearch = false
		group, err = cl.GetGroup(args)
		require.NoError(t, err)
		require.NotNil(t, group)
		require.Equal(t, args.Id, group.Id)
		require.NotNil(t, group.Members)
		require.Len(t, group.Members, 1)
	})
}

func Test_AddGroupMembers(t *testing.T) {
	cl := newMockClient(&Config{})
	require.NoError(t, cl.Connect())

	t.Run("BadArgs", func(t *testing.T) {
		num, err := cl.AddGroupMembers("entryForErr", "user1")
		require.Error(t, err)
		require.Equal(t, 0, num)
	})
	t.Run("Err", func(t *testing.T) {
		num, err := cl.AddGroupMembers("groupFake", "user1")
		require.Error(t, err)
		require.Equal(t, 0, num)
	})
	t.Run("NonExistsUser", func(t *testing.T) {
		added, err := cl.AddGroupMembers("group1", "userFake")
		require.NoError(t, err)
		require.Equal(t, 0, added)
	})
	t.Run("ErrorUser", func(t *testing.T) {
		num, err := cl.AddGroupMembers("group1", "entryForErr")
		require.Error(t, err)
		require.Equal(t, 0, num)
	})
	t.Run("AlreadyMember", func(t *testing.T) {
		added, err := cl.AddGroupMembers("group1", "user1")
		require.NoError(t, err)
		require.Equal(t, 0, added)
	})
	t.Run("Ok", func(t *testing.T) {
		added, err := cl.AddGroupMembers("group1", "userToAdd")
		require.NoError(t, err)
		require.Equal(t, 1, added)
	})
}

func Test_popAddGroupMembers(t *testing.T) {
	group := &Group{
		Members: []GroupMember{{DN: "oldone"}},
	}
	t.Run("Nil", func(t *testing.T) {
		notChanged := popAddGroupMembers(group, nil)
		require.Equal(t, group.MembersDn(), notChanged)
	})
	t.Run("Ok", func(t *testing.T) {
		changed := popAddGroupMembers(group, []string{"newone"})
		require.Equal(t, []string{"oldone", "newone"}, changed)
	})
}

func Test_DeleteGroupMembers(t *testing.T) {
	cl := newMockClient(&Config{})
	require.NoError(t, cl.Connect())

	t.Run("BadArgs", func(t *testing.T) {
		num, err := cl.DeleteGroupMembers("entryForErr", "user1")
		require.Error(t, err)
		require.Equal(t, 0, num)
	})
	t.Run("Err", func(t *testing.T) {
		num, err := cl.DeleteGroupMembers("groupFake", "user1")
		require.Error(t, err)
		require.Equal(t, 0, num)
	})
	t.Run("NonExistsUser", func(t *testing.T) {
		deleted, err := cl.DeleteGroupMembers("group1", "userFake")
		require.NoError(t, err)
		require.Equal(t, 0, deleted)
	})
	t.Run("ErrorUser", func(t *testing.T) {
		num, err := cl.DeleteGroupMembers("group1", "entryForErr")
		require.Error(t, err)
		require.Equal(t, 0, num)
	})
	t.Run("AlreadyNotAMember", func(t *testing.T) {
		deleted, err := cl.DeleteGroupMembers("group1", "user2")
		require.NoError(t, err)
		require.Equal(t, 0, deleted)
	})
	t.Run("Ok", func(t *testing.T) {
		deleted, err := cl.DeleteGroupMembers("group1", "user1")
		require.NoError(t, err)
		require.Equal(t, 1, deleted)
	})
}

func Test_popDelGroupMembers(t *testing.T) {
	group := &Group{
		Members: []GroupMember{
			{DN: "oldone"},
			{DN: "oldone2"},
			{DN: "oldone3"},
		},
	}
	t.Run("Nil", func(t *testing.T) {
		notChanged := popDelGroupMembers(group, nil)
		require.Equal(t, group.MembersDn(), notChanged)
	})
	t.Run("Ok", func(t *testing.T) {
		changed := popDelGroupMembers(group, []string{"oldone2"})
		require.Equal(t, []string{"oldone", "oldone3"}, changed)
	})
}

func Test_Group_MembersDn(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		g := &Group{Members: []GroupMember{}}
		require.Nil(t, g.MembersDn())
	})
	t.Run("Ok", func(t *testing.T) {
		g := &Group{Members: []GroupMember{{Id: "someId", DN: "someDn"}}}
		require.NotNil(t, g.MembersDn())
		require.Contains(t, g.MembersDn(), g.Members[0].DN)
	})
}

func Test_Group_MembersId(t *testing.T) {
	t.Run("Nil", func(t *testing.T) {
		g := &Group{Members: []GroupMember{}}
		require.Nil(t, g.MembersId())
	})
	t.Run("Ok", func(t *testing.T) {
		g := &Group{Members: []GroupMember{{Id: "someId", DN: "someDn"}}}
		require.NotNil(t, g.MembersId())
		require.Contains(t, g.MembersId(), g.Members[0].Id)
	})
}
