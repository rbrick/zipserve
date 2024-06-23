package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	FlagDirectory = flag.String("dir", "", "select the directory to zip up.")
)

func init() {
	flag.Parse()
}

func main() {
	os.Mkdir("serveme", os.ModePerm)

	fileInfo, err := os.Stat(*FlagDirectory)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("compressing directories...")
	err = filepath.WalkDir(*FlagDirectory, func(path string, entry fs.DirEntry, _ error) error {

		_, file := filepath.Split(path)
		if file == "" {
			return nil
		}

		if entry.IsDir() {
			if filepath.Dir(path) == fileInfo.Name() {
				file, err := os.Create(fmt.Sprintf("serveme/%s.zip", entry.Name()))

				if err != nil {
					log.Fatalln(err)
				}

				writer := zip.NewWriter(file)
				rootPath := path
				err = filepath.WalkDir(rootPath, func(path string, entry fs.DirEntry, _ error) error {

					if entry.IsDir() {
						return nil
					}

					split := strings.Split(path, string(filepath.Separator))

					pathname := filepath.Join(split[1:]...)

					fileContent, err := os.ReadFile(path)
					if err != nil {
						return err
					}

					fileWriter, err := writer.Create(pathname)

					if err != nil {
						return err
					}

					_, err = fileWriter.Write(fileContent)

					if err != nil {
						log.Fatalln(err)
					}

					return nil
				})

				if err != nil {
					return err
				}

				writer.Close()
				return nil
			}

			return nil
		}

		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("serving files over HTTP")
	http.Handle("/", http.FileServer(http.Dir("serveme/")))

	http.ListenAndServe(":8080", nil)
}
