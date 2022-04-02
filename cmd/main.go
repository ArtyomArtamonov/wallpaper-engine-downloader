package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/ArtyomArtamonov/wallpaper-engine-downloader/internal/downloader"
	"github.com/ArtyomArtamonov/wallpaper-engine-downloader/internal/extractor"
	"github.com/BurntSushi/toml"
)

func main() {
	input := prepareInput()
	
	dwnl := downloader.NewSteamWorkshopDownloader()

	wallpaper := dwnl.Download(input.wallpaperWorkshopLink)

	extractor.Extract(input.myprojectPath, wallpaper)

	// Delete temporary directory and all files in it
	downloadDir :=  filepath.Join(filepath.Dir(wallpaper))
	log.Printf("Removing temp folder %s", downloadDir)
	err := os.RemoveAll(downloadDir)
	if err != nil {
		log.Fatal(err.Error())
	}
}

type Input struct {
	myprojectPath string
	wallpaperWorkshopLink string
}

type config struct {
	MyprojectPath string
}

func prepareInput() *Input {
	input := &Input{}
	// Workshop link from flag
	workshopLink := flag.String("wallpaper", "", "Link to wallpaper in steam workshop")
	
	wallpaperEngineProjectPath := flag.String("myproject", "", "where wallpapers will be unziped to")

	// Trying to find specified config.toml file
	configPath := flag.String("config", "", "Used to provide downloader with data such as wallpaper engine project folder etc.")
	flag.Parse()

	if *workshopLink == "" {
		log.Fatal("-wallpaper flag should not be empty")
	}

	if *configPath != "" {
		file, err := os.ReadFile(*configPath)
		if err != nil {
			log.Fatalf("error: could not find file with name %s", *configPath)
		}
		var config config
		err = toml.Unmarshal(file, &config)
		if err != nil {
			log.Fatal(err.Error())
		}

		input.myprojectPath = config.MyprojectPath
	} else {

		if *workshopLink == "" {
			log.Fatalf("'config' or 'wallpaper' have to be specified")
		}
		input.wallpaperWorkshopLink = *workshopLink

		if *wallpaperEngineProjectPath == "" {
			log.Fatalf("'config' or 'myproject' have to be specified")
		}
		input.myprojectPath = *wallpaperEngineProjectPath
	}

	input.wallpaperWorkshopLink = *workshopLink

	return input	
}

