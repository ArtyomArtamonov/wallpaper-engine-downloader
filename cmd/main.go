package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ArtyomArtamonov/wallpaper-engine-downloader/internal/downloader"
	"github.com/ArtyomArtamonov/wallpaper-engine-downloader/internal/extractor"
	"github.com/BurntSushi/toml"
)

var workshopLink = flag.String("wallpaper", "", "Link to wallpaper in steam workshop")
var wallpaperEngineProjectPath = flag.String("myproject", "", "where wallpapers will be unziped to")
var configPath = flag.String("config", "", "Used to provide downloader with data such as wallpaper engine project folder etc.")

func main() {
	input := prepareInput()

	for {
		if input.wallpaperWorkshopLink == "" {
			input.wallpaperWorkshopLink = askForWallpaperLinkFromInput()
		}

		dwnl := downloader.NewSteamWorkshopDownloader()

		wallpaper := dwnl.Download(input.wallpaperWorkshopLink)

		extractor.Extract(input.myproject, wallpaper)

		// Delete temporary directory and all files in it
		downloadDir := filepath.Join(filepath.Dir(wallpaper))
		log.Printf("Removing temp folder %s", downloadDir)
		err := os.RemoveAll(downloadDir)
		if err != nil {
			log.Fatal(err.Error())
		}

		input.wallpaperWorkshopLink = ""
	}
}

type Input struct {
	myproject         string
	wallpaperWorkshopLink string
}

type config struct {
	MyprojectPath string
}

func prepareInput() *Input {
	input := &Input{}

	flag.Parse()

	input.wallpaperWorkshopLink = *workshopLink
	input.myproject = *wallpaperEngineProjectPath

	var file []byte
	var err error

	if *configPath == "" {
		file, err = getDefaultConfig()
		if err != nil {
			log.Fatal(err.Error())
		}

	} else {
		file, err = os.ReadFile(*configPath)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	var config config
	err = toml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	input.myproject = config.MyprojectPath
	input.wallpaperWorkshopLink = *workshopLink

	if input.myproject == "" {
		log.Fatal("ERROR: Could not find myproject path." +
		"Please specify it in config as 'myproject' or pass as flag -myproject")
	}

	return input
}

func askForWallpaperLinkFromInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter wallpaper workshop link (or n/N to exit): ")
	text, _ := reader.ReadString('\n')
	text = strings.ReplaceAll(text, "\r\n", "")

	if strings.ToLower(text) == "n" {
		os.Exit(0)
	}

	return text
}

func getDefaultConfig() ([]byte, error) {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err.Error())
	}

	thisPath := filepath.Join(filepath.Dir(execPath))
	defaultConfigPath := filepath.Join(thisPath, "config.toml")

	return os.ReadFile(defaultConfigPath)
}
