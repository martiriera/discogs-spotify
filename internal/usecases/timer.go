package usecases

import (
	"log"
	"time"
)

func StartTimer(name string) func() {
	t := time.Now()
	return func() {
		d := time.Since(t)
		log.Println(name, "took", d)
	}
}
