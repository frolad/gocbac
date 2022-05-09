# gocbac
Simple Golang Content Based Access Control system based to Generics

**Go 1.18 or higher required**

## Description
CBAC allows you to declare list of accesses and function which sets this accesses by content and instance of "USER"

then you can use 3 built-in methods to get access value depending of your needs

## Usage
to prepare CBAC you need 3 types and setter function

### Init
First declare 3 types for Access, Content and User
> you don't have to use custom types, but it makes usage of the package a bit easier if you do
``` go
type Access string

type Content uint64

type User string
```

Then you need to create policies setter function
``` go
func policiesSetter(contentList []Content, user User, requestedAccesses []Access) (gocbac.AccessSetter[Access, Content], error) {
    ...
}
```

Inside policy setter you have to return AccessSetter[Access, Content] function and error
> error will be forwareded to access getters (GetPolicies, GetPolicy, GetAccess) so you can return DB errors and so on if you'd like to handle them, otherwise just return `nil`

Before the `return` you can preprare data (e.g. query DB etc) using contentList, user and requestedAccesses. 
> generally you can skip `requestedAccesses`, but if you would like to optimize data fetching then you can use them to fetch only data which is required to set `requestedAccesses`

Then inside AccessSetter[Access, Content] set access for the content depending on access and data your fetch earlier
``` go
func policiesSetter(contentList []Content, user User, requestedAccesses []Access) (gocbac.AccessSetter[Access, Content], error) {
    // fetch data using contentList, user and requestedAccesses
    myData := ...
    return func(content Content, access Access) bool {
        // determine rather "access" should be true or false using "content", "user" and "access"

        return true | false
    }
}
```

Once policiesSetter is ready you should be ready to init CBAC
``` go
cbac := gocbac.InitCBAC(
	policiesSetter,
	AccessCanView,
	AccessCanEdit,
	AccessCanDelete,
)
```

And use it by calling next this methods:

#### Get the list of policies for the list of content and user
``` go
policies, err := cbac.GetPolicies([]Content{1, 2}, "foo@bar.com")
// or with limited accesses
policies, err := cbac.GetPolicies([]Content{1, 2}, "foo@bar.com", AccessCanView, AccessCanEdit)
```
policies is the map of policies where key is the type of Content: `map[Content]Policy`

#### Get policy for the content and user
``` go
policy, err := cbac.GetPolicy(1, "foo@bar.com")
// or with limited accesses
policy, err := cbac.GetPolicy(1, "foo@bar.com", AccessCanView, AccessCanEdit)
```
policy is the map of booleans where key is the type of Access: `map[Access]bool`

#### Get access for the content and user
``` go
access, err := cbac.GetAccess(1, "foo@bar.com", AccessCanView)
```
access is boolean

### Full example
``` go
import (
	"fmt"

	"github.com/frolad/gocbac"
)

// declare types
type Access string

type Content uint64

type User string

// declare accesses
const (
	AccessCanView   Access = "can_view"
	AccessCanEdit   Access = "can_edit"
	AccessCanDelete Access = "can_delete"
)

// declare setter
func policiesSetter(
	contentList []Content,
	user User,
	requestedAccesses []Access,
) (gocbac.AccessSetter[Access, Content], error) {
	// do content preparation for the list content, users and accesses (e.g. DB queries etc)
	contentPublic := map[Content]bool{
		1: true,
	}

	contentOwners := map[Content]User{
		1: "foo@bar.com",
		2: "bar@foo.com",
	}

	// then fill the access depending on the content
	return func(Content Content, access Access) bool {
		switch access {

		case AccessCanView:
			if _, ok := contentPublic[Content]; ok {
				return true
			} else if owner, ok := contentOwners[Content]; ok {
				return owner == user
			}

			return false

		case AccessCanEdit, AccessCanDelete:
			if owner, ok := contentOwners[Content]; ok {
				return owner == user
			}

			return false
		}

		return false
	}, nil
}

func main() {
	// init cbac
	cbac := gocbac.InitCBAC(
		policiesSetter,
		AccessCanView,
		AccessCanEdit,
		AccessCanDelete,
	)

	// use it
	// by list
	policies, err := cbac.GetPolicies([]Content{1, 2}, "foo@bar.com")
	if err != nil {
		panic(err)
	}
	for content, policy := range policies {
		fmt.Printf("GetPolicies: content: %v, user: foo@bar.com, value: %v\n", content, policy[AccessCanView])
	}

	// by policy
	policy, err := cbac.GetPolicy(1, "foo@bar.com", AccessCanView, AccessCanEdit)
	if err != nil {
		panic(err)
	}
	if policy[AccessCanView] || policy[AccessCanEdit] {
		fmt.Printf("GetPolicy: user had view or edit access\n")
	}

	// by access
	has, err := cbac.GetAccess(1, "foo@bar.com", AccessCanView)
	if err != nil {
		panic(err)
	}

	if has {
		fmt.Printf("GetAccess: has access\n")
	} else {
		fmt.Printf("GetAccess: no access\n")
	}
}
```
