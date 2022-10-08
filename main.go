package main

import (
	"fmt"
	"math"
	"math/rand"
)

const maxPeople = 4

func main() {
	metars, err := getMETARs()
	if err != nil {
		panic(fmt.Sprintf("error getting METARs: %q", err))
	}

	selected := selectMETAR(metars)

	fmt.Printf("Let's plan flight to the Airport: %q\n", selected.StationID)

	fmt.Println("METAR:")
	fmt.Println(selected.Raw)
	fmt.Println()

	tafs, err := getTAFs()
	if err != nil {
		panic(fmt.Sprintf("error getting TAFS: %q", err))
	}

	printTAF(tafs, selected.StationID)
	fmt.Println()

	printLoad()
}

func printLoad() {
	n := rand.Intn(maxPeople) + 1 // always at least one person
	fmt.Println("number of people on board:", n)

	for i := 0; i < n; i++ {
		fmt.Printf("person %d weight: %d lb\n", i+1, assumeWeight(200.0, 70.0))
	}

	fmt.Println()
	fmt.Println("number of bags on board:", n)

	for i := 0; i < n+2; i++ { // up to n+ 2 bags
		fmt.Printf("bag %d weight: %d lb\n", i+1, assumeWeight(10.0, 10.0))
	}
}

func assumeWeight(desiredMean, desiredStdDev float64) int {
	sample := rand.NormFloat64()*desiredStdDev + desiredMean
	cut := math.Min(math.Max(10.0, sample), 400.0) // [20, 400]
	return int(math.Ceil(cut))
}
