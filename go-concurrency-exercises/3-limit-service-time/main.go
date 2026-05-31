//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	var done chan bool = make(chan bool)
	start := time.Now()

	go func() {
		process()
		done <- true
	}()

	for {
		select {
		case <-done:
			return true
		default:
			elapsed := int64(time.Since(start).Seconds())
			atomic.AddInt64(&u.TimeUsed, atomic.LoadInt64(&u.TimeUsed)+elapsed)

			if u.TimeUsed > 10 && !u.IsPremium {
				fmt.Printf("User %v - Time %v\n", u.ID, u.TimeUsed)
				return false
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}

}

func main() {
	RunMockServer()
}
