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
//TODO if value unchanged provide error
func listSimps(recs [][]string, symptom string) []string {
	x := len(recs[1]) - 1
	progs := make([]string, 0)
	fmt.Printf(symptom)
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
	return progs
}

func listSimpsMult(recs [][]string, symptoms []string) {
	progSet := make([]string, 0)
	finProgs := make(map[string]int)
	for _, symptom := range symptoms {
		tempProgs := listSimps(recs, symptom)
		progSet = append(progSet, tempProgs...)
	}
	for _, index := range progSet {
		finProgs[index] = finProgs[index] + 1
	}
	fmt.Println(finProgs)
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
	listSimpsMult(records, []string{"itching", "skin_rash"})
}
