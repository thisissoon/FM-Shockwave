// Volume Management

package shockwave

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/bklimt/volume"
	"gopkg.in/redis.v3"
)

// Holds the current volume level
var CURRENT_LEVEL int

const (
	SET_VOLUME_EVENT    string = "set_volume"
	CHANGE_VOLUME_EVENT string = "volume_changed"
)

// JSON Structure for a Volume Change event
type VolumeEvent struct {
	Event  string `json:"event"`
	Volume int    `json:"volume"`
}

// Volume Manage
type VolumeManager struct {
	RedisClient  *redis.Client
	RedisChannel *string
	MaxVolume    *int
	MinVolume    *int
}

// Constructs a new Volume Manager Type
func NewVolumeManager(r *redis.Client, c *string, max *int, min *int) *VolumeManager {
	return &VolumeManager{
		RedisClient:  r,
		RedisChannel: c,
		MaxVolume:    max,
		MinVolume:    min,
	}
}

func (v *VolumeManager) Subscribe() {
	// Subscribe to channel, exiting the program on fail
	pubsub := v.RedisClient.PubSub()
	err := pubsub.Subscribe(*v.RedisChannel)
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

// Process Volume Event
func (v *VolumeManager) processMessage(m []byte) error {
	event := &VolumeEvent{}
	err := json.Unmarshal(m, event)
	if err != nil {
		return err
	}

	if event.Event == SET_VOLUME_EVENT {
		if err := v.setVolume(event.Volume); err != nil {
			log.Println(err)
		}
	}

	return nil
}

// Set the Volume on the device & Publishes Volume Changed Event
func (v *VolumeManager) setVolume(l int) error {
	// Validate intended level
	if l > 100 || l < 0 {
		return errors.New(fmt.Sprintf("%v is not between 0 and 100", l))
	}

	// Convert out ints to floats
	min := float64(*v.MinVolume)
	max := float64(*v.MaxVolume)
	vol := float64(l) // Percentage

	// Calculate the adjusted volume level - Rounding to the nearest whole number
	actual := int(math.Floor((vol*((max-min)/100) + min) + .5))
	log.Println(actual)
	log.Println(fmt.Sprintf("Set level to: %v%% (%v%%)", vol, actual))

	// Store the new Volume Level
	CURRENT_LEVEL = actual

	// Set the Volume on the Device
	volume.SetVolume("PCM", actual)

	// Create Messages
	message, err := json.Marshal(&VolumeEvent{
		Event:  CHANGE_VOLUME_EVENT,
		Volume: l,
	})

	// Publish Change Event
	log.Println("Publish Volume Change Event")
	err = v.RedisClient.Publish(*v.RedisChannel, string(message[:])).Err()
	if err != nil {
		return err
	}

	return nil
}
