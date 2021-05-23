package adc

import (
	"errors"
	"fmt"
	"sync"

	"github.com/dlampsi/generigo"
	"github.com/go-ldap/ldap/v3"
)

// Active Direcotry group.
type Group struct {
	DN         string                 `json:"dn"`
	Id         string                 `json:"id"`
	Attributes map[string]interface{} `json:"attributes"`
	Members    []GroupMember          `json:"members"`
}

// Active Direcotry member info.
type GroupMember struct {
	DN string `json:"dn"`
	Id string `json:"id"`
}

// Returns string attribute by attribute name.
// Returns empty string if attribute not exists or it can't be covnerted to string.
func (g *Group) GetStringAttribute(name string) string {
	for att, val := range g.Attributes {
		if att == name {
			if s, ok := val.(string); ok {
				return s
			}
		}
	}
	return ""
}

type GetGroupArgs struct {
	// Group ID to search.
	Id string `json:"id"`
	// Optional group DN. Overwrites ID if provided in request.
	Dn string `json:"dn"`
	// Optional group attributes to overwrite attributes in client config.
	Attributes []string `json:"attributes"`
	// Skip search of group members data. Can improve request time.
	SkipMembersSearch bool `json:"skip_members_search"`
}

func (args GetGroupArgs) Validate() error {
	if args.Id == "" && args.Dn == "" {
		return errors.New("neither of ID of DN provided")
	}
	return nil
}

func (cl *Client) GetGroup(args GetGroupArgs) (*Group, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}

	filter := fmt.Sprintf(cl.cfg.Groups.FilterById, args.Id)
	if args.Dn != "" {
		filter = fmt.Sprintf(cl.cfg.Groups.FilterByDn, ldap.EscapeFilter(args.Dn))
	}

	req := &ldap.SearchRequest{
		BaseDN:       cl.cfg.Groups.SearchBase,
		Scope:        ldap.ScopeWholeSubtree,
		DerefAliases: ldap.NeverDerefAliases,
		TimeLimit:    int(cl.cfg.Timeout.Seconds()),
		Filter:       filter,
		Attributes:   cl.cfg.Groups.Attributes,
	}
	if args.Attributes != nil {
		req.Attributes = args.Attributes
	}

	entry, err := cl.searchEntry(req)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	result := &Group{
		DN:         entry.DN,
		Id:         entry.GetAttributeValue(cl.cfg.Groups.IdAttribute),
		Attributes: make(map[string]interface{}, len(entry.Attributes)),
	}
	for _, a := range entry.Attributes {
		result.Attributes[a.Name] = entry.GetAttributeValue(a.Name)
	}

	if !args.SkipMembersSearch {
		members, err := cl.getGroupMembers(entry.DN)
		if err != nil {
			return nil, fmt.Errorf("can't get group members: %s", err.Error())
		}
		result.Members = members
	}

	return result, nil
}

func (cl *Client) getGroupMembers(dn string) ([]GroupMember, error) {
	req := &ldap.SearchRequest{
		BaseDN:       cl.cfg.Users.SearchBase,
		Scope:        ldap.ScopeWholeSubtree,
		DerefAliases: ldap.NeverDerefAliases,
		TimeLimit:    int(cl.cfg.Timeout.Seconds()),
		Filter:       fmt.Sprintf(cl.cfg.Groups.FilterMembersByDn, ldap.EscapeFilter(dn)),
		Attributes:   []string{cl.cfg.Users.IdAttribute},
	}
	entries, err := cl.searchEntries(req)
	if err != nil {
		return nil, err
	}
	var result []GroupMember
	for _, e := range entries {
		result = append(result, GroupMember{
			DN: e.DN,
			Id: e.GetAttributeValue(cl.cfg.Groups.IdAttribute),
		})
	}
	return result, nil
}

// Returns list of group members DNs.
func (g *Group) MembersDn() []string {
	var result []string
	for _, m := range g.Members {
		result = append(result, m.DN)
	}
	return result
}

// Returns list of group members IDs.
func (g *Group) MembersId() []string {
	var result []string
	for _, m := range g.Members {
		result = append(result, m.Id)
	}
	return result
}

