package cmd

import (
	"image/png"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/tjgq/sane"

	"github.com/nint8835/brscan-to-paperless/pkg/utils"
)

var testingCmd = &cobra.Command{
	Use:    "testing",
	Short:  "Command for testing purposes.",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		err := sane.Init()
		checkErr(err, "Failed to initialize SANE")
		defer sane.Exit()

		devs, err := sane.Devices()
		checkErr(err, "Failed to list SANE devices")

		var brotherDevice *sane.Device

		for _, dev := range devs {
			slog.Info("Found SANE device", "name", dev.Name, "vendor", dev.Vendor, "model", dev.Model, "type", dev.Type)

			if dev.Vendor == "Brother" {
				brotherDevice = &dev
				break
			}
		}
		if brotherDevice == nil {
			slog.Error("No Brother device found")
			os.Exit(1)
		}

		slog.Info("Using Brother device", "name", brotherDevice.Name, "model", brotherDevice.Model)

		conn, err := sane.Open(brotherDevice.Name)
		checkErr(err, "Failed to open SANE device")
		defer conn.Close()

		inf, err := conn.SetOption("source", "flatbed")
		checkErr(err, "Failed to set source option")
		slog.Info("Set source option", "info", inf)

		slog.Info("Attempting to scan...")
		image, err := conn.ReadImage()
		checkErr(err, "Failed to read image from SANE device")

		slog.Info("Successfully read image from SANE device", "bounds", image.Bounds())

		out, err := os.Create("test.png")
		checkErr(err, "Failed to create output file")
		defer utils.DeferredClose(out)

		err = png.Encode(out, image)
		checkErr(err, "Failed to encode image to PNG")

		slog.Info("Successfully wrote image to test.png")
	},
}

func init() {
	rootCmd.AddCommand(testingCmd)
}
