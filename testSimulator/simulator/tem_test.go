package simulator_test

import (
	"fmt"
	"testing"
	"time"
)

func Test001(t *testing.T) {
	ticker := time.NewTicker(time.Second * 1)
	for range ticker.C {
		time.Sleep(time.Second * 5)
		fmt.Println(len(ticker.C))
	}
}

func Myfunc()
