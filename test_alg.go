package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

//prints list of prognosises
func listProgs(recs [][]string) {
	x := len(recs[1]) - 1
	for row := range recs {
		//records[row][column]
		fmt.Printf(recs[row][x] + ", ")
	}
	fmt.Println()
}

//lists prognosises for a select symptom
func listSimps(recs [][]string, symptom string) {
	x := len(recs[1]) - 1
	progs := make([]string, 0)
	fmt.Println(symptom)
	row := 0
	//TODO SORT THIS SO IT DOESNT RELATE TO ITCHING / ACTUAL VALUE
	sympIndex := 0
	for row < x {
		if symptom == recs[0][row] {
			sympIndex = row
			//fmt.Println(symptom, sympIndex)
			break
		}
		row++
	}
	for row := range recs {
		if recs[row][sympIndex] == "1" {
			progs = append(progs, recs[row][x])
		}
	}
	fmt.Println(progs)
}

func main() {
	csvfile, err := os.Open("data/testing.csv")
	if err != nil {
		log.Fatalln(err)
	}
	r := csv.NewReader(csvfile)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}
	//listProgs(records)
	listSimps(records, "itching")
}
