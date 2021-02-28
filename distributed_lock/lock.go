package distributed_lock

// lock.go

import (
	"errors"
	"fmt"
)

// Lock - interface
type Lock interface {
	// Lock - return ErrKeyConflict if lock has been acquired
	Lock() error
	UnLock() error
}

const (
	// dslock prefix
	DSLOCKPREFIX = "gm_dslock"
)

var (
	ErrKeyConflict = errors.New("dslock key conflicts") // error key Conflicts
)

// makeKey
func makeKey(key string) string {
	return fmt.Sprintf("%s_%s", DSLOCKPREFIX, key)
}
