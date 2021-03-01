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

func Test_GetUserRequest_Validate(t *testing.T) {
	var req *GetUserRequest
	err := req.Validate()
	require.Error(t, err)

	req = &GetUserRequest{}
	err1 := req.Validate()
	require.Error(t, err1)

	req = &GetUserRequest{Id: "fake"}
	errOk := req.Validate()
	require.NoError(t, errOk)
}

func Test_Client_GetUser(t *testing.T) {
	cl := New(&Config{}, WithLdapClient(&mockClient{}))

	err := cl.Connect()
	require.NoError(t, err)

	var badReq *GetUserRequest
	_, badReqErr := cl.GetUser(badReq)
	require.Error(t, badReqErr)

	req := &GetUserRequest{Id: "user2", SkipGroupsSearch: true}
	_, err = cl.GetUser(req)
	require.Error(t, err)

	req = &GetUserRequest{Id: "user1", SkipGroupsSearch: true}
	user, err := cl.GetUser(req)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, req.Id, user.Id)
	require.Nil(t, user.Groups)

	req.Attributes = []string{"something"}
	user, err = cl.GetUser(req)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, req.Id, user.Id)
	require.Nil(t, user.Groups)

	req.SkipGroupsSearch = false
	user, err = cl.GetUser(req)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, req.Id, user.Id)
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
