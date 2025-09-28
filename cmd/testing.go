package cmd

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"log/slog"
	"os"

	"codeberg.org/go-pdf/fpdf"
	"github.com/fewebahr/sane"
	"github.com/spf13/cobra"
)

func imageToPNGReader(img image.Image) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		err := png.Encode(pw, img)
		if err != nil {
			_ = pw.CloseWithError(fmt.Errorf("failed to encode image to PNG: %w", err))
			return
		}
		err = pw.Close()
		if err != nil {
			slog.Error("Failed to close pipe writer", "err", err)
		}
	}()
	return pr
}

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

		imageNum := 1

		pdf := fpdf.New("P", "pt", "", "")

		for {
			slog.Info("Attempting to scan...")
			image, err := conn.ReadImage()
			if errors.Is(err, sane.ErrEmpty) {
				break
			}
			checkErr(err, "Failed to read image from SANE device")

			slog.Info("Successfully read image from SANE device", "bounds", image.Bounds())

			pdf.AddPageFormat(
				"P",
				fpdf.SizeType{
					Wd: float64(image.Bounds().Dx()),
					Ht: float64(image.Bounds().Dy()),
				},
			)
			pdf.RegisterImageOptionsReader(
				fmt.Sprintf("image-%d", imageNum),
				fpdf.ImageOptions{
					ReadDpi:   true,
					ImageType: "PNG",
				},
				imageToPNGReader(image),
			)
			pdf.ImageOptions(
				fmt.Sprintf("image-%d", imageNum),
				0,
				0,
				float64(image.Bounds().Dx()),
				float64(image.Bounds().Dy()),
				false,
				fpdf.ImageOptions{
					ReadDpi:   true,
					ImageType: "PNG",
				},
				0,
				"",
			)

			imageNum++
		}

		err = pdf.OutputFileAndClose("test-output.pdf")
		checkErr(err, "Failed to save PDF file")

		slog.Info("Successfully saved PDF file", "path", "test-output.pdf")
	},
}

func init() {
	rootCmd.AddCommand(testingCmd)
}
