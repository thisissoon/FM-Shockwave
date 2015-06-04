// Volume Management

package shockwave

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/redis.v3"
)

// Volume Manage
type VolumeManager struct {
	ReidsClient  *redis.Client
	RedisChannel *string
	MaxVolume    *int
	MinVolume    *int
}

// Constructs a new Volume Manager Type
func NewVolumeManager(r *redis.Client, c *string, max *int, min *int) *VolumeManager {
	return &VolumeManager{
		ReidsClient:  r,
		RedisChannel: c,
		MaxVolume:    max,
		MinVolume:    min,
	}
}

func (v *VolumeManager) Consume() {
	// Subscribe to channel, exiting the program on fail
	pubsub := v.RedisClient.PubSub()
	err := pubsub.Subscribe(v.RedisChannel)
	if err != nil {
		log.Fatalln(err)
	}

	// Ensure connection the channel is closed on exit
	defer pubsub.Close()

	// Loop to recieve events
	for {
		msg, err := pubsub.Receive() // recieve a message from the channel
		if err != nil {
			log.Println(err)
		} else {
			switch m := msg.(type) { // Switch the mesage type
			case *redis.Subscription:
				log.Println(fmt.Sprintf("%s: %s", strings.Title(m.Kind), m.Channel))
			case *redis.Message:
				err := v.processMessage([]byte(m.Payload))
				if err != nil {
					log.Println(err)
				}
			default:
				log.Println(fmt.Sprintf("Unknown message: %#v", m))
			}
		}
	}
}

func Foo() {
	fmt.Println("Foo")
}
