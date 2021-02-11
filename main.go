package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

//Form Struct - for v1 - not used currently.
type Form struct {
	Address       string `json:"address"`
	Name          string `json:"name"`
	Dob           string `json:"dob"`
	Nok           string `json:"nok"`
	Chinumber     string `json:"chinumber"`
	Allergies     string `json:"allergies"`
	DateSubmitted string `json:"dateSubmitted"`
	Id            string `json:"id"`
	Illness       string `json:"illness"`
	Pain          string `json:"pain"`
	Priority      string `json:"priority"`
	Seen          bool   `json:"seen"`
}

//Patient struct - used to get patient JSON strings when sent to server. Converts to GO strings.
type Patient struct {
	Address   string `json:"address"`
	Allergies string `json:"allergies"`
	Chinumber string `json:"chinumber"`
	Dob       string `json:"dob"`
	Donor     bool   `json:"donor"`
	Name      string `json:"name"`
	Nok       string `json:"nok"`
}

//Request struct - used to split patient and form when sent to server.
type Request struct {
	Patient Patient
	Form    Form1
}

// Login struct - used when recieving a login request from app.
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Siteid   string `json:"siteid"`
}

//Site struct - used to get site JSON strings when sent to server. Converts to GO Strings.
type Site struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Siteid  string `json:"siteid"`
}

//LoggedInUser struct - used to get non senstive JSON strings from firebase.
type LoggedInUser struct {
	Username string `json:"username"`
	Rank     string `json:"rank"`
}

//Form1 struct - used to get form JSON strings when sent to server. Converts to GO strings.
type Form1 struct {
	Ailment       string `json:"ailment"`
	DateSubmitted string `json:"dateSubmitted"`
	Pain          int64  `json:"pain"`
	Patient       string `json:"patient"`
	Priority      string `json:"priority"`
	Seen          bool   `json:"seen"`
	Approved      bool   `json:"approved"`
	DocID         string `json:"docID"`
}

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
		//Convert to Form1
		var p Form1

		json.Unmarshal([]byte(body), &p)
		//Firebase update.
		_, err = client.Collection("form").Doc(p.DocID).Update(ctx, []firestore.Update{
			{
				Path:  "Approved",
				Value: p.Approved,
			},
			{
				Path:  "Ailment",
				Value: p.Ailment,
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

//Func deleteForm - deletes form when requested.
func test(w http.ResponseWriter, r *http.Request) {
	//Set headers for response
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	//Get vars from request.

	switch r.Method {
	case "GET":
		w.Write([]byte(`{"success":true}`))
	case "POST":
		//Read request

		w.Write([]byte(`{"success":true}`))

	default:
		w.Write([]byte(`{"success":true}`))
	}

}

//Func deleteForm - deletes form when requested.
func deleteForm(w http.ResponseWriter, r *http.Request) {
	//Set response headers
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	var f []Form1
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {

			log.Fatalf("Failed to iterate: %v", err)
		}
		//Convert to Form1
		var nyData Form1
		if err := doc.DataTo(&nyData); err != nil {
			// TODO: Handle error.
		}
		//Conversion
		jsonString, _ := json.Marshal(doc.Data())

		s := Form1{}
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
	var c Patient
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
	var nyData Site
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
	var c LoggedInUser
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
	//Get request vars
	params := mux.Vars(r)

	//Alert - id recieved
	color.Yellow("ID Recieved ✔️")
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
	iter := client.Collection("form").Documents(ctx)
	var f []Form1
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {

			log.Fatalf("Failed to iterate: %v", err)
		}
		//Convert to form
		var nyData Form1
		if err := doc.DataTo(&nyData); err != nil {
			// TODO: Handle error.
		}
		jsonString, _ := json.Marshal(doc.Data())

		s := Form1{}
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
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
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
		var p Login
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

//func Main - this is where it all starts
func main() {
	//Inital setup
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

	//color.Green("Retrieving Firebase collection...⏳")
	//log.Println("Retrieving Firebase collection...")
	//iter := client.Collection("forms").Documents(ctx)
	//color.Green("Retrieved Firebase collection...✔️")
	//log.Println("Retrieved Firebase collection...")
	//Firebase iterator
	// for {
	// 	doc, err := iter.Next()
	// 	if err == iterator.Done {
	// 		break
	// 	}
	// 	if err != nil {

	// 		log.Fatalf("Failed to iterate: %v", err)
	// 	}
	// 	var nyData Form
	// 	if err := doc.DataTo(&nyData); err != nil {
	// 		// TODO: Handle error.
	// 	}
	// 	// jsonString, _ := json.Marshal(doc.Data())

	// 	// s := Form{}
	// 	// //convert to Form struct
	// 	// json.Unmarshal(jsonString, &s)
	// 	// //fmt.Println(s.DateSubmitted)
	// 	// //fmt.Println(string(jsonString))

	// 	// //fmt.Println(doc.Data())
	// 	// //fmt.Println(nyData)
	// }
	//color.Green("Forms initialised ✔️")

	//Routers - the api aspect of the application.
	r := mux.NewRouter()
	r.HandleFunc("/getPatient/{id}", getPatient).Methods("POST")
	r.HandleFunc("/getVisits/{id}", getVisits).Methods("POST")
	r.HandleFunc("/getUser/{id}", getUser).Methods("POST")
	r.HandleFunc("/getPatients", getPatients).Methods("POST")
	r.HandleFunc("/updateForm/{id}", updateForm).Methods("POST")
	r.HandleFunc("/deleteForm/{id}", deleteForm).Methods("POST")
	r.HandleFunc("/getSite/{id}", getSite).Methods("POST")
	r.HandleFunc("/login/", authLogin).Methods("POST")
	r.HandleFunc("/test", test).Methods("POST")
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

			var p Request

			//NOTE & denotes value passed by reference ie value is changed by function.
			json.Unmarshal([]byte(body), &p)
			pat := p.Patient
			form := p.Form
			ref := client.Collection("patient").NewDoc()

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write([]byte("{\"patient\": \"}"))
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
			fmt.Println(form)
			_, err2 := ref2.Set(ctx, form)
			if err != nil {
				// Handle any errors in an appropriate way, such as returning them.
				log.Printf("An error has occurred: %s", err2)
			}
			color.Yellow("Patient and Form added to DB ✔️")
			log.Println("Patient and Form added to DB : ", ref.ID)
		default:

		}
	})

	//Starts and opens port allowing connections on runtime.
	//NOTE JB - I think that this port should be 443
	handler := cors.Default().Handler(r)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
