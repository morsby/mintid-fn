package function

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/morsby/mintid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getAPISecret(secretName string) (secretBytes []byte, err error) {
	path := "/var/openfaas/secrets/" + secretName

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return nil, nil
	}

	// read from the openfaas secrets folder
	secretBytes, err = ioutil.ReadFile(path)

	return secretBytes, err
}

var db *gorm.DB

type Person struct {
	gorm.Model
	CPR      string `json:"cpr" gorm:"type:varbinary(512)"`
	Password string `json:"password" gorm:"type:varbinary(512)"`
	Code     string `json:"code" gorm:"unique"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	key, err := getAPISecret("aes")
	if err != nil {
		panic(err)
	}

	dsn, err := getAPISecret("mysql")
	if err != nil {
		panic(err)
	}

	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN: string(dsn) + "?charset=utf8&parseTime=True&loc=Local",
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&Person{})

	if r.Method == http.MethodGet {
		if r.URL.Path == "/" {
			w.Write([]byte("Du er på forsiden"))
		} else {
			var person Person
			db.First(&person, "code = ?", r.URL.Path[1:])
			cpr, err := decrypt([]byte(person.CPR), key)
			if err != nil {
				http.Error(w, "error decrypting cpr: "+err.Error(), http.StatusInternalServerError)
			}

			pwd, err := decrypt([]byte(person.Password), key)
			if err != nil {
				http.Error(w, "error decrypting password", http.StatusInternalServerError)
			}

			auth, err := mintid.Login(string(cpr), string(pwd))
			if err != nil {
				http.Error(w, "error logging into MedarbejderNet", http.StatusInternalServerError)
			}

			shifts, _ := auth.Fetch("202101010000", "202112310000")

			calendar, _ := mintid.CreateCalendar(shifts, "PUBLISH", "bf", "Fri efter vagt", "Blank dag")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(calendar))

		}
	} else if r.Method == http.MethodPost {
		var person Person
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error unmarshalling request", http.StatusBadRequest)
		}
		json.Unmarshal(body, &person)

		// generate uuid
		uuidWithHyphhen := uuid.New()
		uuidStr := strings.Replace(uuidWithHyphhen.String(), "-", "", -1)
		person.Code = uuidStr

		// encrypt info
		cpr, err := encrypt([]byte(person.CPR), key)
		if err != nil {
			http.Error(w, "error encrypting cpr", http.StatusInternalServerError)
		}
		person.CPR = string(cpr)

		pwd, err := encrypt([]byte(person.Password), key)
		if err != nil {
			http.Error(w, "error encrypting pwd", http.StatusInternalServerError)
		}
		person.Password = string(pwd)

		db.Create(&person)

		w.Write([]byte("Din kalender findes på: https://faasd.morsby.dk/function/mintid/" + person.Code))
	}

}
