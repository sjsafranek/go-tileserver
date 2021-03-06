package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

import "maptiles"

type Config struct {
	Cache  string            `json:"cache"`
	Engine string            `json:"engine"`
	Layers map[string]string `json:"layers"`
	Port   int               `json:"port"`
}

var (
	config Config
	// port string
	// db_cache string
	config_file   string
	print_version bool
)

// Serve a single stylesheet via HTTP. Open view_tileserver.html in your browser
// to see the results.
// The created tiles are cached in an sqlite database (MBTiles 1.2 conform) so
// successive access a tile is much faster.
func TileserverWithCaching(layer_config map[string]string) {
	bind := fmt.Sprintf("0.0.0.0:%v", config.Port)
	t := maptiles.NewTileServerMux(config.Cache)

	for i := range layer_config {
		t.AddMapnikLayer(i, layer_config[i])
	}

	// maptiles.Ligneous.Info("Connecting to postgres database:")
	// maptiles.Ligneous.Info("*** ", config.Cache)
	maptiles.Ligneous.Info(fmt.Sprintf("Magic happens on port %v...", config.Port))
	srv := &http.Server{
		Addr:         bind,
		Handler:      t.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	maptiles.Ligneous.Error(srv.ListenAndServe())
}

func init() {
	// TODO: add config file
	flag.StringVar(&config_file, "c", "", "tile server config")
	flag.BoolVar(&print_version, "v", false, "version")
	flag.Parse()
	if print_version {
		fmt.Println(maptiles.SERVER_NAME + "-" + maptiles.VERSION)
		os.Exit(1)
	}
}

func getConfig() {
	// check if file exists!!!
	if _, err := os.Stat(config_file); err == nil {

		file, err := ioutil.ReadFile(config_file)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(file, &config)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

		maptiles.Ligneous.Debug(config)
	} else {
		fmt.Println("Config file not found")
		os.Exit(1)
	}
}

// Before uncommenting the GenerateOSMTiles call make sure you have
// the necessary OSM sources. Consult OSM wiki for details.
func main() {
	getConfig()
	TileserverWithCaching(config.Layers)
}
