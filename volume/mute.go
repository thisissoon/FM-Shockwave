// Handle Mute Toggling

package volume

import (
	log "github.com/Sirupsen/logrus"
	device "github.com/thisissoon/volume"
)

var CURRENT_LEVEL int

type MuteManager struct {
	Channel    chan bool
	MixerName  *string
	DeviceName *string
}

func (m *MuteManager) Run() {
	for {
		active := <-m.Channel
		if err := m.set(active); err != nil {
			log.Errorf("Failed to set mute state", err)
		}
	}
}

func (m *MuteManager) set(active bool) error {
	if active {
		// Set sound level to 0 on Mute
		log.Info("Mute Volume")
		CURRENT_LEVEL, _ = device.GetVolume(*m.DeviceName, *m.MixerName)
		device.SetVolume(*m.DeviceName, *m.MixerName, 0)
	} else {
		// Restore sound level to 0 on Mute
		log.Infof("Restore Volume to: %v", CURRENT_LEVEL)
		device.SetVolume(*m.DeviceName, *m.MixerName, CURRENT_LEVEL)
	}

	return nil
}

// Construct a new mute manager
func NewMuteManager(c chan bool, m *string, d *string) *MuteManager {
	return &MuteManager{
		Channel:    c,
		MixerName:  m,
		DeviceName: d,
	}
}
