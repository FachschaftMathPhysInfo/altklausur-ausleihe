package main

import (
	"log"

	api "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func applyWatermark() {
	onTop := true
	update := false

	// Stamp all odd pages of the pfd in red at the right border of the document
	wm, err := api.TextWatermark("test123 1. Mai 2021", "font:Courier, points:40, col: 1 0 0, rot:-90, sc:1 abs, opacity:0.4, pos: l, offset: -190 0", onTop, update, pdfcpu.POINTS)
	api.AddWatermarksFile("in.pdf", "", []string{"odd"}, wm, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// add the mathphys logo to the top-left corner? :P
	wm, err = api.PDFWatermark("MathPhysLogoInfo.pdf", "pos:tr, rot:0, sc:0.5 abs, offset: -10 -10, opacity:0.5", onTop, update, pdfcpu.POINTS)
	api.AddWatermarksFile("in.pdf", "", nil, wm, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	applyWatermark()
}
