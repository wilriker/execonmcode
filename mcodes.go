package execonmcode

import (
	"errors"
	"fmt"
	"strconv"
)

type MCodes []int64

func (m *MCodes) String() string {
	return fmt.Sprintf("%v", *m)
}

func (m *MCodes) Len() int {
	return len(*m)
}

func (m *MCodes) Set(v string) error {
	mc, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return err
	}
	if mc < 0 {
		return errors.New("-mCode must be >= 0")
	}
	*m = append(*m, mc)
	return nil
}
