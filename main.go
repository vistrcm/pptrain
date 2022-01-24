package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const headerLines = 5

type metar struct {
	Raw       string
	StationID string
}

type taf struct {
	Raw       string
	StationID string
}

type tafMap map[string]taf

func main() {
	metars, err := getMETARs()
	if err != nil {
		panic(fmt.Sprintf("error getting METARs: %q", err))
	}

	selected := selectMETAR(metars)
	fmt.Println(selected.Raw)
	fmt.Println()

	tafs, err := getTAFs()
	if err != nil {
		panic(fmt.Sprintf("error getting TAFS: %q", err))
	}

	selectedTAF, presented := tafs[selected.StationID]
	if !presented {
		fmt.Printf("Can't find TAF for %q\n", selected.StationID)
		return
	}
	printTAF(selectedTAF.Raw)
}

func printTAF(raw string) {
	for _, token := range strings.Split(raw, " ") {
		if isChangeIndicator(token) {
			fmt.Println()
			fmt.Print("     ")  // changes indented
		}
		fmt.Print(token)
		fmt.Print(" ")
	}
	fmt.Println()
}

func isChangeIndicator(token string) bool {
	if strings.HasPrefix(token, "FM") {
		return true
	}
	if token == "TEMPO" || token == "BECMG" || token == "PROB" {
		return true
	}
	return false
}

func selectMETAR(metars []metar) metar {
	rand.Seed(time.Now().Unix())
	l := len(metars)
	for {
		s := rand.Intn(l)
		candidate := metars[s]
		if strings.HasPrefix(candidate.StationID, "K") {
			return candidate
		}
	}
}

func getMETARs() ([]metar, error) {
	resp, err := http.Get("https://www.aviationweather.gov/adds/dataserver_current/current/metars.cache.csv")
	if err != nil {
		return nil, fmt.Errorf("failed requesting metars cache")
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	err = readHeader(reader, headerLines)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(reader)
	if err := readCSVHeader(r); err != nil {
		return nil, err
	}

	metars := []metar{}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return metars, fmt.Errorf("error parsing CSV: %w", err)
		}

		metar := metar{Raw: record[0], StationID: record[1]}
		metars = append(metars, metar)
	}

	return metars, nil

}

func readCSVHeader(r *csv.Reader) error {
	_, err := r.Read()

	if err != nil {
		return fmt.Errorf("error parsing CSV header: %w", err)
	}

	return nil
}

func getTAFs() (tafMap, error) {
	resp, err := http.Get("https://www.aviationweather.gov/adds/dataserver_current/current/tafs.cache.csv")
	if err != nil {
		return nil, fmt.Errorf("failed requesting metars cache")
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	err = readHeader(reader, headerLines)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(reader)
	r.FieldsPerRecord = -1 // for some reason TAFs has variable amount of fields

	if err := readCSVHeader(r); err != nil {
		return nil, err
	}

	tafs := make(tafMap)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return tafs, fmt.Errorf("error parsing CSV: %w", err)
		}

		stationID := record[1]
		_, prst := tafs[stationID]
		if prst {
			return nil, fmt.Errorf("unexpected duplicate taf for %q", stationID)
		}

		taf := taf{Raw: record[0], StationID: stationID}
		tafs[stationID] = taf
	}

	return tafs, nil

}

func readHeader(reader *bufio.Reader, skip int) error {

	for i := 0; i <= skip-1; i++ {
		_, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading header line %d", i)
		}
	}

	return nil
}
