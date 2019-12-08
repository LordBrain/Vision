package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/LordBrain/Vision/cmd"

	"github.com/google/logger"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("vision", "Tool to create static image albums")

	create     = app.Command("create", "Create a new album")
	createPath = create.Arg("path", "Path to images").Required().ExistingDir()
	// createTemplate = create.Flag("template", "Template to use").String()
	createWidth = create.Flag("width", "Resized image width").Short('w').Int()

	update         = app.Command("update", "Update existing album")
	updatePath     = update.Arg("path", "Path to album").Required().ExistingDir()
	updateTemplate = update.Flag("template", "Template to use").String()
)

func main() {
	defer logger.Init("visonLog", true, false, ioutil.Discard).Close()
	logger.SetFlags(log.LstdFlags)
	// kingpin.Parse()

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// Create new album
	case create.FullCommand():
		logger.Info("Creating Albums. If there are a lot of images, this may take a while.")
		path, _ := filepath.Abs(*createPath)
		startPath := path
		logger.Info("Albums Root Path: ", startPath)
		folders := cmd.GetFolders(path)

		allAlbums := cmd.GenAlbums(startPath, folders)
		logger.Info("Number of Albums: ", len(allAlbums))

		for _, album := range allAlbums {
			logger.Info("Album Name: ", album.BetterName)
			logger.Infof("Number of images in %s: %d", album.BetterName, len(album.AlbumImages))

			//Resize images and create thumbnails
			for _, imageName := range album.AlbumImages {
				imagePath := filepath.Join(album.Path, imageName.Name)
				if strconv.Itoa(*createWidth) != "0" {
					cmd.GenImages(imagePath, *createWidth)
				} else {
					cmd.GenImages(imagePath, 800)
				}

			}

			//Create html files

			cmd.CreateIndex(album)

			fmt.Println("----")
		}

	// Update album
	case update.FullCommand():
		println("Updating")
	}
}
