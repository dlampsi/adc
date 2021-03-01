# adc

[![build-img]][build-url]
[![doc-img]][doc-url]
[![coverage-img]][coverage-url]

Active Directory client library.

The library is a wrapper around  [go-ldap/ldap](github.com/go-ldap/ldap) module that provides a more convient client for Active Directory.

## Usage

Import module in your go app:

```go
import "github.com/dlampsi/adc"
```

### Getting started

```go
// Init client
cl := adc.New(&adc.Config{
    URL:         "ldaps://my.ad.site:636",
    SearchBase:  "OU=some,DC=company,DC=com",
    Bind: &adc.BindAccount{
        DN:       "CN=admin,DC=company,DC=com",
        Password: "***",
    },
})

// Connect
if err := cl.Connect(); err != nil {
    // Handle error
}

// Search for a user
user, err := cl.GetUser(&adc.GetUserRequest{Id:"userId"})
if err != nil {
    // Handle error
}
if user == nil {
    // Handle not found
}
fmt.Println(user)

// Search for a group
group, err := cl.GetGroup(&adc.GetGroupequest{Id:"groupId"})
if err != nil {
    // Handle error
}
if group == nil {
    // Handle not found
}
fmt.Println(group)

// Add new users to group members
added, err := cl.AddGroupMembers("groupId", "newUserId1", "newUserId2", "newUserId3")
if err != nil {
    // Handle error
}
fmt.Printf("Added %d members", added)


// Delete users from group members
deleted, err := cl.DeleteGroupMembers("groupId", "userId1", "userId2")
if err != nil {
    // Handle error
}
fmt.Printf("Deleted %d users from group members", deleted)

```

### Default config file

By default client initializes with default config file. You can find it in [DefaultUsersConfigs()](config.go) func.

### Check auth by creds

Custom check authentification for provided credentials:

```go
if err := cl.CheckAuthByDN("CN=user,DC=company,DC=com", "password"); err != nil {
    // Handle bad credentials error
}
```

### Custom search base

You can set custom search base for user/group in config:

```go
cfg := &adc.Config{
    URL:         "ldaps://my.ad.site:636",
    SearchBase:  "OU=some,DC=company,DC=com",
    Bind: &adc.BindAccount{DN: "CN=admin,DC=company,DC=com", Password: "***"},
    Users: &adc.UsersConfigs{
        SearchBase: "OU=users_base,DC=company,DC=com",,
    },
}

cl := New(cfg)

if err := cl.Connect(); err != nil {
    // Handle error
}
```


### Custom entries attributes

You can parse custom attributes to client config to fetch those attributes during users or groups fetch:
```go
// Append new attributes to existsing user attributes
cl.Config().AppendUsesAttributes("manager")

// Search for a user
user, err := cl.GetUser(&adc.GetUserRequest{Id:"userId"})
if err != nil {
    // Handle error
}
if user == nil {
    // Handle not found
}

// Get custom attribute
userManager := exists.GetStringAttribute("manager")
fmt.Println(userManager)
```

Also you can parse custom attributes during each get requests:
```go
user, err := cl.GetUser(&adc.GetUserRequest{Id: "userId", Attributes: []string{"manager"}})
if err != nil {
    // Handle error
}
// Get custom attribute
userManager := exists.GetStringAttribute("manager")
fmt.Println(userManager)
```


### Custom search filters

You can parse custom search filters to client config:

```go
cfg := &adc.Config{
    URL:         "ldaps://my.ad.site:636",
    SearchBase:  "OU=some,DC=company,DC=com",
    Bind: &adc.BindAccount{DN: "CN=admin,DC=company,DC=com", Password: "***"},
    Users: &adc.UsersConfigs{
        FilterById: "(&(objectClass=person)(cn=%v))",
    },
    Groups: &adc.GroupsConfigs{
        FilterById: "(&(objectClass=group)(cn=%v))",
    },
}
cl := New(cfg)
if err := cl.Connect(); err != nil {
    // Handle error
}
```

## Contributing

1. Create new PR from `main` branch
2. Request review from maintainers

## License

[MIT License](LICENSE).


[build-img]: https://github.com/dlampsi/adc/workflows/build/badge.svg
[build-url]: https://github.com/dlampsi/adc/actions
[coverage-img]: https://codecov.io/gh/dlampsi/adc/branch/main/graph/badge.svg
[coverage-url]: https://codecov.io/gh/dlampsi/adc
[doc-img]: https://pkg.go.dev/badge/dlampsi/adc
[doc-url]: https://pkg.go.dev/github.com/dlampsi/adc