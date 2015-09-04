// Main Package - CLI Interface

package main

import (
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thisissoon/FM-Shockwave/event"
	"github.com/thisissoon/FM-Shockwave/mute"
	"github.com/thisissoon/FM-Shockwave/socket"
	"github.com/thisissoon/FM-Shockwave/volume"
)

// Long Command Description
var shockWaveLongDesc = `Shockwave Manages the Volume Levels on the System`

// Cobra Base Command
var ShockWaveCmd = &cobra.Command{
	Use:   "shockwave",
	Short: "Volume Managment Service",
	Long:  shockWaveLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Starting Shockwave")

		// Create Channels
		eventChannel := make(chan []byte)
		volumeChannel := make(chan int)
		muteChannel := make(chan bool)

		// Consume events from Perceptor
		perceptor := socket.NewPerceptorService(
			viper.GetString("perceptor_addr"),
			viper.GetString("secret"),
			eventChannel)
		go perceptor.Run()

		// Event Handler
		eventHandler := event.NewHandler(eventChannel, muteChannel, volumeChannel)
		go eventHandler.Run()

		// Volume Manager
		volumeManager := volume.NewVolumeManager(&volume.VolumeManagerOpts{
			Channel:       volumeChannel,
			MaxVolume:     viper.GetInt("max"),
			MinVolume:     viper.GetInt("min"),
			MixerName:     viper.GetString("mixer"),
			DeviceName:    viper.GetString("device"),
			PerceptorAddr: viper.GetString("perceptor_addr"),
			Secret:        viper.GetString("secret"),
		})
		go volumeManager.Run()

		// Mute Manager
		muteManager := mute.NewMuteManager(
			muteChannel,
			viper.GetString("mixer"),
			viper.GetString("device"),
			viper.GetString("perceptor_addr"),
			viper.GetString("secret"))
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
	// Load config from File
	log.SetLevel(log.WarnLevel)

	// Defaults
	viper.SetDefault("log_level", "warn")
	viper.SetDefault("perceptor_address", "localhost:9000")
	viper.SetDefault("secret", "CHANGE_ME")
	viper.SetDefault("max", 100)
	viper.SetDefault("min", 0)
	viper.SetDefault("device", "default")
	viper.SetDefault("mixer", "PCM")

	// From file
	viper.SetConfigName("config")           // name of config file (without extension)
	viper.AddConfigPath("/etc/shockwave/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.shockwave") // call multiple times to add many search paths
	viper.AddConfigPath("$PWD/.shockwave")  // call multiple times to add many search paths
	err := viper.ReadInConfig()             // Find and read the config file
	if err != nil {                         // Handle errors reading the config file
		log.Warnf("No config file found or is not properly formatted: %s", err)
	}

	// Switch Log Level
	switch viper.GetString("log_level") {
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
}

// Application Entry Point
func main() {
	ShockWaveCmd.Execute()
}
