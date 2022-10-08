package main

import (
	"fmt"
)

func main() {
	metars, err := getMETARs()
	if err != nil {
		panic(fmt.Sprintf("error getting METARs: %q", err))
	}

	selected := selectMETAR(metars)

	fmt.Printf("Let's plan flight to the Airport: %q\n", selected.StationID)

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
