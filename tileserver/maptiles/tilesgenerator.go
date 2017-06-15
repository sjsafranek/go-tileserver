package maptiles

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
)

// TileCache struct for PostgreSQL MBTile database.
// MBTiles 1.2-compatible Tile Db with multi-layer support.
// Was named Mbtiles before, hence the use of *m in methods.
type TileCache struct {
	requestChan chan TileFetchRequest
	insertChan  chan TileFetchResult
	layerIds    map[string]int
	qc          chan bool
}

// NewTileCache creates TileCache struct.
// Creates database tables and initializes tile request channels.
func NewTileCache(path string) *TileCache {
	m := TileCache{}
	m.insertChan = make(chan TileFetchResult)
	m.requestChan = make(chan TileFetchRequest)
	go m.Run()
	return &m
}

// exists returns whether the given file or directory exists or not
func (self *TileCache) exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (self *TileCache) createDirectory(filepath string) {
	os.MkdirAll(filepath, os.ModePerm)
}

// Close tile request channels.
func (self *TileCache) Close() {
	close(self.insertChan)
	close(self.requestChan)
	if self.qc != nil {
		<-self.qc // block until channel qc is closed (meaning Run() is finished)
	}
}

// InsertQueue gets tile insert channel.
func (self TileCache) InsertQueue() chan<- TileFetchResult {
	return self.insertChan
}

// RequestQueue gets tile request channel.
func (self TileCache) RequestQueue() chan<- TileFetchRequest {
	return self.requestChan
}

// Run runs tile generation.
// Best executed in a dedicated go routine.
func (self *TileCache) Run() {
	self.qc = make(chan bool)
	for {
		select {
		case r := <-self.requestChan:
			self.fetch(r)
		case i := <-self.insertChan:
			self.insert(i)
		}
	}
	self.qc <- true
}

// insert tile request into database table.
func (self *TileCache) insert(i TileFetchResult) {
	i.Coord.setTMS(true)
	x, y, zoom, l := i.Coord.X, i.Coord.Y, i.Coord.Zoom, i.Coord.Layer

	filepath := fmt.Sprintf("cache/%v/%v/%v/", l, zoom, x)
	if !self.exists(filepath) {
		self.createDirectory(filepath)
	}

	img, _, _ := image.Decode(bytes.NewReader(i.BlobPNG))

	imgFileName := fmt.Sprintf("./cache/%v/%v/%v/%v.png", l, zoom, x, y)
	out, err := os.Create(imgFileName)
	if nil != err {
		Ligneous.Error("error during insert", err)
		return
	}

	err = png.Encode(out, img)
	if nil != err {
		Ligneous.Error("error during insert", err)
		return
	}

	Ligneous.Trace(fmt.Sprintf("INSERT BLOB %v %v %v %v", l, zoom, x, y))
}

// fetch gets cached tile from database.
func (self *TileCache) fetch(r TileFetchRequest) {
	// FIND ON DISK
	r.Coord.setTMS(true)
	zoom, x, y, l := r.Coord.Zoom, r.Coord.X, r.Coord.Y, r.Coord.Layer
	result := TileFetchResult{r.Coord, nil}

	imgFileName := fmt.Sprintf("./cache/%v/%v/%v/%v.png", l, zoom, x, y)
	if self.exists(imgFileName) {

		dat, err := ioutil.ReadFile(imgFileName)
		if nil != err {
			result.BlobPNG = nil
		} else {
			result.BlobPNG = dat
			Ligneous.Trace(fmt.Sprintf("REUSE BLOB %v %v %v %v", l, zoom, x, y))
		}

	} else {
		result.BlobPNG = nil
	}

	r.OutChan <- result
}
