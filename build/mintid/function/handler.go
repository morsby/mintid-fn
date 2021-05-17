package function

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/morsby/mintid"
)

func getAPISecret(secretName string) (secretBytes []byte, err error) {
	// read from the openfaas secrets folder
	secretBytes, err = ioutil.ReadFile("/var/openfaas/secrets/" + secretName)

	return secretBytes, err
}

func Handle(w http.ResponseWriter, r *http.Request) {
	cpr, err := getAPISecret("cpr")
	if err != nil {
		http.Error(w, "error getting credentials", http.StatusInternalServerError)
	}
	fmt.Println("CPR:" + string(cpr))
	pwd, err := getAPISecret("pwd")
	if err != nil {
		http.Error(w, "error getting credentials", http.StatusInternalServerError)
	}
	person, err := mintid.Login(string(cpr), string(pwd))
	if err != nil {
		panic(err)
	}

	shifts, _ := person.Fetch("202101010000", "202112310000")

	calendar, _ := mintid.CreateCalendar(shifts, "PUBLISH", "bf", "Fri efter vagt", "Blank dag")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(calendar))
}
