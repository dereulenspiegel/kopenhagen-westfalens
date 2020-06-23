package main

import (
	"archive/zip"
	"encoding/xml"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	rootDocName = "doc.kml"
)

var (
	importPath = flag.String("in", "", "Import base path")
	outPath    = flag.String("out", "./content/post", "Where to generate page bundles")
)

type KML struct {
	Document *KMLDocument `xml:"Document"`
}

type KMLDocument struct {
	Name        string     `xml:"name"`
	Description string     `xml:"description"`
	Folder      *KMLFolder `xml:"Folder"`
}

type KMLFolder struct {
	Name       string       `xml:"name"`
	Placemarks []*Placemark `xml:"Placemark"`
}

type Placemark struct {
	Name         string        `xml:"name"`
	Description  string        `xml:"description"`
	Point        *Point        `xml:"Point"`
	ExtendedData *ExtendedData `xml:"ExtendedData"`
}

type Point struct {
	Coordinates string `xml:"coordinates"`
}

type ExtendedData struct {
	Data *Data `xml:"Data"`
}

type Data struct {
	Value string `xml:"value"`
}

func main() {
	flag.Parse()

	zipReader, err := zip.OpenReader(*importPath)
	if err != nil {
		log.Fatalf("Failed to open doc.kml (%s): %s", *importPath, err)
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		if f.Name == rootDocName {
			zf, err := f.Open()
			if err != nil {
				log.Fatalf("Failed to open file inside of kmz (%s): %w", f.Name, err)
			}
			defer zf.Close()
			kml := &KML{}
			if err := xml.NewDecoder(zf).Decode(kml); err != nil {
				log.Fatalf("Failed to decode doc.kml: %w", err)
			}
			generatePageBundles(kml.Document.Folder.Placemarks)
		}
	}

}

func generatePageBundles(placemarks []*Placemark) {
	for _, p := range placemarks {
		log.Printf("Found placemark: %s Description: %s", p.Name, p.Description)
		pageBundleFolderName := sanitizePlacemarkName(p.Name)
		pageBundlePath := filepath.Join(*outPath, pageBundleFolderName)
		if f, err := os.Stat(pageBundlePath); os.IsExist(err) && f.IsDir() {
			log.Printf("Page bundle for %s already exists, skipping", p.Name)
			continue
		}

		if err := os.Mkdir(pageBundlePath, 0766); err != nil {
			log.Fatalf("Failed to create page bundle directory %s: %s", pageBundlePath, err)
		}

		indexMd := generateFrontMatter(p)
		indexMd = indexMd + p.Description
		indexMdFile, err := os.Create(filepath.Join(pageBundlePath, "index.html"))
		if err != nil {
			log.Fatalf("Failed to create index html file %s: %s", filepath.Join(pageBundlePath, "index.html"), err)
		}
		defer indexMdFile.Close()
		_, err = indexMdFile.Write([]byte(indexMd))
		if err != nil {
			log.Fatalf("Failed to write index HTML file: %s", err)
		}
	}
}

func generateFrontMatter(placemark *Placemark) string {
	frontmatter := "---\n"
	frontmatter = frontmatter + "title: \"" + placemark.Name + "\"\n"
	frontmatter = frontmatter + "draft: false\n"
	frontmatter = frontmatter + "Date: " + time.Now().Format(time.RFC3339) + "\n"
	frontmatter = frontmatter + "infra_problem: \"problem\"\n"
	frontmatter = frontmatter + "source: mymaps\n"
	coords := strings.Split(strings.ReplaceAll(placemark.Point.Coordinates, "\n", ""), ",")
	if len(coords) < 2 {
		log.Fatalf("Coordinates (%s) have an invalid format", placemark.Point.Coordinates)
	}
	frontmatter = frontmatter + "location:\n"
	frontmatter = frontmatter + "  lat: " + coords[1] + "\n"
	frontmatter = frontmatter + "  lon: " + coords[0] + "\n"
	frontmatter = frontmatter + "---\n"
	return frontmatter
}

func sanitizePlacemarkName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, ".", "")
	name = strings.ReplaceAll(name, "/", "--")
	return name
}
