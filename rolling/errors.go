package rolling

import "errors"

var (
	ErrRollingFileIsAlreadyClosed = errors.New("rolling file is already closed")
)
