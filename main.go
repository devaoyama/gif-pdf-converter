package main

import (
	"bytes"
	"github.com/signintech/gopdf"
	"html/template"
	"image/gif"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func convertPDF(gifFile *gif.GIF) *gopdf.GoPdf {
	rect := gopdf.Rect{
		W: float64(gifFile.Config.Width + 10),
		H: float64(gifFile.Config.Height + 10),
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: rect, Unit: gopdf.UnitPT})

	for _, frame := range gifFile.Image {
		bounds := frame.Bounds()
		imageWidth := bounds.Max.X
		imageHeight := bounds.Max.Y
		imageRect := gopdf.Rect{
			W: float64(imageWidth),
			H: float64(imageHeight),
		}
		rand.Seed(time.Now().UnixNano())
		f := new(bytes.Buffer)
		png.Encode(f, frame)
		pdf.AddPage()
		ih, _ := gopdf.ImageHolderByReader(f)
		pdf.ImageByHolder(ih, 0, 0, &imageRect)
	}
	return &pdf
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("template/view.html")
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func convertHandler(w http.ResponseWriter, r * http.Request) {
	files, reader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer files.Close()

	if reader.Header.Get("Content-Type") != "image/gif" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	gifFiles, err := gif.DecodeAll(files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pdf := convertPDF(gifFiles)

	err = pdf.Write(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func main() {

	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/convert", convertHandler)
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
