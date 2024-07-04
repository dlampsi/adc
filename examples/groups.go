package examples

import (
	"fmt"

	"github.com/dlampsi/adc"
)

func mainGroups() {
	cfg := &adc.Config{
		URL: "ldaps://my.ad.site:636",
		Bind: &adc.BindAccount{
			DN:       "CN=admin,DC=company,DC=com",
			Password: "***",
		},
		SearchBase: "OU=default,DC=company,DC=com",
	}

	cl := adc.New(cfg)

	if err := cl.Connect(); err != nil {
		panic(err)
	}

	/* -------------- Create -------------- */

	createReq := adc.CreateGroupArgs{
		Id: "exampleGroupId",
		Attributes: map[string][]string{
			"description": {"Example group"},
		},
	}
	if err := cl.CreateGroup(createReq); err != nil {
		panic(err)
	}

	/* -------------- Search -------------- */

	getReq := adc.GetGroupArgs{
		Id: "exampleGroupId",
	}
	group, err := cl.GetGroup(getReq)
	if err != nil {
		panic(err)
	}
	fmt.Println(group)

	/* -------------- Add group members -------------- */

	added, err := cl.AddGroupMembers("exampleGroupId", "newUserId1", "newUserId2", "newUserId3")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Added %d members", added)

	/* -------------- Delete group members -------------- */

	deleted, err := cl.DeleteGroupMembers("exampleGroupId", "userId1", "userId2")
	if err != nil {
		// Handle error
	}
	fmt.Printf("Deleted %d users from group members", deleted)

	/* -------------- Delete members -------------- */

	if err := cl.DeleteGroup(group.DN); err != nil {
		panic(err)
	}
}
