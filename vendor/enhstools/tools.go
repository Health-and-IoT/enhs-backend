package enhstools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	mailgun "github.com/mailgun/mailgun-go"
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

// ListProgs Returns list of prognoses
// Usage:
//		 var records [][]string
//     func main() {
//     		 var prognoses []string
//	 			 --->
// 		 		 population of records from csv file
//	 		 	 --->
//         prognoses = enhstools.ListProgs(records)
// 				 fmt.Println(prognoses)
//     }
//
// Output: array[Prognosis1 Prognosis2 ... ]
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
// Usage:
//		 var records [][]string
//     func main() {
// 		 		 var allSymptoms []byte
//	 			 --->
// 		 		 population of records from csv file
//	 		 	 --->
//         allSymptoms = enhstools.ListAllSimps(records)
// 				 fmt.Println(string(allSymptoms))
//     }
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

//listSimps Returns all prognoses related to a symptom if symptom exists
// Private Method
// Usage:
//		 var records [][]string
//     func main() {
// 		 		 var symptom string
//	 			 --->
// 		 		 population of records from csv file
// 			 	 value given to symptom
//	 		 	 --->
//     		 var prognoses []string
//         prognoses = enhstools.listSimps(records, symptom)
// 				 fmt.Println(prognoses)
//     }
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
// Usage:
//		 var records [][]string
//     func main() {
// 				 var symptoms []string
//	 			 --->
// 		 		 population of records from csv file
// 			 	 symptoms given set of values
//	 		 	 --->
//     		 var prognoses []string
//         prognoses = enhstools.listSimpsMult(records, symptoms)
// 				 fmt.Println(string(prognoses))
//     }
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
	log.Println("ListSimpsMult Func called. Returned value: ", string(finProgsJSON))
	return finProgsJSON
}

// Mail - Sends Email Receipt of form
// Cut from GRAEME HILL's test_mail.go file.
func Mail(domain string, mailAPIKey, recipient string, sender string, locID string, formID string) {
	mg := mailgun.NewMailgun(domain, mailAPIKey)
	subject := "Receipt: Form Received"
	body := ""

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)
	message.SetTemplate("newmessage-2021-02-21.181649")
	message.AddTemplateVariable("location_id", locID)
	message.AddTemplateVariable("form_id", formID)
	message.AddTemplateVariable("sub_time", string(time.Now()))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Email Sent. ID: %s Resp: %s\n", id, resp)
}