// Adds provided accounts IDs to provided group members. Returns number of addedd accounts.
func (cl *Client) AddGroupMembers(groupId string, membersIds ...string) (int, error) {
	group, err := cl.GetGroup(GetGroupArgs{Id: groupId})
	if err != nil {
		return 0, fmt.Errorf("can't get group: %s", err.Error())
	}
	if group == nil {
		return 0, fmt.Errorf("group '%s' not found by ID", groupId)
	}

	ch := make(chan string, len(membersIds))
	errCh := make(chan error, len(membersIds))
	wg := &sync.WaitGroup{}

	for _, id := range membersIds {
		wg.Add(1)
		go func(userId string, ch chan<- string, errCh chan<- error, wg *sync.WaitGroup) {
			defer wg.Done()
			user, err := cl.GetUser(GetUserArgs{Id: userId})
			if err != nil {
				errCh <- fmt.Errorf("can't get account '%s': %s", userId, err.Error())
				return
			}
			if user == nil {
				cl.logger.Debugf("Account '%s' being added to '%s' wasn't found",
					userId, groupId)
				return
			}
			if user.IsGroupMember(groupId) {
				cl.logger.Debugf("The adding account '%s' is already a member of the group '%s'",
					userId, groupId)
				return
			}
			ch <- user.DN
		}(id, ch, errCh, wg)
	}
	wg.Wait()
	close(errCh)
	close(ch)

	for err := range errCh {
		if err != nil {
			return 0, err
		}
	}

	var toAdd []string
	for dn := range ch {
		toAdd = append(toAdd, dn)
	}
	if len(toAdd) == 0 {
		return 0, nil
	}

	newMembers := popAddGroupMembers(group, toAdd)

	cl.logger.Debugf("Adding new group members to '%s'; Old count: %d; New count: %d",
		groupId, len(group.MembersId()), len(newMembers))

	if err := cl.updateAttribute(group.DN, "member", newMembers); err != nil {
		return 0, err
	}

	return len(toAdd), nil
}

func popAddGroupMembers(g *Group, toAdd []string) []string {
	if len(toAdd) == 0 {
		return g.MembersDn()
	}
	result := make([]string, 0, len(g.Members)+len(toAdd))
	result = append(result, g.MembersDn()...)
	result = append(result, toAdd...)
	return result
}

// Deletes provided accounts IDs from provided group members. Returns number of deleted from group members.
func (cl *Client) DeleteGroupMembers(groupId string, membersIds ...string) (int, error) {
	group, err := cl.GetGroup(GetGroupArgs{Id: groupId})
	if err != nil {
		return 0, fmt.Errorf("can't get group: %s", err.Error())
	}
	if group == nil {
		return 0, fmt.Errorf("group '%s' not found by ID", groupId)
	}

	ch := make(chan string, len(membersIds))
	errCh := make(chan error, len(membersIds))
	wg := &sync.WaitGroup{}

	for _, id := range membersIds {
		wg.Add(1)
		go func(userId string, ch chan<- string, errCh chan<- error, wg *sync.WaitGroup) {
			defer wg.Done()
			user, err := cl.GetUser(GetUserArgs{Id: userId})
			if err != nil {
				errCh <- fmt.Errorf("can't get account '%s': %s", userId, err.Error())
				return
			}
			if user == nil {
				cl.logger.Debugf("Account '%s' being deleted from '%s' wasn't found",
					userId, groupId)
				return
			}
			if !user.IsGroupMember(groupId) {
				cl.logger.Debugf("The deleting account '%s' already isn't a member of the group '%s'",
					userId, groupId)
				return
			}
			ch <- user.DN
		}(id, ch, errCh, wg)
	}
	wg.Wait()
	close(errCh)
	close(ch)

	for err := range errCh {
		if err != nil {
			return 0, err
		}
	}

	var toDel []string
	for dn := range ch {
		toDel = append(toDel, dn)
	}
	if len(toDel) == 0 {
		return 0, nil
	}

	newMembers := popDelGroupMembers(group, toDel)

	cl.logger.Debugf("Deleting members from group '%s'; Old count: %d; New count: %d",
		groupId, len(group.MembersId()), len(newMembers))

	if err := cl.updateAttribute(group.DN, "member", newMembers); err != nil {
		return 0, err
	}

	return len(toDel), nil
}

func popDelGroupMembers(g *Group, toDel []string) []string {
	result := []string{}
	for _, memberDN := range g.MembersDn() {
		if !generigo.StringInSlice(memberDN, toDel) {
			result = append(result, memberDN)
		}
	}
	return result
}
