package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type metar struct {
	Raw string
}

func main() {
	 metars, err := getMETARs()
	 if err != nil {
		 panic(fmt.Sprintf("error getting METARs: %q", err))
	 }

	 selected := selectMETAR(metars)
	 fmt.Println(selected.Raw)
}

func selectMETAR(metars []metar) metar{
	rand.Seed(time.Now().Unix())
	l := len(metars)
	s := rand.Intn(l)
	return metars[s]
}

func getMETARs() ([]metar, error) {
	resp, err := http.Get("https://www.aviationweather.gov/adds/dataserver_current/current/metars.cache.csv")
	if err != nil {
		return nil, fmt.Errorf("failed requesting metars cache")
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	err = readHeader(reader)
	if err != nil{
		return nil, err
	}

	r := csv.NewReader(reader)

	metars := []metar{}


	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return metars, fmt.Errorf("error parsing CSV: %w", err)
		}

		metar :=  metar{Raw: record[0]}
		metars = append(metars, metar)
	}

	return metars, nil

}

func readHeader(reader *bufio.Reader) error {
	_, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading header on 'errors' line")
	}
	//fmt.Printf("errors: %q\n", errors)

	_, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading headers on 'warnings' line")
	}
	//fmt.Printf("warnings: %q\n", warnings)

	_, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading headers on 'timing' line")
	}
	//fmt.Printf("timing: %q\n", timing)

	_, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading headers on 'dataSource' line")
	}
	//fmt.Printf("dataSource: %q\n", dataSource)

	_, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading headers on 'numResults' line")
	}
	//fmt.Printf("numResults: %q\n", numResults)

	_, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading headers on 'csvHeader' line")
	}
	//fmt.Printf("csvHeader: %q\n", csvHeader)

	return nil
}
