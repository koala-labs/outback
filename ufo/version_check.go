package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const VERSION_FILE = "https://s3.amazonaws.com/fuzz-ufo/current_version.json"

type Version struct {
	Version string `json:"version"`
}

func AssertCurrentVersion(ufoVersion string) error {
	res, err := http.Get(VERSION_FILE)

	if err != nil {
		return ErrCouldNotAssertVersion
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return ErrCouldNotAssertVersion
	}

	v := &Version{}
	err = json.Unmarshal(body, v)

	if err != nil {
		return ErrCouldNotAssertVersion
	}

	if v.Version != ufoVersion {
		return ErrUFOOutOfDate
	}

	return nil
}