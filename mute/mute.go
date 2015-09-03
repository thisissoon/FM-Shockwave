// Handle Mute Toggling

package mute

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	device "github.com/thisissoon/volume"
)

var CURRENT_LEVEL int

type MuteManager struct {
	Channel       chan bool
	MixerName     *string
	DeviceName    *string
	PerceptorAddr *string
	Secret        *string
}

type PerceptorPayload struct {
	Active bool `json:"active"`
}

func (m *MuteManager) Run() {
	for {
		active := <-m.Channel
		if err := m.set(active); err != nil {
			log.Errorf("Failed to set mute state", err)
		} else {
			if err := m.put(active); err != nil {
				log.Errorf("Failed to Update Perceptor", err)
			}
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

func (m *MuteManager) put(active bool) error {
	// Build URL
	url := fmt.Sprintf("http://%s/mute", *m.PerceptorAddr)

	// Generate Payload
	payload, _ := json.Marshal(PerceptorPayload{
		Active: active,
	})
	log.Debugf("Mute Payload: %s", payload)

	// Create Request
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	// Generate HMAC
	mac := hmac.New(sha256.New, []byte(*m.Secret))
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

// Construct a new mute manager
func NewMuteManager(c chan bool, m *string, d *string, p *string, s *string) *MuteManager {
	return &MuteManager{
		Channel:       c,
		MixerName:     m,
		DeviceName:    d,
		PerceptorAddr: p,
		Secret:        s,
	}
}
