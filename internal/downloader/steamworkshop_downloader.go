package downloader

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const WALLPAPER_ENGINE_STEAM_APP_ID = "431960"

type SteamWorkshopDownloader struct {
	link string
}

func NewSteamWorkshopDownloader() *SteamWorkshopDownloader {
	return &SteamWorkshopDownloader{
		link: "http://steamworkshop.download/online/steamonline.php",
	}
}

func (s SteamWorkshopDownloader) Download(wallpaper string) string {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err.Error())
	}
	downloadFolder, err := ioutil.TempDir(filepath.Dir(execPath), "temp")

	if err != nil {
		log.Fatal(err.Error())
	}
	id, err := getIdFromWorkshopLink(wallpaper)
	if err != nil {
		log.Fatal(err.Error())
	}

	data := url.Values{
		"app": {WALLPAPER_ENGINE_STEAM_APP_ID},
		"item": {id},
	}
	resp, err := http.PostForm(s.link, data)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer resp.Body.Close()

	downloadLink := getDownloadLink(resp)
	archive, err := http.Get(downloadLink)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer archive.Body.Close()
	
	var result []byte
	archive.Body.Read(result)
	out, _ := os.Create(path.Join(downloadFolder, id + ".zip"))
	defer out.Close()
	io.Copy(out, archive.Body)

	return out.Name()
}

func getDownloadLink(resp *http.Response) string {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	link, exists := doc.Find("a").Attr("href")
	if !exists {
		log.Fatal("Could not found download link")
	}

	return link
}

func getIdFromWorkshopLink(link string) (string, error) {
	values, err := url.ParseQuery(link)
	if err != nil {
		return "", errors.New("could not fetch id from url")
	}

	for key, val := range values {
		if strings.Contains(key, "id") {
			return val[0], nil
		}
	} 

	return "", errors.New("string not found")
}
