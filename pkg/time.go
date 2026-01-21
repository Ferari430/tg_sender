package pkg

import (
	"log"
	"time"
)

func TimeSpent(t time.Time) time.Duration {
	sec := time.Since(t)
	log.Printf("Потрачено %v секунд", sec)
	return sec
}
