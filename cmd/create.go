package cmd

import (
	"fmt"
	//image libraries.
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

//GetFolders Gets all the sub folders and returns a slice with the path and folder name
func GetFolders(folder string) []Folders {
	var allFolders []Folders
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		var currentFolder Folders
		if info.IsDir() {
			currentFolder.ParentDir = filepath.Dir(path)
			currentFolder.ParentName = filepath.Base(currentFolder.ParentDir)
			currentFolder.Path = path
			currentFolder.Name = info.Name()
			allFolders = append(allFolders, currentFolder)

		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	var tmpFolders []Folders

	for _, theFolders := range allFolders {
		shouldBeHidden := regexp.MustCompile(`^\.+?`).FindString(theFolders.Name)
		if theFolders.Name == "visionimg" {
			shouldBeHidden = "visonimg"
		}
		if shouldBeHidden != "" {
			//Remove folders that start with a "." or is called "visionimg"
			for theFoldersPosition, removeHidden := range allFolders {
				if strings.Contains(removeHidden.Path, theFolders.Path) {
					allFolders[theFoldersPosition].Remove = true
				}
			}
		}

	}
	for _, removed := range allFolders {
		if !removed.Remove {
			tmpFolders = append(tmpFolders, removed)
		}
	}
	return tmpFolders
}

//RootAlbum gets the root path for albums to generate.
func RootAlbum(name string) Album {
	rootAlbum := Album{Name: name}
	return rootAlbum
}

//GenAlbums finds the images and folders.
func GenAlbums(startPath string, folders []Folders) []Album {

	var allAlbums []Album
	for _, things := range folders {
		var newAlbum Album
		betterName := strings.ReplaceAll(things.Name, "_", " ")
		newAlbum.Name = things.Name
		newAlbum.BetterName = betterName
		newAlbum.Path = things.Path
		var images []AlbumImages

		dir, _ := ReadDir(newAlbum.Path)

		for _, fileThings := range dir {
			var imageName AlbumImages
			if !fileThings.IsDir() {
				// match only these file names
				if filepath.Ext(fileThings.Name()) == ".jpg" || filepath.Ext(fileThings.Name()) == ".jpeg" || filepath.Ext(fileThings.Name()) == ".png" || filepath.Ext(fileThings.Name()) == ".gif" || filepath.Ext(fileThings.Name()) == ".tiff" {
					imageName.Name = fileThings.Name()
					images = append(images, imageName)
				}
			}
			newAlbum.AlbumImages = images
		}

		if startPath != things.Path {
			newAlbum.ParentAlbum = things.ParentName
		}

		allAlbums = append(allAlbums, newAlbum)
	}
	for _, albumDetails := range allAlbums {
		for _, subalbumDetails := range allAlbums {
			if subalbumDetails.ParentAlbum == albumDetails.Name {
				var newSubalbum SubAlbum
				betterName := strings.ReplaceAll(subalbumDetails.Name, "_", " ")
				newSubalbum.Name = subalbumDetails.Name
				newSubalbum.BetterName = betterName
				newSubalbum.PathName = subalbumDetails.Path
				// pick random images from subalbumDetails, then add those to newSubalbum images
				var randomSubImage []AlbumImages
				rand.Seed(time.Now().UnixNano())
				if len(subalbumDetails.AlbumImages) >= 4 {

					randomize := rand.Perm(len(subalbumDetails.AlbumImages))
					for _, v := range randomize[:4] {
						randomSubImage = append(randomSubImage, subalbumDetails.AlbumImages[v])
					}
				} else {
					randomSubImage = append(randomSubImage, subalbumDetails.AlbumImages[0])
				}

				newSubalbum.AlbumImages = randomSubImage
				newSubalbum.ImageCount = len(subalbumDetails.AlbumImages)
				for albumNumber, addSub := range allAlbums {
					if addSub.Name == albumDetails.Name {
						allAlbums[albumNumber].SubAlbum = append(allAlbums[albumNumber].SubAlbum, newSubalbum)
					}
				}

			}
		}
	}
	return allAlbums
}

// ReadDir reads the directory named by dirname and returns
// a list of directory entries sorted by filename.
// https://flaviocopes.com/go-list-files/
func ReadDir(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}

//GenImages will crate the thumbnails and smaller sized images.
func GenImages(imagePath string, width int) error {
	//Verify the img directory exists.
	directoryPath := filepath.Dir(imagePath)
	imageName := filepath.Base(imagePath)
	imgDir := filepath.Join(directoryPath, "visionimg")
	if _, err := os.Stat(imgDir); os.IsNotExist(err) {
		os.MkdirAll(imgDir, os.ModePerm)
	}

	src, err := imaging.Open(imagePath)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
		return err
	}
	resize := imaging.Resize(src, width, 0, imaging.Lanczos)
	resizeName := filepath.Join(imgDir, "resize_"+imageName)
	thumb := imaging.Thumbnail(src, 500, 333, imaging.Lanczos)
	thumbName := filepath.Join(imgDir, "thumb_"+imageName)
	err = imaging.Save(resize, resizeName)
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
		return err
	}
	err = imaging.Save(thumb, thumbName)
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
		return err
	}
	return nil
}
