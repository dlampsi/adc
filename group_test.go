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
	var req *GetGroupequest
	err := req.Validate()
	require.Error(t, err)

	req = &GetGroupequest{}
	err1 := req.Validate()
	require.Error(t, err1)

	req = &GetGroupequest{Id: "fake"}
	errOk := req.Validate()
	require.NoError(t, errOk)
}

func Test_Client_GetGroup(t *testing.T) {
	cl := New(&Config{}, WithLdapClient(&mockClient{}))

	err := cl.Connect()
	require.NoError(t, err)

	var badReq *GetGroupequest
	_, badReqErr := cl.GetGroup(badReq)
	require.Error(t, badReqErr)

	req := &GetGroupequest{Id: "entryForErr", SkipMembersSearch: true}
	_, err = cl.GetGroup(req)
	require.Error(t, err)

	req = &GetGroupequest{Id: "groupFake", SkipMembersSearch: true}
	group, err := cl.GetGroup(req)
	require.NoError(t, err)
	require.Nil(t, group)

	dnReq := &GetGroupequest{Dn: "OU=group1,DC=company,DC=com", SkipMembersSearch: true}
	groupByDn, err := cl.GetGroup(dnReq)
	require.NoError(t, err)
	require.NotNil(t, groupByDn)
	require.Equal(t, dnReq.Dn, groupByDn.DN)

	req = &GetGroupequest{Id: "group1", SkipMembersSearch: true}
	group, err = cl.GetGroup(req)
	require.NoError(t, err)
	require.NotNil(t, group)
	require.Equal(t, req.Id, group.Id)

	req.Attributes = []string{"something"}
	group, err = cl.GetGroup(req)
	require.NoError(t, err)
	require.NotNil(t, group)
	require.Equal(t, req.Id, group.Id)

	req.SkipMembersSearch = false
	group, err = cl.GetGroup(req)
	require.NoError(t, err)
	require.NotNil(t, group)
	require.Equal(t, req.Id, group.Id)
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
	cl := New(&Config{}, WithLdapClient(&mockClient{}))

	_, err := cl.AddGroupMembers("group2", "user1")
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
	cl := New(&Config{}, WithLdapClient(&mockClient{}))

	_, err := cl.DeleteGroupMembers("group2", "user1")
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
	deleted, err = cl.DeleteGroupMembers("group1", "user3")
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
