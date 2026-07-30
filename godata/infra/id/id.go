package id

import (
	"sync"

	"github.com/sony/sonyflake"
)

var (
	sf   *sonyflake.Sonyflake
	once sync.Once
)

func initSonyflake() {
	once.Do(func() {
		sf = sonyflake.NewSonyflake(sonyflake.Settings{})
	})
}

func GenerateID() (uint64, error) {
	initSonyflake()
	return sf.NextID()
}

func MustGenerateID() uint64 {
	id, err := GenerateID()
	if err != nil {
		panic("failed to generate ID: " + err.Error())
	}
	return id
}
