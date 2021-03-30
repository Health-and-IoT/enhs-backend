package main

import (
	"container/heap"
	"container/list"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"

	"enhstools"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/fatih/color"
	"github.com/gorilla/mux"

	"github.com/rs/cors"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var emailData enhstools.EmailData
var records [][]string
var queue list.List

//Func update - used to update values before sending to firebase. Used mainly when creating new form and getting patients ID.
func update(v interface{}, updates map[string]string) {
	rv := reflect.ValueOf(v).Elem()
	for key, val := range updates {
		fv := rv.FieldByName(key)
		fv.SetString(val)
	}
}

//Func updateForm - used to update form when requested from web app.
func updateForm(w http.ResponseWriter, r *http.Request) {
	//Set headers for response
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//Get vars from request.
	params := mux.Vars(r)

	//Alert - id recieved.
	color.Yellow("ID Recieved ✔️")
	color.Yellow(params["id"])
	//Firebase setup
	ctx := context.Background()
	sa := option.WithCredentialsFile("sk.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	//Get post method.
	switch r.Method {
	case "GET":

	case "POST":
		//Read request
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {

			log.Fatal(err)
		}

		//fmt.Println(string(body))
		//Convert to Form
		var p enhstools.Form

		json.Unmarshal([]byte(body), &p)
		//Firebase update.
		_, err = client.Collection("form").Doc(p.DocID).Update(ctx, []firestore.Update{
			{
				Path:  "Approved",
				Value: p.Approved,
			},
			{
				Path:  "Symptoms",
				Value: p.Symptoms,
			},
			{
				Path:  "Pain",
				Value: p.Pain,
			},
			{
				Path:  "Priority",
				Value: p.Priority,
			},
			{
				Path:  "Seen",
				Value: p.Seen,
			},
		})
		if err != nil {
			// Handle any errors in an appropriate way, such as returning them.
			log.Printf("An error has occurred: %s", err)
		}

		//Alert
		color.Green("Form Updated - " + p.DocID)
		//Send response.
		w.Write([]byte(`{"success":true}`))

	default:

	}

}

//func retSymptoms JB Returns a JSON array of symptoms
func retSymptoms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	simps := enhstools.ListAllSimps(records)
	//simpsJSON, _ := json.Marshal(simps)
	w.Write(simps)

}

//Func deleteForm - deletes form when requested.
func test(w http.ResponseWriter, r *http.Request) {
	//Set headers for response
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//Get vars from request.

	switch r.Method {
	case "GET":
		color.Yellow("Test recieved ✔️")
		w.Write([]byte(`{"success":true}`))
	case "POST":
		//Read request
		color.Green("Test recieved ✔️")
		w.Write([]byte(`{"success":true}`))

	default:
		color.Red("Test recieved ✔️")
		w.Write([]byte(`{"success":true}`))
	}

}

//Func deleteForm - deletes form when requested.
func deleteForm(w http.ResponseWriter, r *http.Request) {
	//Set response headers
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//Get request vars
	params := mux.Vars(r)
	//Alert - delete request recieved.
	color.Yellow("Delete Request Recieved ✔️")

	//Firebase setup.
	ctx := context.Background()
	sa := option.WithCredentialsFile("sk.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	//Delete doc from firebase based on ID recieved in request.
	_, err2 := client.Collection("form").Doc(params["id"]).Delete(ctx)
	if err2 != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}
	//Alert - deleted
	color.Red("Deleted document - ")
	color.Red(params["id"])
	//Send response.
	w.Write([]byte(`{"success":true}`))

}

//Func getVisits - get visits when requested
func getVisits(w http.ResponseWriter, r *http.Request) {
	//Set response headers
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//Get request vars
	params := mux.Vars(r)

	//Alert - id recieved
	color.Yellow("ID Recieved ✔️")
	color.Yellow(params["id"])
	//Setup firebase
	ctx := context.Background()
	sa := option.WithCredentialsFile("sk.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	//Iterate through form from firebase where patient id matches.
	iter := client.Collection("form").Where("Patient", "==", params["id"]).Documents(ctx)
	var f []enhstools.Form
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {

			log.Fatalf("Failed to iterate: %v", err)
		}
		//Convert to Form
		var nyData enhstools.Form
		if err := doc.DataTo(&nyData); err != nil {
			// TODO: Handle error.
		}
		//Conversion
		jsonString, _ := json.Marshal(doc.Data())

		s := enhstools.Form{}
		//convert to Form struct
		json.Unmarshal(jsonString, &s)
		//fmt.Println(s.DateSubmitted)
		//fmt.Println(string(jsonString))

		// //fmt.Println(doc.Data())
		// //fmt.Println(nyData)
		//Append to array for visits
		f = append(f, s)

	}
	//Send visits as an array in JSON format.
	json.NewEncoder(w).Encode(f)

}

//GetPatient - gets patient by id.
func getPatient(w http.ResponseWriter, r *http.Request) {
	//Set response header
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//Get request vars
	params := mux.Vars(r)

	//Alert - id recieved.
	color.Yellow("ID Recieved ✔️")
	color.Yellow(params["id"])
	//Setup firebase
	ctx := context.Background()
	sa := option.WithCredentialsFile("sk.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	//Get patient based on id.
	dsnap, err2 := client.Collection("patient").Doc(params["id"]).Get(ctx)
	if err2 != nil {
		log.Fatal(err2)
	}

	//convert to patient struct
	var c enhstools.Patient
	dsnap.DataTo(&c)

	//Log and send response
	fmt.Printf("Document data: %#v\n", c)
	log.Println("Patient retrieved - ", params["id"])
	json.NewEncoder(w).Encode(c)
}

//Get Site - gets site based on ID.
func getSite(w http.ResponseWriter, r *http.Request) {
	//Set response headers
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//Get request vars
	params := mux.Vars(r)

	//Alert - id recieved
	color.Yellow("ID Recieved ✔️")
	color.Yellow(params["id"])
	//Setup firebase
	ctx := context.Background()
	sa := option.WithCredentialsFile("sk.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	//Iterate through firebase collection of sites until id matches.
	iter := client.Collection("sites").Where("siteid", "==", params["id"]).Documents(ctx)
	//Convert to site struct
	var nyData enhstools.Site
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {

			log.Fatalf("Failed to iterate: %v", err)
		}

		if err := doc.DataTo(&nyData); err != nil {
			// TODO: Handle error.
		}
		//Conversion
		jsonString, _ := json.Marshal(doc.Data())

		//convert to Form struct
		json.Unmarshal(jsonString, &nyData)
		//fmt.Println(s.DateSubmitted)
		//fmt.Println(string(jsonString))

		// //fmt.Println(doc.Data())
		// //fmt.Println(nyData)

	}
	//Log and respond
	fmt.Printf("Document data: %#v\n", nyData)
	json.NewEncoder(w).Encode(nyData)

}

//Func getUser - gets user account based on id.

func getUser(w http.ResponseWriter, r *http.Request) {
	//Set response headers
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//Get request vars
	params := mux.Vars(r)

	//Alert - id recieved.
	color.Yellow("ID Recieved ✔️")
	color.Yellow(params["id"])
	//Setup firebase
	ctx := context.Background()
	sa := option.WithCredentialsFile("sk.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	//Find user based on id
	dsnap, err2 := client.Collection("users").Doc(params["id"]).Get(ctx)
	if err2 != nil {
		log.Fatal(err2)
	}

	//Make into logged user
	var c enhstools.LoggedInUser
	dsnap.DataTo(&c)

	//Record and respond
	fmt.Printf("Document data: %#v\n", c)
	log.Println("User retrieved - ", params["id"])
	json.NewEncoder(w).Encode(c)
}

//Func getPatients - get all forms for patients
func getPatients(w http.ResponseWriter, r *http.Request) {
	//Set response headers
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//Get request vars
	params := mux.Vars(r)

	//Alert - id recieved
	color.Yellow("ID Recieved 1 ✔️")
	color.Yellow(params["id"])
	//Setup Firebase
	ctx := context.Background()
	sa := option.WithCredentialsFile("sk.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	//Iterate through forms.
	iter := client.Collection("form").Where("SiteID", "==", params["id"]).Documents(ctx)
	var f []enhstools.Form
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {

			log.Fatalf("Failed to iterate: %v", err)
		}
		//Convert to form
		var nyData enhstools.Form
		if err := doc.DataTo(&nyData); err != nil {
			// TODO: Handle error.
		}
		jsonString, _ := json.Marshal(doc.Data())

		s := enhstools.Form{}
		//convert to Form struct
		json.Unmarshal(jsonString, &s)
		//fmt.Println(s.DateSubmitted)
		//fmt.Println(string(jsonString))

		// //fmt.Println(doc.Data())
		// //fmt.Println(nyData)
		//Append to form array
		f = append(f, s)

	}
	//Respond
	json.NewEncoder(w).Encode(f)

}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	//Set response headers
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//Get request vars
	params := mux.Vars(r)

	//Alert - id recieved
	color.Yellow("ID Recieved 1 ✔️")
	color.Yellow(params["id"])
	//Setup Firebase
	ctx := context.Background()
	sa := option.WithCredentialsFile("sk.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	//Iterate through forms.
	iter := client.Collection("events").Documents(ctx)
	var f []enhstools.Event
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {

			log.Fatalf("Failed to iterate: %v", err)
		}
		//Convert to form
		var nyData enhstools.Event
		if err := doc.DataTo(&nyData); err != nil {
			// TODO: Handle error.
		}
		jsonString, _ := json.Marshal(doc.Data())

		s := enhstools.Event{}
		//convert to Form struct
		json.Unmarshal(jsonString, &s)
		//fmt.Println(s.DateSubmitted)
		//fmt.Println(string(jsonString))

		// //fmt.Println(doc.Data())
		// //fmt.Println(nyData)
		//Append to form array
		f = append(f, s)

	}
	//Respond
	json.NewEncoder(w).Encode(f)

}

//Func authLogin - used to authorise user and log them into application.
func authLogin(w http.ResponseWriter, r *http.Request) {
	//Set response headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Type")
	//Setup firebase
	ctx := context.Background()
	sa := option.WithCredentialsFile("sk.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	switch r.Method {
	case "GET":

	case "POST":
		//Read request

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {

			log.Fatal(err)
		}

		//fmt.Println(string(body))

		//Setup login var
		var p enhstools.Login
		//Convert to struct
		json.Unmarshal([]byte(body), &p)
		color.Yellow("Login Attempt Recieved from: " + p.Username)
		log.Println("Login Attempt Recieved from: " + p.Username)
		//Iterate through collection to see if details match
		iter := client.Collection("users").Where("username", "==", p.Username).Where("password", "==", p.Password).Where("siteid", "==", p.Siteid).Documents(ctx)

		for {
			doc, err2 := iter.Next()

			if err2 == iterator.Done {

				break
			}

			//Match - respond and log
			color.Green("Login Successful from: " + p.Username)
			log.Println("Login Successful from: " + p.Username)
			w.Write([]byte(`{"success":true, "id":"` + doc.Ref.ID + `"}`))

		}

	default:

	}

}

//Work in progress for real-time updates.
func checkForNewForm() {
	ctx := context.Background()
	sa := option.WithCredentialsFile("sk.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// currentTime := time.Now().Add(-10 * time.Minute)
	// snapIter := client.Collection("form").OrderBy("DateSubmitted", firestore.Asc).Snapshots(ctx)

	// defer snapIter.Stop()

	// for {
	// 	snap, err := snapIter.Next()
	// 	if err != nil {
	// 		color.Yellow("", err)
	// 	}

	// 	for _, diff := range snap.Changes {

	// 		jsonString, _ := json.Marshal(diff.Doc.Data())

	// 		s := Form{}
	// 		//convert to Form struct
	// 		json.Unmarshal(jsonString, &s)

	// 		if s.DateSubmitted > currentTime.Format("02/01/2006, 15:04:05") {

	// 			color.Cyan("New Form Received 🔔")
	// 			color.Cyan("Date and time submitted: %v", s.DateSubmitted)
	// 		} else {

	// 		}

	// 	}

	// }
	cols, err := client.Collections(context.Background()).GetAll()

	for _, col := range cols {
		iter := col.Snapshots(context.Background())
		defer iter.Stop()

		for {
			doc, err := iter.Next()
			if err != nil {
				if err == iterator.Done {
					break
				}

			}

			for _, change := range doc.Changes {
				// access the change.Doc returns the Document,
				// which contains Data() and DataTo(&p) methods.
				switch change.Kind {
				case firestore.DocumentAdded:
					// on added it returns the existing ones.
					//isNew := change.Doc.CreateTime.After(l.startTime)
					// [...]
				case firestore.DocumentModified:
					// [...]
				case firestore.DocumentRemoved:
					// [...]
				}
			}
		}
	}
}

type IntHeap []int

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

//func Main - this is where it all starts
func main() {
	//QRCode
	//err := qrcode.WriteFile("hi", qrcode.Medium, 256, "qr.png")
	//Inital setup
	//Start new queue
	nums := []int{3, 2, 20, 5, 3, 1, 2, 5, 6, 9, 10, 4}

	// initialize the heap data structure
	h := &IntHeap{}

	// add all the values to heap, O(n log n)
	for _, val := range nums { // O(n)
		heap.Push(h, val) // O(log n)
	}

	// print all the values from the heap
	// which should be in ascending order
	//for i := 0; i < len(nums); i++ {
	//	fmt.Printf("%d,", heap.Pop(h).(int))
	//}
	queue := list.New()
	color.Green("Backend server started! ✔️")
	//Check for logfile - if none, create one.
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("Backend server started!")

	//Setup firebase
	ctx := context.Background()
	sa := option.WithCredentialsFile("sk.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	//JAKE BLOCK ADDED 18/02/21
	csvfile, err := os.Open("data/testing.csv")
	if err != nil {
		log.Fatalln(err)
	}
	reader := csv.NewReader(csvfile)
	records, err = reader.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}
	defer csvfile.Close()
	log.Println("CSV Data Loaded.")
	jsonFile, err := ioutil.ReadFile("email.json")
	if err != nil {
		fmt.Println(err)
	}
	_ = json.Unmarshal([]byte(jsonFile), &emailData)
	log.Println("Email Config Loaded.")

	r := mux.NewRouter()
	r.HandleFunc("/getPatient/{id}", getPatient).Methods("POST")
	r.HandleFunc("/getVisits/{id}", getVisits).Methods("POST")
	r.HandleFunc("/getUser/{id}", getUser).Methods("POST")
	r.HandleFunc("/getPatients/{id}", getPatients).Methods("POST")
	r.HandleFunc("/getAllEvents", getAllEvents).Methods("POST")
	r.HandleFunc("/updateForm/{id}", updateForm).Methods("POST")
	r.HandleFunc("/deleteForm/{id}", deleteForm).Methods("POST")
	r.HandleFunc("/getSite/{id}", getSite).Methods("POST")
	r.HandleFunc("/login/", authLogin).Methods("POST")

	r.HandleFunc("/symptoms", retSymptoms).Methods("GET")
	r.HandleFunc("/test", test).Methods("GET")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case "GET":

		case "POST":
			color.Yellow("New Form Recieved ✔️")

			body, err := ioutil.ReadAll(r.Body)

			if err != nil {

				log.Fatal(err)
			}

			//fmt.Println(string(body))

			var p enhstools.Request

			//NOTE & denotes value passed by reference ie value is changed by function.
			json.Unmarshal([]byte(body), &p)
			pat := p.Patient
			form := p.Form
			ref := client.Collection("patient").NewDoc()

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write([]byte(`{"success":true}`))
			_, err1 := ref.Set(ctx, pat)
			if err != nil {
				// Handle any errors in an appropriate way, such as returning them.
				log.Printf("An error has occurred: %s", err1)
			}

			ref2 := client.Collection("form").NewDoc()
			updates := map[string]string{"Patient": ref.ID}

			update(&form, updates)
			updates1 := map[string]string{"DocID": ref2.ID}

			update(&form, updates1)

			str := enhstools.ListSimpsMult(records, form.Symptoms)
			form.ProgList = string(str)

			//Mail method with to be added variables (Domain, mailAPIKey and Sender)
			//form.Email to be used instead of hardcoded email.

			enhstools.Mail(emailData.Domain, emailData.APIKey, "40315515@live.napier.ac.uk", emailData.Sender, form.SiteID, form.DocID)
			log.Println("Email Sent for: ", form.DocID)
			fmt.Println(form)
			_, err2 := ref2.Set(ctx, form)
			if err != nil {
				// Handle any errors in an appropriate way, such as returning them.
				log.Printf("An error has occurred: %s", err2)
			}
			color.Yellow("Patient and Form added to DB ✔️")
			log.Println("Patient and Form added to DB : ", ref2.ID)
			queue.PushBack(ref2.ID)

			color.Green("Queue")
			// Dequeue
			//front := queue.Front()
			//fmt.Println(front.Value)
			// This will remove the allocated memory and avoid memory leaks
			//queue.Remove(front)
			for e := queue.Front(); e != nil; e = e.Next() {

				fmt.Println("Queue:", e.Value.(string))
			}
		default:

		}
	})

	//Starts and opens port allowing connections on runtime.
	handler := cors.Default().Handler(r)

	if err := http.ListenAndServeTLS(":8080", "/etc/apache2/certificate/apache-certificate.crt", "/etc/apache2/certificate/apache.key", handler); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
