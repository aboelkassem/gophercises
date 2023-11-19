package main

import (
	"log"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

const (
	docHeaderTitle    = "INVOICE"
	docHeaderImage    = "gopher.png"
	docHeaderContacts = "(814) 977-7566\njon@calhoun.io\ngophercises.com"
	docHeaderAddress  = "123 Fake St.\nSome Town, PA\n12345"
)

func main() {
	pdf := gofpdf.New("P", "mm", "A4", "")

	lm, tm, rm, _ := pdf.GetMargins()
	pw, _ := pdf.GetPageSize()

	pdf.SetHeaderFunc(func() {
		pdf.SetFont("Arial", "B", 36)
		titleWidth := pdf.GetStringWidth(docHeaderTitle)
		_, titleHeight := pdf.GetFontSize()
		pdf.CellFormat(titleWidth, titleHeight, docHeaderTitle, "", 0, "L", false, 0, "")

		pdf.SetFont("Arial", "B", 12)
		addressWidth := strWidth(pdf, docHeaderAddress)
		_, addressHeight := pdf.GetFontSize()
		pdf.SetXY(pw-rm-addressWidth, tm)
		pdf.MultiCell(addressWidth, addressHeight, docHeaderAddress, "", "R", false)

		contactsWidth := strWidth(pdf, docHeaderContacts)
		_, contactsHeight := pdf.GetFontSize()
		pdf.SetXY(pw-rm-addressWidth-5-contactsWidth, tm)
		pdf.MultiCell(contactsWidth, contactsHeight, docHeaderContacts, "", "R", false)

		var imageOptions gofpdf.ImageOptions
		pdf.RegisterImageOptions(docHeaderImage, imageOptions)
		imageInfo := pdf.GetImageInfo(docHeaderImage)
		imageWidth := imageInfo.Width() * titleHeight / imageInfo.Height()
		imageX := lm + titleWidth + (pw-lm-titleWidth-imageWidth-contactsWidth-5-addressWidth-rm)/2
		pdf.ImageOptions("gopher.png", imageX, tm, 0, titleHeight, false, imageOptions, 0, "")
	})

	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	if err := pdf.OutputFileAndClose("output.pdf"); err != nil {
		log.Fatalf("Failed to create pdf: %v", err)
	}
}

func strWidth(pdf *gofpdf.Fpdf, text string) float64 {
	var maxWidth float64
	for _, line := range strings.Split(text, "\n") {
		w := pdf.GetStringWidth(line)
		if w > maxWidth {
			maxWidth = w
		}
	}
	return maxWidth + 2
}
