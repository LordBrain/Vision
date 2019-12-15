package cmd

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobuffalo/packr"
	"github.com/google/logger"
)

//HTML template stuff goes here.

//CreateIndex pulls in the template to the package.
func CreateIndex(album Album) {
	box := packr.NewBox("../template")

	if _, err := os.Stat(album.Path + "/visionimg"); os.IsNotExist(err) {
		os.Mkdir(album.Path+"/visionimg", 0777)
	}

	placeHolder, err := box.Find("album-placeholder.jpg")
	if err != nil {
		logger.Error("here?", err)
	}
	placeholderFile, err := os.Create(album.Path + "/visionimg/album-placeholder.jpg")
	if err != nil {
		logger.Error(err)
	}
	_, err = placeholderFile.Write(placeHolder)
	// err = ioutil.WriteFile(album.Path+"/visonimg/album-placeholder.jpg", placeHolder, 0777)
	if err != nil {
		logger.Error(err)
	}
	placeholderFile.Close()

	t, err := template.New("index").Parse(box.String("index.tmpl"))

	if err != nil {
		logger.Error(err)
	}
	indexFile := filepath.Join(album.Path, "index.html")
	f, err := os.Create(indexFile)
	if err != nil {
		logger.Error(err)
	}
	err = t.Execute(f, album)
	if err != nil {
		logger.Error(err)
	}
	f.Close()

	// Read files in styles directory
	stylesheets, err := ioutil.ReadDir("./template/styles")
	if err != nil {
		logger.Error(err)
	}

	// Loop through files and add them to the visionimg folder
	for _, file := range stylesheets {

		stylesheet := []byte(box.String("styles/" + file.Name()))
		if err != nil {
			logger.Error(err)
		}
		stylesheetFilePath := filepath.Join(album.Path, "visionimg/"+file.Name())
		err = ioutil.WriteFile(stylesheetFilePath, stylesheet, 0777)
		if err != nil {
			logger.Error(err)
		}
	}
}
