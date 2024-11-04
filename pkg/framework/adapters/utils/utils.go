package utils

import (
	"fmt"
	"io"
)

func CloseResource(c io.Closer, err *error) {
	if cerr := c.Close(); cerr != nil {
		if *err == nil {
			*err = cerr
		} else {
			*err = fmt.Errorf("%v; %v", *err, cerr)
		}
	}
}
