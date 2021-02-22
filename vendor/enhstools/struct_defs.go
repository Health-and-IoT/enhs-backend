package enhstools

//Form Struct - for v1 - not used currently.
type Form struct {
	Address       string `json:"address"`
	Name          string `json:"name"`
	Dob           string `json:"dob"`
	Nok           string `json:"nok"`
	Chinumber     string `json:"chinumber"`
	Allergies     string `json:"allergies"`
	DateSubmitted string `json:"dateSubmitted"`
	ID            string `json:"id"`
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
	Ailment       []string `json:"ailment"`
	DateSubmitted string   `json:"dateSubmitted"`
	Pain          int64    `json:"pain"`
	Patient       string   `json:"patient"`
	Priority      string   `json:"priority"`
	Seen          bool     `json:"seen"`
	Approved      bool     `json:"approved"`
	DocID         string   `json:"docID"`
	ProgList      string   `json:"progList"`
	FinProg       string   `json:"finProg"`
}
