package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/cheggaaa/pb/v3"
	"gopkg.in/alecthomas/kingpin.v2"
)

type flags struct {
	url     *string
	workers *int
}

func main() {
	f := flags{
		url: kingpin.Arg("url", "URL of the thread").
			Required().
			String(),
		workers: kingpin.Flag("concurrent", "Number of concurrent downloaders").
			Short('c').
			Default("10").
			Int(),
	}
	_ = kingpin.Parse()

	u, err := NewThreadURLFromString(*f.url)
	if err != nil {
		log.Fatal(err)
	}

	th, err := threadInfo(u)
	if err != nil {
		log.Fatal(err)
	}

	path := u.id
	os.Mkdir(path, os.ModePerm)

	d, err := NewDownloader(path, 10)
	if err != nil {
		log.Fatal(err)
	}

	total := 0
	for _, post := range th.Posts {
		total += len(post.Files)
	}

	bar := pb.StartNew(total)
	i := 1
	for _, post := range th.Posts {
		for _, file := range post.Files {
			name := strconv.Itoa(i)
			ext := filepath.Ext(file.Path)
			d.Download(name+ext, "https://2ch.hk"+file.Path)
			i++
			bar.Increment()
		}
	}
	bar.Finish()
}
