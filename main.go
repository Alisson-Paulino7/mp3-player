package main

import (
	"io"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Tocador MP3")

	var player *oto.Player
	var decoder *minimp3.Decoder
	var context *oto.Context

	playButton := widget.NewButton("Tocar", func() {
		if decoder != nil {
			go func() {
				buffer := make([]byte, 8192)
				for {
					n, err := decoder.Read(buffer)
					if err == io.EOF {
						break
					}
					if err != nil {
						break
					}
					player.Write(buffer[:n])
				}
			}()
		}
	})

	openButton := widget.NewButton("Abrir", func() {
		dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err == nil && read != nil {
				if player != nil {
					player.Close()
				}
				if decoder != nil {
					decoder.Close()
				}
				if context != nil {
					context.Close()
				}

				file, _ := os.Open(read.URI().Path())
				decoder, _ = minimp3.NewDecoder(file)
				context, _ = oto.NewContext(decoder.SampleRate, decoder.Channels, 2, 8192)
				player = context.NewPlayer()
			}
		}, myWindow).Show()
	})

	content := container.NewVBox(
		openButton,
		playButton,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(500, 900))
	myWindow.ShowAndRun()
}
