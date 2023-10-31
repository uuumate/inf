package sql

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	currentPath        string
	getCurrentPathOnce sync.Once
)

func sqlFormat(v ...interface{}) string {
	if len(v) != 6 {
		return ""
	}

	caller := v[1]
	longMsDuration := v[2].(time.Duration)
	sql := v[3]
	rowsAffected := v[5]

	return fmt.Sprintf("(%s) %s [%0.2fms] [%d rows affected]", caller.(string)[len(getCurrentPath()):], sql, float64(longMsDuration.Microseconds())/10e2, rowsAffected)
}

func getCurrentPath() string {
	getCurrentPathOnce.Do(func() {
		currentPath, _ = os.Getwd()
		currentPath = currentPath + "/"
	})

	return currentPath
}
