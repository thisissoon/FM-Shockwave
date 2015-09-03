// Main Package - CLI Interface

package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/thisissoon/FM-Shockwave/event"
	"github.com/thisissoon/FM-Shockwave/socket"
	"github.com/thisissoon/FM-Shockwave/volume"
)

// Flag Variable Holders
var (
	perceptor_addr string
	secret         string
	max_volume     int
	min_volume     int
	mixer          string
	device         string
)

// Long Command Description
var shockWaveLongDesc = `Shockwave Manages the Volume Levels on the System`

// Cobra Base Command
var ShockWaveCmd = &cobra.Command{
	Use:   "shockwave",
	Short: "Volume Managment Service",
	Long:  shockWaveLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting Shockwave")
		// Channel to listen for OS Signals
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill)
		// Run for ever unless we get a signal
		for sig := range signals {
			log.Println(sig)
			os.Exit(1)
		}
	},
}

// Set Command Line Flags
func init() {
	ShockWaveCmd.Flags().StringVarP(&perceptor_addr, "perceptor", "p", "perceptor.thisissoon.fm", "Perceptor Server Address")
	ShockWaveCmd.Flags().StringVarP(&secret, "secret", "s", "CHANGE_ME", "Client Secret Ket")
	ShockWaveCmd.Flags().IntVarP(&max_volume, "max_volume", "", 100, "Max Volume Level")
	ShockWaveCmd.Flags().IntVarP(&min_volume, "min_volume", "", 0, "Min Volume Level")
	ShockWaveCmd.Flags().StringVarP(&device, "device", "d", "default", "Audio Device Name")
	ShockWaveCmd.Flags().StringVarP(&mixer, "mixer", "m", "PCM", "Audio Mixer Name")
}

// Application Entry Point
func main() {
	// Create Channels
	eventChannel := make(chan []byte)
	volumeChannel := make(chan int)
	muteChannel := make(chan bool)

	// Consume events from Perceptor
	perceptor := socket.NewPerceptorService(&perceptor_addr, &secret, eventChannel)
	go perceptor.Run()

	// Event Handler
	eventHandler := event.NewHandler(eventChannel, muteChannel, volumeChannel)
	go eventHandler.Run()

	// Volume Manager
	volumeManager := volume.NewVolumeManager(&volume.VolumeManagerOpts{
		Channel:    volumeChannel,
		MaxVolume:  &max_volume,
		MinVolume:  &min_volume,
		MixerName:  &mixer,
		DeviceName: &device,
	})
	go volumeManager.Run()

	// Mute Manager
	muteManager := volume.NewMuteManager(muteChannel, &mixer, &device)
	go muteManager.Run()

	ShockWaveCmd.Execute()
}
