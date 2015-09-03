// Event Handler which deligates the event to the appropriate service

package event

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
)

const (
	VOLUME_EVENT string = "set_volume"
	MUTE_EVENT   string = "set_mute"
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
			log.Errorf("Failed to unmarshal event", err)
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
				log.Errorf("Failed to unmarshal Volume event", err)
				continue
			}
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
				log.Errorf("Failed to unmarshal Mute event", err)
				continue
			}
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
