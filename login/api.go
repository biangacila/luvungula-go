package login

import (
	"encoding/json"
	"fmt"
	"github.com/biangacila/luvungula-go/global"
	"net/http"
	"strings"
)

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

func UserDetailHandler(w http.ResponseWriter, r *http.Request) {
	myAuthorization := ""
	for k, v := range r.Header {
		myKey := fmt.Sprintf("%v", k)
		myVal := fmt.Sprintf("%v", v[0])
		if strings.Contains(myKey, "Authorization") {
			arr := strings.Split(myVal, "Bearer ")
			myAuthorization = arr[1]

		}
	}

	//todo check if the token is generate by us
	user, isFind := TokensList[myAuthorization]
	if !isFind {
		w.WriteHeader(http.StatusUnauthorized)
		my := make(map[string]interface{})
		my["error"] = "Unauthorised access to this resource > not find"
		str, _ := json.Marshal(my)
		strJson := string(str)
		fmt.Fprint(w, strJson)
		return
	}

	my := make(map[string]interface{})
	my["STATUS"] = "OK"
	my["REQUEST"] = myAuthorization
	my["RESULT"] = user

	JsonResponse(my, w)

}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	global.RecoverMe("LoginHandler")
	var serviceLogin ServiceLogin
	//decode request into UserCredentials struct
	err := json.NewDecoder(r.Body).Decode(&serviceLogin)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Error in request")
		return
	}

	//validate user credentials
	boo, user, tokenString, msg := serviceLogin.NewLogin()
	if !boo {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, msg)
		return
	}

	fmt.Println("-->OK Request login from: ", serviceLogin.Username, serviceLogin.Password)
	//response := Token{tokenString,user}
	data := make(map[string]interface{})
	data["boo"] = boo
	data["user"] = user
	data["token"] = tokenString
	data["msg"] = msg

	my := make(map[string]interface{})
	my["STATUS"] = "OK"
	my["REQUEST"] = serviceLogin
	my["DATA"] = data

	JsonResponse(data, w)
}

//AUTH TOKEN VALIDATION
func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	myAuthorization := ""
	for k, v := range r.Header {
		myKey := fmt.Sprintf("%v", k)
		myVal := fmt.Sprintf("%v", v[0])
		if strings.Contains(myKey, "Authorization") {
			arr := strings.Split(myVal, "Bearer ")
			myAuthorization = arr[1]

		}
	}

	//todo check if the token is generate by us
	my := make(map[string]interface{})
	my["error"] = "Unauthorised access to this resource"
	str, _ := json.Marshal(my)
	strJson := string(str)
	_, isFind := TokensList[myAuthorization]
	if !isFind {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, strJson)
		return
	}

	if isVal, err := IsValidToken(myAuthorization); !isVal {
		my["error"] = "Unauthorised access to this resource " + err.Error()
		str, _ := json.Marshal(my)
		strJson = string(str)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, strJson)
		return
	}

	next(w, r)

}

//HELPER FUNCTIONS

func JsonResponse(response interface{}, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
