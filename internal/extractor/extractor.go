package extractor

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Extract(destination string, archiveName string) {
	archive, err := zip.OpenReader(archiveName)
    if err != nil {
        panic(err)
    }
    defer archive.Close()

    for _, f := range archive.File {
        filePath := filepath.Join(destination, f.Name)
        fmt.Println("unzipping file ", filePath)

        if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
            fmt.Println("invalid file path")
            return
        }
        if f.FileInfo().IsDir() {
            fmt.Println("creating directory...")
            os.MkdirAll(filePath, os.ModePerm)
            continue
        }

        if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
            panic(err)
        }

        dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
        if err != nil {
            panic(err)
        }
        defer dstFile.Close()

        fileInArchive, err := f.Open()
        if err != nil {
            panic(err)
        }
        defer fileInArchive.Close()

        if _, err := io.Copy(dstFile, fileInArchive); err != nil {
            panic(err)
        }
    }
}
