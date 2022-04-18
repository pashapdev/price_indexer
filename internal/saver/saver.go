package saver

import (
	"log"
	"time"
)

type saver struct {
}

func New() *saver {
	return &saver{}
}

func (s *saver) Save(currentMinute time.Time, data string) {
	log.Printf("%d, %s\n", currentMinute.UTC().Unix(), data)
}
