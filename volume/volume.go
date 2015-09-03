// Handle setting volume level

package volume

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"

	log "github.com/Sirupsen/logrus"
	device "github.com/thisissoon/volume"
)

type VolumeManagerOpts struct {
	Channel       chan int
	MaxVolume     *int
	MinVolume     *int
	MixerName     *string
	DeviceName    *string
	PerceptorAddr *string
	Secret        *string
}

type VolumeManager struct {
	opts *VolumeManagerOpts
}

type PerceptorPayload struct {
	Level int `json:"level"`
}

// Consumes messages from the EventChannel
func (v *VolumeManager) Run() {
	for {
		level := <-v.opts.Channel
		if err := v.set(level); err != nil {
			log.Errorf("Failed to set Volume", err)
		} else {
			if err := v.put(level); err != nil {
				log.Errorf("Failed to Update Perceptor", err)
			}
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
	log.Infof("Set level to: %v%% (%v%%)", level, volume)

	// Set the Volume on the Device
	device.SetVolume(*v.opts.DeviceName, *v.opts.MixerName, volume)

	return nil
}

func (v *VolumeManager) put(level int) error {
	// Build URL
	url := fmt.Sprintf("http://%s/volume", *v.opts.PerceptorAddr)

	// Generate Payload
	payload, _ := json.Marshal(PerceptorPayload{
		Level: level,
	})

	// Create Request
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	// Generate HMAC
	mac := hmac.New(sha256.New, []byte(*v.opts.Secret))
	mac.Write(payload)
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// Add Header
	req.Header.Add("Signature", fmt.Sprintf("%s:%s", "shockwave", sig))

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	log.Infof("Response from Perceptor: %s", resp.Status)

	return nil
}

// Construct a new volume manager
func NewVolumeManager(opts *VolumeManagerOpts) *VolumeManager {
	return &VolumeManager{
		opts: opts,
	}
}
