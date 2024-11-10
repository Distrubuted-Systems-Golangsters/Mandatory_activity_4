package peer

import (
	"log"
	"time"
)

func CriticalSection() {
	time.Sleep(10 * time.Second)
	log.Printf("Hello world\n")
}