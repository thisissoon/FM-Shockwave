// Main Package - CLI Interface

package main

import (
	"log"
	"os"
	"os/signal"

	"gopkg.in/redis.v3"

	"github.com/spf13/cobra"
	"github.com/thisissoon/FM-Shockwave"
)

// Flag Variable Holders
var (
	redis_address string
	redis_channel string
	max_volume    int
	min_volume    int
	mixer         string
	device        string
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

		// Create a redis client
		redis_client := redis.NewClient(&redis.Options{
			Network: "tcp",
			Addr:    redis_address,
		})

		// Create Volume Manager
		v := shockwave.NewVolumeManager(
			redis_client,
			&redis_channel,
			&max_volume,
			&min_volume,
			&mixer,
			&device)
		go v.Subscribe() // Subscribe to the redis Pub/Sub channel and consume the messages

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
	ShockWaveCmd.Flags().StringVarP(&redis_address, "redis", "r", "127.0.0.1:6379", "Redis Server Address")
	ShockWaveCmd.Flags().StringVarP(&redis_channel, "channel", "c", "", "Redis Channel Name")
	ShockWaveCmd.Flags().IntVarP(&max_volume, "max_volume", "", 100, "Max Volume Level")
	ShockWaveCmd.Flags().IntVarP(&min_volume, "min_volume", "", 0, "Min Volume Level")
	ShockWaveCmd.Flags().StringVarP(&device, "device", "d", "default", "Audio Device Name")
	ShockWaveCmd.Flags().StringVarP(&mixer, "mixer", "m", "PCM", "Audio Mixer Name")
}

// Application Entry Point
func main() {
	ShockWaveCmd.Execute()
}
