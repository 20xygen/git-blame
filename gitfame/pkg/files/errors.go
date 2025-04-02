package files

import "fmt"

type ErrorWalk struct {
	Err error
}

func (e ErrorWalk) Error() string {
	return fmt.Sprintf("walking the directory failed (%v)", e.Err)
}

type ErrorRelativePath struct {
	E error
}

func (e ErrorRelativePath) Error() string {
	return fmt.Sprintf("getting relative path failed (%v)", e.E)
}
