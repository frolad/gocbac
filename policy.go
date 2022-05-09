package gocbac

type Policies[A comparable, C comparable] map[C]Policy[A]

func (p Policies[A, C]) Set(setter func(content C, access A) bool) Policies[A, C] {
	for content, policy := range p {
		for access := range policy {
			p[content][access] = setter(content, access)
		}
	}

	return p
}

// policy is the simple map of accesses and their bool values
type Policy[A comparable] map[A]bool
