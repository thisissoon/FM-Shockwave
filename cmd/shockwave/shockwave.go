// Main Package - CLI Interface

package main

import (
	"github.com/spf13/cobra"
	"github.com/thisissoon/FM-Shockwave"
)

var shockWaveLongDesc = `Shockwave Manages the Volume Levels on the System`

var ShockWaveCmd = &cobra.Command{
	Use:   "shockwave",
	Short: "Volume Managment Service",
	Long:  shockWaveLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		shockwave.Foo()
	},
}

func main() {
	ShockWaveCmd.Execute()
}
