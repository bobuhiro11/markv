package markv

import (
	"bufio"
	"bytes"
	"github.com/disintegration/imaging"
	"github.com/kevin-cantwell/dotmatrix"
	"github.com/nfnt/resize"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/crypto/ssh/terminal"
	"image"
	"image/draw"
	"math"
	"net/http"
	"os"
)

func RenderImage(url string) string {
	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		panic(err)
	}
	defer res.Body.Close()
	img, _, err := image.Decode(res.Body)
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	dotmatrix.NewPrinter(writer, &dotmatrix.Config{
		Filter: &filter{},
		Drawer: func() draw.Drawer {
			return draw.FloydSteinberg
		}()}).Print(img)

	writer.Flush()
	return b.String()
}

type filter struct {
	table *tablewriter.Table
}

func (f *filter) Filter(img image.Image) image.Image {
	img = imaging.Invert(img)

	cols, rows := terminalDimensions()
	dx, dy := img.Bounds().Dx(), img.Bounds().Dy()

	scale := math.Min(1.0, math.Min(
		float64(cols*2)/float64(dx),
		float64(rows*4)/float64(dy)))

	width := uint(scale * float64(img.Bounds().Dx()))
	height := uint(scale * float64(img.Bounds().Dy()))

	return resize.Resize(width, height, img, resize.NearestNeighbor)
}

func terminalDimensions() (int, int) {
	cols, rows := 27, 8

	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		tw, th, err := terminal.GetSize(int(os.Stdout.Fd()))
		if err == nil {
			cols = tw / 3
			rows = (th - 1) / 3
		}
	}

	return cols, rows
}
