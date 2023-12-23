package main

import (
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

func Run(tmplSrcPath string, outputPath string) {
	tmpl, err := template.ParseFiles(tmplSrcPath)
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
		},
	}
	tmpl.Execute(fp, data)
}

func main() {
	Run("tmpl.html", "output.html")
}
