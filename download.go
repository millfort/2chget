package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type downloader struct {
	workers  chan chan paramaters
	workload chan paramaters
	path     string

	client *http.Client
}

type paramaters struct {
	filename string
	path     string
}

func NewDownloader(path string, n int) (*downloader, error) {
	d := &downloader{
		workers:  make(chan chan paramaters, n),
		workload: make(chan paramaters),
		path:     path,
		client:   http.DefaultClient,
	}
	go d.start(n)
	return d, nil
}

func (d *downloader) start(n int) {
	for i := 0; i < n; i++ {
		ch := make(chan paramaters)
		go func(c chan paramaters) {
			for u := range c {
				d.download(u.filename, u.path)
				d.workers <- c
			}
		}(ch)
		d.workers <- ch
	}
	for u := range d.workload {
		for ch := range d.workers {
			ch <- u
			break
		}
	}
}

func (d *downloader) download(filename, u string) {
	path := filepath.Join(d.path, filename)
	DownloadFile(path, u)
}

func (d *downloader) Download(filename, path string) {
	d.workload <- paramaters{
		filename: filename,
		path:     path,
	}
}

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}
