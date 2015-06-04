// Volume Management

package shockwave

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/thisissoon/volume"
	"gopkg.in/redis.v3"
)

// Holds the current volume level
var CURRENT_LEVEL int

// State Keys
const (
	VOLUME_STATE_KEY string = "fm:player:volume"
	MUTE_STATE_KEY   string = "fm:player:mute"
)

// Event Names
const (
	SET_VOLUME_EVENT    string = "set_volume"
	SET_MUTE_EVENT      string = "set_mute"
	CHANGE_VOLUME_EVENT string = "volume_changed"
	CHANGE_MUTE_EVENT   string = "mute_changed"
)

// Base Event
type Event struct {
	Event string `json:"event"`
	Value json.RawMessage
}

// JSON Structure for a Volume Change event
type VolumeEvent struct {
	Event  string `json:"event"`
	Volume int    `json:"volume"`
}

// JSON Structure for a Mute Change event
type MuteEvent struct {
	Event string `json:"event"`
	Mute  bool   `json:"mute"`
}

// Volume Manage
type VolumeManager struct {
	RedisClient  *redis.Client
	RedisChannel *string
	MaxVolume    *int
	MinVolume    *int
	MixerName    *string
	DeviceName   *string
}

// Constructs a new Volume Manager Type
func NewVolumeManager(r *redis.Client, c *string, max *int, min *int, mixer *string, device *string) *VolumeManager {
	return &VolumeManager{
		RedisClient:  r,
		RedisChannel: c,
		MaxVolume:    max,
		MinVolume:    min,
		MixerName:    mixer,
		DeviceName:   device,
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
	event := &Event{}
	err := json.Unmarshal(m, event)
	if err != nil {
		return err
	}

	// Switch the raw event type
	switch event.Event {
	case SET_VOLUME_EVENT:
		if err := v.setVolume(&m); err != nil {
			log.Println(err)
		}
	case SET_MUTE_EVENT:
		if err := v.setMute(&m); err != nil {
			log.Println(err)
		}
	}

	return nil
}

// Set the Volume on the device & Publishes Volume Changed Event
func (v *VolumeManager) setVolume(m *[]byte) error {
	var err error

	// Unmarshal the JSON
	event := &VolumeEvent{}
	err = json.Unmarshal(*m, event)
	if err != nil {
		return err
	}

	// Convert out ints to floats
	min := float64(*v.MinVolume)
	max := float64(*v.MaxVolume)
	vol := float64(event.Volume) // Percentage

	// Validate intended level
	if vol > 100 || vol < 0 {
		return errors.New(fmt.Sprintf("%v is not between 0 and 100", vol))
	}

	// Calculate the adjusted volume level - Rounding to the nearest whole number
	actual := int(math.Floor((vol*((max-min)/100) + min) + .5))
	log.Println(fmt.Sprintf("Set level to: %v%% (%v%%)", vol, actual))

	// Store the new Volume Level
	CURRENT_LEVEL = actual

	// Set the Volume on the Device
	volume.SetVolume(*v.DeviceName, *v.MixerName, actual)

	// Create Messages
	message, err := json.Marshal(&VolumeEvent{
		Event:  CHANGE_VOLUME_EVENT,
		Volume: event.Volume,
	})

	// Set Volume State Redis Key
	err = v.RedisClient.Set(VOLUME_STATE_KEY, strconv.Itoa(event.Volume), 0).Err()
	if err != nil {
		return err
	}

	// Publish Change Event
	log.Println("Publish Volume Change Event")
	err = v.RedisClient.Publish(*v.RedisChannel, string(message[:])).Err()
	if err != nil {
		return err
	}

	// Set Mute State to 0
	if err := v.publishMuteChangeEvent(false); err != nil {
		return err
	}

	return nil
}

// Set the mute level
func (v *VolumeManager) setMute(m *[]byte) error {
	var err error

	// Unmarshal the JSON
	event := &MuteEvent{}
	err = json.Unmarshal(*m, event)
	if err != nil {
		return err
	}

	if event.Mute {
		// Set sound level to 0 on Mute
		CURRENT_LEVEL, _ = volume.GetVolume(*v.DeviceName, *v.MixerName)
		volume.SetVolume(*v.DeviceName, *v.MixerName, 0)
	} else {
		// Restore sound level to 0 on Mute
		volume.SetVolume(*v.DeviceName, *v.MixerName, CURRENT_LEVEL)
	}

	// Set Mute State
	if err := v.publishMuteChangeEvent(event.Mute); err != nil {
		return err
	}

	return nil
}

// Publush Mute State
func (v *VolumeManager) publishMuteChangeEvent(state bool) error {
	var err error

	log.Println(fmt.Sprintf("Set %v: %v", MUTE_STATE_KEY, state))
	// Set Mute State Redis Key to value of state
	err = v.RedisClient.Set(MUTE_STATE_KEY, strconv.Itoa(Btoi(state)), 0).Err()
	if err != nil {
		return err
	}

	// Create Messages
	message, err := json.Marshal(&MuteEvent{
		Event: CHANGE_MUTE_EVENT,
		Mute:  state,
	})
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("Publish %v: %v", CHANGE_MUTE_EVENT, state))
	err = v.RedisClient.Publish(*v.RedisChannel, string(message[:])).Err()
	if err != nil {
		return err
	}

	return nil
}

// Convert bool to int
func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
