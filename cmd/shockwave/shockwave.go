// Main Package - CLI Interface

package main

import (
	"github.com/spf13/cobra"
	"github.com/thisissoon/FM-Shockwave"
)

// Flag Variable Holders
var (
	redis_address string
	redis_channel string
	max_volume    int
	min_volume    int
)

// Long Command Description
var shockWaveLongDesc = `Shockwave Manages the Volume Levels on the System`

// Cobra Base Command
var ShockWaveCmd = &cobra.Command{
	Use:   "shockwave",
	Short: "Volume Managment Service",
	Long:  shockWaveLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		shockwave.Foo()
	},
}

// Set Command Line Flags
func init() {
	ShockWaveCmd.Flags().StringVarP(&redis_address, "redis", "r", "127.0.0.1:6379", "Redis Server Address")
	ShockWaveCmd.Flags().StringVarP(&redis_channel, "channel", "c", "", "Redis Channel Name")
	ShockWaveCmd.Flags().IntVarP(&max_volume, "max_volume", "", 100, "Max Volume Level")
	ShockWaveCmd.Flags().IntVarP(&min_volume, "min_volume", "", 0, "Min Volume Level")
}

// Application Entry Point
func main() {
	ShockWaveCmd.Execute()
}
