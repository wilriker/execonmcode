package execonmcode

import (
	"errors"
	"fmt"
	"strings"
)

type Commands []string

func (c *Commands) String() string {
	return fmt.Sprintf("%v", *c)
}

func (c *Commands) Set(v string) error {
	*c = append(*c, v)
	return nil
}

func (c *Commands) Len() int {
	return len(*c)
}

func (c *Commands) Get(i int) (string, []string, error) {
	if i >= c.Len() {
		return "", nil, errors.New("Index out of range")
	}
	cmd := (*c)[i]
	s := strings.Split(cmd, " ")
	a := make([]string, 0)
	if len(s) > 1 {
		for _, p := range s[1:] {
			pp := strings.TrimSpace(p)
			if pp != "" {
				a = append(a, p)
			}
		}
	}
	return s[0], a, nil
}
