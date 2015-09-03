// Event Handler which deligates the event to the appropriate service

package event

import (
	"encoding/json"
	"fmt"
)

const (
	VOLUME_EVENT string = "volume_changed"
	MUTE_EVENT   string = "mute_changed"
)

type Event struct {
	Type string `json:"event"`
}

type VolumeEvent struct {
	Level int `json:"volume"`
}

type MuteEvent struct {
	Active bool `json:"mute"`
}

type Handler struct {
	eventChannel  chan []byte
	muteChannel   chan bool
	volumeChannel chan int
}

// Consumes event message channel, decoding the event and passing to
// the relevant service
func (h *Handler) Run() {
	for {
		msg := <-h.eventChannel
		var event Event
		if err := json.Unmarshal(msg, &event); err != nil {
			fmt.Println(err)
			continue
		}
		switch event.Type {
		// It's a volume change event
		case VOLUME_EVENT:
			var volume struct {
				Event
				VolumeEvent
			}
			if err := json.Unmarshal(msg, &volume); err != nil {
				fmt.Println(err)
				continue
			}
			// Push to the VolumeChannel for the VolumeManager to
			// set the volume: TODO: Log It
			go func(level int) {
				h.volumeChannel <- level
			}(volume.Level)
		// It's a Mute Event
		case MUTE_EVENT:
			var mute struct {
				Event
				MuteEvent
			}
			if err := json.Unmarshal(msg, &mute); err != nil {
				fmt.Println(err)
				continue
			}
			// Push to the VolumeChannel for the VolumeManager to
			// set the volume: TODO: Log It
			go func(active bool) {
				h.muteChannel <- active
			}(mute.Active)
		}
	}
}

// Constructs a new Event Handler
func NewHandler(e chan []byte, m chan bool, v chan int) *Handler {
	return &Handler{
		eventChannel:  e,
		muteChannel:   m,
		volumeChannel: v,
	}
}
