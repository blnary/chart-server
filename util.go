package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func get(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, BPM_SERVER_URL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating GET request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending GET request: %w", err)
	}
	return resp, nil
}

func getBPM(song *Song) (*InternalGetBPMReply, error) {
	resp, err := get("bpm/?filename=" + song.Location)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var reply InternalGetBPMReply
	err = json.NewDecoder(resp.Body).Decode(&reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

func relErr(a float64, b float64, unit float64) float64 {
	delta := a - b
	if delta < 0 {
		delta = -delta
	}
	for delta > unit {
		delta -= unit
	}
	if delta > unit/2 {
		delta = unit - delta
	}
	return delta
}

func eqErr(a float64, b float64) float64 {
	if a < b {
		t := a
		a = b
		b = t
	}
	for a > b {
		a -= b
	}
	if a > b/2 {
		a = b - a
	}
	return a
}
