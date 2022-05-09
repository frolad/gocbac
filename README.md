# gocbac
Simple Golang Content Based Access Control system

## Usage
``` go
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
	ContentList []Content,
	User User,
	RequestedAccesses []Access,
) (AccessSetter[Access, Content], error) {
    // do content preparation for the list content, users and accesses (e.g. DB queries etc)
    contentPublic := map[Content]bool{
        1: true,
    }

    contentOwners := map[Content]string{
        1: "foo@bar.com",
        2: "bar@foo.com",
    }

    // then fill the access depending on the content
	return func(Content Content, access Access) bool {
        switch access {
		
        case AccessCanView:
			if _, ok := contentPublic[Content]; ok {
                return true
            } else if user, ok := contentOwners[Content]; ok {
                return owner == User;
            }

            return false
        
        case AccessCanEdit, AccessCanDelete:
            if user, ok := contentOwners[Content]; ok {
                return owner == User;
            }

            return false
		}
	}, nil
}

func main() {
    // init cbac
    var cbac = InitCBAC(
        policiesSetter,
        AccessCanView,
        AccessCanEdit,
        AccessCanDelete,
    )

    // use it

    // by list
    policies, err := cbac.GetPolicies([]Content{1, 2}, "foo@bar.com")
    if err != nil {
		// error handling
	}
    for _, policy := range policies {
        if policy[AccessCanView] {
            // do something
        }
    }


    // by policy
    policy, err := cbac.GetPolicy([]Content{1, 2}, "foo@bar.com", AccessCanView, AccessCanEdit)
    if err != nil {
		// error handling
	}
    if policy[AccessCanView] || policy[AccessCanEdit] {
        // do something
    }


    // by access
    has, err := cbac.GetAccess(1, "foo@bar.com", AccessCanView)
	if err != nil {
		// error handling
	}

    if has {
        // user has AccessCanView access to the content
    } else {
        // otherwise
    }
}
```
