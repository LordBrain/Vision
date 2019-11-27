package cmd

//Album information
type Album struct {
	Name        string
	BetterName  string
	Path        string
	ParentAlbum string
	SubAlbum    []SubAlbum
	AlbumImages []AlbumImages
}

//SubAlbum information
type SubAlbum struct {
	Name        string
	BetterName  string
	PathName    string
	AlbumImages []AlbumImages
	ImageCount  int
}

//AlbumImages name
type AlbumImages struct {
	Name string
}

//Folders information
type Folders struct {
	Path       string
	Name       string
	ParentDir  string
	ParentName string
	Remove     bool
}
