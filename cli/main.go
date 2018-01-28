package main

import (
	"github.com/nmi/markv"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"os"
)

func convert(in string) string {
	b := []byte(in)
	r := &markv.Render{}

	return string(blackfriday.Run(
		b,
		blackfriday.WithExtensions(blackfriday.CommonExtensions),
		blackfriday.WithRenderer(r)))
}

func main() {
	str, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	b := []byte(str)
	r := &markv.Render{}

	if _, err := os.Stdout.Write(blackfriday.Run(
		b,
		blackfriday.WithExtensions(blackfriday.CommonExtensions),
		blackfriday.WithRenderer(r))); err != nil {
		panic(err)
	}
}
