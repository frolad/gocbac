package gocbac

import (
	"errors"
	"fmt"

	"golang.org/x/exp/maps"
)

// Possible errors
// ErrNoContent in case there is no such content in the policy
// ErrNoAccess in case there is no such access in the policy(s)
var (
	ErrNoContent = errors.New("no such content")
	ErrNoAccess  = errors.New("no such access")
)

// PoliciesSetter function which populates policies for the list of content, error will be passed to the executor (e.g. GetPolicies)
// AccessSetter function which populates access for the peace of content
type (
	PoliciesSetter[A, C comparable, O any] func(ContentList []C, On O, requestedAccesses []A) (AccessSetter[A, C], error)
	AccessSetter[A, C comparable]          func(content C, access A) bool
)

// GetPolicies get the list of policies for the list of content on On instance (optionally for the list of accesses)
// GetPolicy get the policy for the content on On instance (optionally for the list of accesses)
// GetAccess get the access for the content on On instance for the specific access
type CBAC[A, C comparable, O any] interface {
	GetPolicies(ContentList []C, On O, Accesses ...A) (Policies[A, C], error)
	GetPolicy(Content C, On O, Accesses ...A) (Policy[A], error)
	GetAccess(Content C, On O, Access A) (bool, error)
}

// cbac stores policies setter as well as list of accesses
type cbac[A, C comparable, O any] struct {
	accesses map[A]bool
	setter   PoliciesSetter[A, C, O]
}

// Init CBAC where:
// A - type of access
// C - type of content
// O - type of instance of which access should be checked
// Setter - policies setter function
// Accesses - list of needed accesses
func InitCBAC[A, C comparable, O any](Setter PoliciesSetter[A, C, O], Accesses ...A) CBAC[A, C, O] {
	return &cbac[A, C, O]{
		accesses: SliceToBoolMap(Accesses),
		setter:   Setter,
	}
}

// GetPolicies get the list of policies for the list of content on On instance (optionally for the list of accesses)
func (c *cbac[A, C, O]) GetPolicies(ContentList []C, On O, requestedAccesses ...A) (Policies[A, C], error) {
	possibleAccesses, err := c.cleanUpReqeustAccesses(requestedAccesses)
	if err != nil {
		return Policies[A, C]{}, err
	}

	// get initial policies
	policies, err := c.preparePolicies(ContentList, possibleAccesses)
	if err != nil {
		return policies, err
	}

	acessSetter, err := c.setter(ContentList, On, possibleAccesses)
	if err != nil {
		return policies, err
	}

	for content, policy := range policies {
		for access := range policy {
			policies[content][access] = acessSetter(content, access)
		}
	}

	// remove unrequested acceses in case getter have set them
	return c.cleanUpPolicies(policies, possibleAccesses), nil
}

// GetPolicy get the policy for the content on On instance (optionally for the list of accesses)
func (c *cbac[A, C, O]) GetPolicy(Content C, On O, Accesses ...A) (Policy[A], error) {
	policies, err := c.GetPolicies([]C{Content}, On, Accesses...)
	if err != nil {
		return Policy[A]{}, err
	}

	if policy, ok := policies[Content]; ok {
		return policy, nil
	}

	return Policy[A]{}, ErrNoContent
}

// GetAccess get the access for the content on On instance for the specific access
func (c *cbac[A, C, O]) GetAccess(Content C, On O, Access A) (bool, error) {
	policies, err := c.GetPolicies([]C{Content}, On, []A{Access}...)
	if err != nil {
		return false, err
	}

	if policy, ok := policies[Content]; ok {
		return policy[Access], nil
	}

	return false, fmt.Errorf("%w: "+fmt.Sprintf("%v", Access), ErrNoAccess)
}

// clean up requested accesses
// in case requested access is not in the list of original accesses - return error
func (c *cbac[A, C, O]) cleanUpReqeustAccesses(requestedAccesses []A) ([]A, error) {
	keys := []A{}

	if len(requestedAccesses) > 0 {
		for _, access := range requestedAccesses {
			if _, ok := c.accesses[access]; ok {
				keys = append(keys, access)
			} else {
				return keys, fmt.Errorf("%w: "+fmt.Sprintf("%v", access), ErrNoAccess)
			}
		}
	} else {
		keys = maps.Keys(c.accesses)
	}

	return keys, nil
}

// created empty list of policies for the list of contents
func (c *cbac[A, C, O]) preparePolicies(ContentList []C, requestedAccesses []A) (Policies[A, C], error) {
	res := Policies[A, C]{}
	for _, id := range ContentList {
		res[id] = MapFill(Policy[A]{}, requestedAccesses, false)
	}

	return res, nil
}

// clean up policies - keep only accesses which were provided in the InitCBAC and ignore the rest
func (c *cbac[A, C, O]) cleanUpPolicies(policies Policies[A, C], requestedAccesses []A) Policies[A, C] {
	for key, policy := range policies {
		cleanPolicy := MapFill(Policy[A]{}, requestedAccesses, false)

		for _, access := range requestedAccesses {
			if val, ok := policy[access]; ok {
				cleanPolicy[access] = val
			}
		}

		policies[key] = cleanPolicy
	}

	return policies
}
