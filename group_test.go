package adc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Group_GetStringAttribute(t *testing.T) {
	u := &Group{
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

func Test_GetGroupRequest_Validate(t *testing.T) {
	var req GetGroupArgs
	err := req.Validate()
	require.Error(t, err)

	req = GetGroupArgs{}
	err1 := req.Validate()
	require.Error(t, err1)

	req = GetGroupArgs{Id: "fake"}
	errOk := req.Validate()
	require.NoError(t, errOk)
}

func Test_Client_GetGroup(t *testing.T) {
	cl := New(&Config{}, withMock())
	err := cl.Connect()
	require.NoError(t, err)

	var badArgs GetGroupArgs
	_, badReqErr := cl.GetGroup(badArgs)
	require.Error(t, badReqErr)

	args := GetGroupArgs{Id: "entryForErr", SkipMembersSearch: true}
	_, err = cl.GetGroup(args)
	require.Error(t, err)

	args = GetGroupArgs{Id: "groupFake", SkipMembersSearch: true}
	group, err := cl.GetGroup(args)
	require.NoError(t, err)
	require.Nil(t, group)

	// Too many entries error
	group, err = cl.GetGroup(GetGroupArgs{
		Id:                "notUniq",
		SkipMembersSearch: true,
		Attributes:        []string{"sAMAccountName"},
	})
	require.Error(t, err)
	require.Nil(t, group)

	// Group with err members get
	group, err = cl.GetGroup(GetGroupArgs{
		Id:                "groupWithErrMember",
		SkipMembersSearch: false,
	})
	require.Error(t, err)
	require.Nil(t, group)

	dnArgs := GetGroupArgs{Dn: "OU=group1,DC=company,DC=com", SkipMembersSearch: true}
	groupByDn, err := cl.GetGroup(dnArgs)
	require.NoError(t, err)
	require.NotNil(t, groupByDn)
	require.Equal(t, dnArgs.Dn, groupByDn.DN)

	args = GetGroupArgs{Id: "group1", SkipMembersSearch: true}
	group, err = cl.GetGroup(args)
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
}

func Test_popAddGroupMembers(t *testing.T) {
	group := &Group{
		Members: []GroupMember{
			{DN: "oldone"},
		},
	}

	notChanged := popAddGroupMembers(group, nil)
	require.Equal(t, group.MembersDn(), notChanged)

	changed := popAddGroupMembers(group, []string{"newone"})
	require.Equal(t, []string{"oldone", "newone"}, changed)
}

func Test_AddGroupMembers(t *testing.T) {
	cl := New(&Config{}, withMock())
	err := cl.Connect()
	require.NoError(t, err)

	_, err = cl.AddGroupMembers("entryForErr", "user1")
	require.Error(t, err)

	_, err = cl.AddGroupMembers("groupFake", "user1")
	require.Error(t, err)

	// Non exists user
	added, err := cl.AddGroupMembers("group1", "userFake")
	require.NoError(t, err)
	require.Equal(t, 0, added)

	// Error user
	_, err = cl.AddGroupMembers("group1", "entryForErr")
	require.Error(t, err)

	// Already member user
	added, err = cl.AddGroupMembers("group1", "user1")
	require.NoError(t, err)
	require.Equal(t, 0, added)

	// Ok user
	added, err = cl.AddGroupMembers("group1", "userToAdd")
	require.NoError(t, err)
	require.Equal(t, 1, added)

}

func Test_popDelGroupMembers(t *testing.T) {
	group := &Group{
		Members: []GroupMember{
			{DN: "oldone"},
			{DN: "oldone2"},
			{DN: "oldone3"},
		},
	}

	notChanged := popDelGroupMembers(group, nil)
	require.Equal(t, group.MembersDn(), notChanged)

	changed := popDelGroupMembers(group, []string{"oldone2"})
	require.Equal(t, []string{"oldone", "oldone3"}, changed)
}

func Test_DeleteGroupMembers(t *testing.T) {
	cl := New(&Config{}, withMock())
	err := cl.Connect()
	require.NoError(t, err)

	_, err = cl.DeleteGroupMembers("entryForErr", "user1")
	require.Error(t, err)

	_, err = cl.DeleteGroupMembers("groupFake", "user1")
	require.Error(t, err)

	// Non exists user
	deleted, err := cl.DeleteGroupMembers("group1", "userFake")
	require.NoError(t, err)
	require.Equal(t, 0, deleted)

	// Error user
	_, err = cl.DeleteGroupMembers("group1", "entryForErr")
	require.Error(t, err)

	// Already not member user
	deleted, err = cl.DeleteGroupMembers("group1", "user2")
	require.NoError(t, err)
	require.Equal(t, 0, deleted)

	// Ok user
	deleted, err = cl.DeleteGroupMembers("group1", "user1")
	require.NoError(t, err)
	require.Equal(t, 1, deleted)

}

func Test_Group_MembersDn(t *testing.T) {
	g := &Group{
		Members: []GroupMember{},
	}
	require.Nil(t, g.MembersDn())

	newMem := GroupMember{Id: "someId", DN: "someDn"}
	g.Members = append(g.Members, newMem)
	require.NotNil(t, g.MembersDn())
	require.Contains(t, g.MembersDn(), newMem.DN)
}

func Test_Group_MembersId(t *testing.T) {
	g := &Group{
		Members: []GroupMember{},
	}
	require.Nil(t, g.MembersId())

	newMem := GroupMember{Id: "someId", DN: "someDn"}
	g.Members = append(g.Members, newMem)
	require.NotNil(t, g.MembersId())
	require.Contains(t, g.MembersId(), newMem.Id)
}
