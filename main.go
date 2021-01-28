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
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

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

type Patient struct {
	Address   string `json:"address"`
	Allergies string `json:"allergies"`
	Chinumber string `json:"chinumber"`
	Dob       string `json:"dob"`
	Donor     bool   `json:"donor"`
	Name      string `json:"name"`
	Nok       string `json:"nok"`
}
type Request struct {
	Patient Patient
	Form    Form1
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Siteid   string `json:"siteid"`
}
type Site struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Siteid  string `json:"siteid"`
}
type LoggedInUser struct {
	Username string `json:"username"`
	Rank     string `json:"rank"`
}
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

func update(v interface{}, updates map[string]string) {
	rv := reflect.ValueOf(v).Elem()
	for key, val := range updates {
		fv := rv.FieldByName(key)
		fv.SetString(val)
	}
}

func updateVisit1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)

	color.Yellow("ID Recieved ✔️")
	color.Yellow(params["id"])
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

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {

			log.Fatal(err)
		}

		//fmt.Println(string(body))

		var p Form1

		json.Unmarshal([]byte(body), &p)

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

		color.Green("Form Updated - " + p.DocID)

		w.Write([]byte(`{"success":true}`))

	default:

	}

}
func getVisits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)

	color.Yellow("ID Recieved ✔️")
	color.Yellow(params["id"])
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
		f = append(f, s)

	}
	json.NewEncoder(w).Encode(f)

}
func getPatient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)

	color.Yellow("ID Recieved ✔️")
	color.Yellow(params["id"])
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

	dsnap, err2 := client.Collection("patient").Doc(params["id"]).Get(ctx)
	if err2 != nil {
		log.Fatal(err2)
	}

	var c Patient
	dsnap.DataTo(&c)

	fmt.Printf("Document data: %#v\n", c)
	log.Println("Patient retrieved - ", params["id"])
	json.NewEncoder(w).Encode(c)
}

func getSite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)

	color.Yellow("ID Recieved ✔️")
	color.Yellow(params["id"])
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

	iter := client.Collection("sites").Where("siteid", "==", params["id"]).Documents(ctx)
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
		jsonString, _ := json.Marshal(doc.Data())

		//convert to Form struct
		json.Unmarshal(jsonString, &nyData)
		//fmt.Println(s.DateSubmitted)
		//fmt.Println(string(jsonString))

		// //fmt.Println(doc.Data())
		// //fmt.Println(nyData)

	}
	fmt.Printf("Document data: %#v\n", nyData)
	json.NewEncoder(w).Encode(nyData)

}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)

	color.Yellow("ID Recieved ✔️")
	color.Yellow(params["id"])
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

	dsnap, err2 := client.Collection("users").Doc(params["id"]).Get(ctx)
	if err2 != nil {
		log.Fatal(err2)
	}

	var c LoggedInUser
	dsnap.DataTo(&c)

	fmt.Printf("Document data: %#v\n", c)
	log.Println("User retrieved - ", params["id"])
	json.NewEncoder(w).Encode(c)
}

func getPatients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)

	color.Yellow("ID Recieved ✔️")
	color.Yellow(params["id"])
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
		f = append(f, s)

	}
	json.NewEncoder(w).Encode(f)
}

func updateVisit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)

	color.Yellow("ID Recieved ✔️")
	color.Yellow(params["id"])
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

	dsnap, err2 := client.Collection("visits").Doc(params["id"]).Get(ctx)
	if err2 != nil {
		log.Fatal(err2)
	}

	var c Patient
	dsnap.DataTo(&c)

	fmt.Printf("Document data: %#v\n", c)
	log.Println("Patient retrieved - ", params["id"])
	json.NewEncoder(w).Encode(c)
}

func authLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
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

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {

			log.Fatal(err)
		}

		//fmt.Println(string(body))

		var p Login

		json.Unmarshal([]byte(body), &p)
		color.Yellow("Login Attempt Recieved from: " + p.Username)
		log.Println("Login Attempt Recieved from: " + p.Username)

		iter := client.Collection("users").Where("username", "==", p.Username).Where("password", "==", p.Password).Where("siteid", "==", p.Siteid).Documents(ctx)

		for {
			doc, err2 := iter.Next()

			if err2 == iterator.Done {

				break
			}

			doc.Data()
			color.Green("Login Successful from: " + p.Username)
			log.Println("Login Successful from: " + p.Username)
			w.Write([]byte(`{"success":true, "id":"` + doc.Ref.ID + `"}`))

		}

	default:

	}

}
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

	currentTime := time.Now().Add(-10 * time.Minute)
	snapIter := client.Collection("form").OrderBy("DateSubmitted", firestore.Asc).Snapshots(ctx)

	defer snapIter.Stop()

	for {
		snap, err := snapIter.Next()
		if err != nil {
			color.Yellow("", err)
		}

		for _, diff := range snap.Changes {

			jsonString, _ := json.Marshal(diff.Doc.Data())

			s := Form{}
			//convert to Form struct
			json.Unmarshal(jsonString, &s)

			if s.DateSubmitted > currentTime.Format("02/01/2006, 15:04:05") {

				color.Cyan("New Form Received 🔔")
				color.Cyan("Date and time submitted: %v", s.DateSubmitted)
			} else {

			}

		}

	}
}

func main() {
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

	r := mux.NewRouter()
	r.HandleFunc("/getPatient/{id}", getPatient).Methods("POST")
	r.HandleFunc("/getVisits/{id}", getVisits).Methods("POST")
	r.HandleFunc("/getUser/{id}", getUser).Methods("POST")
	r.HandleFunc("/getPatients", getPatients).Methods("POST")
	r.HandleFunc("/updateVisit/{id}", updateVisit).Methods("POST")
	r.HandleFunc("/updateV/{id}", updateVisit1).Methods("POST")
	r.HandleFunc("/getSite/{id}", getSite).Methods("POST")
	r.HandleFunc("/login/", authLogin).Methods("POST")

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

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.

	handler := cors.Default().Handler(r)
	http.ListenAndServe(":8080", handler)
	checkForNewForm()

}
