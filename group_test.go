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

	req := &GetGroupequest{Id: "group2", SkipMembersSearch: true}
	_, err = cl.GetGroup(req)
	require.Error(t, err)

	req = &GetGroupequest{Id: "group1", SkipMembersSearch: true}
	group, err := cl.GetGroup(req)
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
