package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

var (
	threadURLTemplate = regexp.MustCompile(`^(https:\/\/|http:\/\/)?2ch\.hk\/([A-Za-z0-9]+)\/res\/([0-9]+)(\.json|\.html)?$`)
)

type threadURL struct {
	board string
	id    string
}

func NewThreadURLFromString(s string) (threadURL, error) {
	submatchs := threadURLTemplate.FindStringSubmatch(s)
	if len(submatchs) == 0 {
		return threadURL{}, fmt.Errorf("wrong url")
	}
	u := threadURL{
		board: submatchs[2],
		id:    submatchs[3],
	}
	return u, nil
}

func (t *threadURL) JSON() string {
	const pattern = `https://2ch.hk/%s/res/%s.json`
	return fmt.Sprintf(pattern, t.board, t.id)
}

type thread struct {
	Posts []post `json:"posts"`
}

type post struct {
	Files []file `json:"files"`
}

type file struct {
	Path string `json:"path"`
}

func threadInfo(u threadURL) (thread, error) {
	var data struct {
		Threads []thread `json:"threads"`
	}
	resp, err := http.Get(u.JSON())
	if err != nil {
		return thread{}, err
	}
	defer resp.Body.Close()
	d := json.NewDecoder(resp.Body)
	err = d.Decode(&data)
	if err != nil {
		return thread{}, err
	}
	if len(data.Threads) != 1 {
		return thread{}, fmt.Errorf("can't found thread")
	}
	return data.Threads[0], nil
}
