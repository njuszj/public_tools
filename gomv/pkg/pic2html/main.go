package main

import (
	_ "embed"
	"html/template"
	"log"
	"os"
)

type Img struct {
	Src string
}

type ImgData struct {
	PageTitle string
	Imgs      []Img
}

func Run(tmplSrc string, outputPath string) {
	tmpl, err := template.New("myteml").Parse(tmplSrc)
	if err != nil {
		log.Fatal(err)
	}
	fp, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	data := ImgData{
		PageTitle: "Demo",
		Imgs: []Img{
			{Src: "DSC_0164.JPG"},
			{Src: "DSC_0181.JPG"},
			{Src: "DSC_0182.JPG"},
			{Src: "DSC_0183.JPG"},
			{Src: "DSC_0185.JPG"},
		},
	}
	tmpl.Execute(fp, data)
}

//go:embed tmpl.html
var templateHTML string

func main() {
	Run(templateHTML, "output.html")
}
