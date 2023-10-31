package main

import (
	"fmt"
	"time"

	"github.com/uuumate/inf/rolling"
)

func main() {
	rolling := rolling.NewRollingFile("logs", "access")
	for i := 0; i < 10; i++ {
		_, _ = rolling.Write([]byte(fmt.Sprintf(`{"a":%d}`, i)))
	}
	_ = rolling.Sync()
	time.Sleep(time.Second * 3)
}
