// Volume Management

package shockwave

import (
	"fmt"

	"gopkg.in/redis.v3"
)

type VolumeManager struct {
	ReidsClient  *redis.Client
	RedisChannel *string
	MaxVolume    *int
	MinVolume    *int
}

func NewVolumeManager(r *redis.Client, c *string, max *int, min *int) *VolumeManager {
	return &VolumeManager{
		ReidsClient:  r,
		RedisChannel: c,
		MaxVolume:    max,
		MinVolume:    min,
	}
}

func Foo() {
	fmt.Println("Foo")
}
