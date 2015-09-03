// Main Package - CLI Interface

package main

import (
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thisissoon/FM-Shockwave/event"
	"github.com/thisissoon/FM-Shockwave/mute"
	"github.com/thisissoon/FM-Shockwave/socket"
	"github.com/thisissoon/FM-Shockwave/volume"
)

// Flag Variable Holders
var (
	perceptorAddr string
	secret        string
	max_volume    int
	min_volume    int
	mixer         string
	device        string
	log_level     string
)

// Long Command Description
var shockWaveLongDesc = `Shockwave Manages the Volume Levels on the System`

// Cobra Base Command
var ShockWaveCmd = &cobra.Command{
	Use:   "shockwave",
	Short: "Volume Managment Service",
	Long:  shockWaveLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		switch log_level {
		case "debug":
			log.SetLevel(log.DebugLevel)
		case "info":
			log.SetLevel(log.InfoLevel)
		case "warn":
			log.SetLevel(log.WarnLevel)
		case "error":
			log.SetLevel(log.ErrorLevel)
		default:
			log.SetLevel(log.WarnLevel)
		}

		log.Info("Starting Shockwave")

		// Create Channels
		eventChannel := make(chan []byte)
		volumeChannel := make(chan int)
		muteChannel := make(chan bool)

		// Consume events from Perceptor
		perceptor := socket.NewPerceptorService(&perceptorAddr, &secret, eventChannel)
		go perceptor.Run()

		// Event Handler
		eventHandler := event.NewHandler(eventChannel, muteChannel, volumeChannel)
		go eventHandler.Run()

		// Volume Manager
		volumeManager := volume.NewVolumeManager(&volume.VolumeManagerOpts{
			Channel:       volumeChannel,
			MaxVolume:     &max_volume,
			MinVolume:     &min_volume,
			MixerName:     &mixer,
			DeviceName:    &device,
			PerceptorAddr: &perceptorAddr,
		})
		go volumeManager.Run()

		// Mute Manager
		muteManager := mute.NewMuteManager(muteChannel, &mixer, &device)
		go muteManager.Run()

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
	ShockWaveCmd.Flags().StringVarP(&perceptorAddr, "perceptor", "p", "perceptor.thisissoon.fm", "Perceptor Server Address")
	ShockWaveCmd.Flags().StringVarP(&secret, "secret", "s", "CHANGE_ME", "Client Secret Ket")
	ShockWaveCmd.Flags().IntVarP(&max_volume, "max_volume", "", 100, "Max Volume Level")
	ShockWaveCmd.Flags().IntVarP(&min_volume, "min_volume", "", 0, "Min Volume Level")
	ShockWaveCmd.Flags().StringVarP(&device, "device", "d", "default", "Audio Device Name")
	ShockWaveCmd.Flags().StringVarP(&mixer, "mixer", "m", "PCM", "Audio Mixer Name")
	ShockWaveCmd.Flags().StringVarP(&log_level, "log_level", "l", "debug", "Log Level")
}

// Application Entry Point
func main() {
	ShockWaveCmd.Execute()
}
