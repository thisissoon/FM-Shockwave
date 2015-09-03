// Handle setting volume level

package volume

import (
	"errors"
	"fmt"
	"log"
	"math"

	device "github.com/thisissoon/volume"
)

type volumeEvent struct {
	Event  string `json:"event"`
	Volume int    `json:"volume"`
}

type VolumeManagerOpts struct {
	Channel    chan int
	MaxVolume  *int
	MinVolume  *int
	MixerName  *string
	DeviceName *string
}

type VolumeManager struct {
	opts *VolumeManagerOpts
}

// Consumes messages from the EventChannel
func (v *VolumeManager) Run() {
	for {
		level := <-v.opts.Channel
		if err := v.set(level); err != nil {
			fmt.Println(err)
		}
	}
}

// Set the Volume to the correct level
func (v *VolumeManager) set(i int) error {
	// Convert out ints to floats
	min := float64(*v.opts.MinVolume)
	max := float64(*v.opts.MaxVolume)
	level := float64(i) // Percentage

	// Validate intended level
	if level > 100 || level < 0 {
		return errors.New(fmt.Sprintf("%v is not between 0 and 100", level))
	}

	// Calculate the adjusted volume level - Rounding to the nearest whole number
	volume := int(math.Floor((level*((max-min)/100) + min) + .5))
	log.Println(fmt.Sprintf("Set level to: %v%% (%v%%)", level, volume))

	// Set the Volume on the Device
	device.SetVolume(*v.opts.DeviceName, *v.opts.MixerName, volume)

	return nil
}

// Construct a new volume manager
func NewVolumeManager(opts *VolumeManagerOpts) *VolumeManager {
	return &VolumeManager{
		opts: opts,
	}
}
