package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
)

// Prints list of prognosises
func listProgs(recs [][]string) []string {
	allProgs := make([]string, 0)
	x := len(recs[1]) - 1
	for row := range recs {
		//records[row][column]
		//fmt.Printf(recs[row][x] + ", ")
		allProgs = append(allProgs, recs[row][x])
	}
	//fmt.Println()
	return allProgs
}

// Returns a list of all Symptoms
func listAllSimps(recs [][]string) []string {
	allSimps := make([]string, 0)
	for col := range recs[0] {
		//records[row][column]
		//fmt.Printf(recs[0][col] + ", ")
		allSimps = append(allSimps, recs[0][col])
	}
	//fmt.Println()
	return allSimps
}

//lists prognosises for a select symptom
func listSimps(recs [][]string, symptom string) []string {
	x := len(recs[1]) - 1
	progs := make([]string, 0)
	//fmt.Printf(symptom)
	row := 0
	sympIndex := -1
	for row < x {
		if symptom == recs[0][row] {
			sympIndex = row
			//fmt.Println(symptom, sympIndex)
			break
		}
		row++
	}
	if sympIndex == -1 {
		return progs
	}
	for row := range recs {
		if recs[row][sympIndex] == "1" {
			progs = append(progs, recs[row][x])
		}
	}
	//fmt.Println(progs)
	return progs
}

// Returns a sorted map of prognoses
func listSimpsMult(recs [][]string, symptoms []string) map[string]int {
	progSet := make([]string, 0)
	finProgs := make(map[string]int)
	for _, symptom := range symptoms {
		tempProgs := listSimps(recs, symptom)
		progSet = append(progSet, tempProgs...)
	}
	for _, index := range progSet {
		finProgs[index] = finProgs[index] + 1
	}
	//fmt.Println(finProgs)
	keys := make([]string, 0, len(finProgs))
	for k := range finProgs {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return finProgs[keys[i]] > finProgs[keys[j]]
	})
	for _, key := range keys {
		fmt.Printf("%-7v %v\n", key, finProgs[key])
	}
	return finProgs
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
	str := listProgs(records)
	//str := listAllSimps(records)
	//str := listSimpsMult(records, []string{"itching", "skin_rash"})
	fmt.Println(str)
}
