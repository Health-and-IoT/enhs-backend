package enhstools

import (
	"encoding/json"
	"log"
	"sort"
)

// Symptom - Symptom struct with fields ID and Name
type Symptom struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Prognosis - Prognosis struct with fields ID, Name and SympCount. SympCount indicates the number of symptoms the patient has that are indicative of each prognosis.
type Prognosis struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	SympMatch int    `json:"sympCount"`
}

// ListProgs Prints list of prognosises
func ListProgs(recs [][]string) []string {
	allProgs := make([]string, 0)
	x := len(recs[1]) - 1
	for row := range recs {

		allProgs = append(allProgs, recs[row][x])
	}
	log.Println("ListProgs Func called. Returned value: ", allProgs)
	return allProgs
}

// ListAllSimps Returns a list of all Symptoms
func ListAllSimps(recs [][]string) []byte {
	allSimps := make([]Symptom, 0)
	var simp Symptom
	for col := range recs[0] {
		simp.ID = col
		simp.Name = recs[0][col]
		allSimps = append(allSimps, simp)
	}
	allSimps = allSimps[:len(allSimps)-1]
	allSimpsJSON, _ := json.Marshal(allSimps)
	log.Println("ListAllSimps Func called. Returned value: ", string(allSimpsJSON))
	return allSimpsJSON
}

//ListSimps Returns prognosises for a select symptom
func listSimps(recs [][]string, symptom string) []string {
	x := len(recs[1]) - 1
	progs := make([]string, 0)
	row := 0
	sympIndex := -1
	for row < x {
		if symptom == recs[0][row] {
			sympIndex = row
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
	log.Println("ListSimps Func called. Returned: ", progs)
	return progs
}

// ListSimpsMult Returns a sorted map of prognoses
func ListSimpsMult(recs [][]string, symptoms []string) []byte {
	progSet := make([]string, 0)
	finProgsMap := make(map[string]int)
	finProgsProg := make([]Prognosis, 0)
	var progObj Prognosis
	for _, symptom := range symptoms {
		tempProgs := listSimps(recs, symptom)
		progSet = append(progSet, tempProgs...)
	}
	for _, index := range progSet {
		finProgsMap[index] = finProgsMap[index] + 1
	}
	keys := make([]string, 0, len(finProgsMap))
	for k := range finProgsMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return finProgsMap[keys[i]] > finProgsMap[keys[j]]
	})
	index := 0
	for prog, count := range finProgsMap {
		progObj.ID = index
		progObj.Name = prog
		progObj.SympMatch = count
		finProgsProg = append(finProgsProg, progObj)
		index++
	}
	finProgsJSON, _ := json.Marshal(finProgsProg)
	return finProgsJSON
}
