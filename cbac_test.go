package gocbac

import (
	"errors"
	"testing"
)

type Access string

type Content struct {
	ID uint64
}

type User struct {
	Email string
}

const (
	AccessCanView   Access = "can_view"
	AccessCanEdit   Access = "can_edit"
	AccessCanDelete Access = "can_delete"
)

func policiesSetter(
	ContentList []Content,
	User User,
	RequestedAccesses []Access,
) (AccessSetter[Access, Content], error) {
	if User.Email == "error@bar.com" {
		return nil, errors.New("Error in setter")
	}

	return func(Content Content, access Access) bool {
		if Content.ID == 1 && User.Email == "foo@bar.com" {
			return true
		}

		return false
	}, nil
}

var cbacInstance = InitCBAC(
	policiesSetter,
	AccessCanView,
	AccessCanEdit,
	AccessCanDelete,
)

func TestCorrectAccess(t *testing.T) {
	has, err := cbacInstance.GetAccess(Content{ID: 1}, User{Email: "foo@bar.com"}, AccessCanView)
	if err != nil {
		t.Error(err)
	}

	if !has {
		t.Error("Access value is incorrect for foo@bar.com")
	}

	has, err = cbacInstance.GetAccess(Content{ID: 1}, User{Email: "bar@foo.com"}, AccessCanView)
	if err != nil {
		t.Error(err)
	}

	if has {
		t.Error("Access value is incorrect for bar@foo.com")
	}
}

func TestIncorrectAccess(t *testing.T) {
	_, err := cbacInstance.GetAccess(Content{ID: 1}, User{Email: "foo@bar.com"}, "random-access")
	if !errors.Is(err, ErrNoAccess) {
		t.Error(err)
	}
}

func TestSetterError(t *testing.T) {
	_, err := cbacInstance.GetAccess(Content{ID: 1}, User{Email: "error@bar.com"}, AccessCanView)
	if err == nil || err.Error() != "Error in setter" {
		t.Error("Incorrect setter error")
	}
}
