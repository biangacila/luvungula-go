package login

import (
	"encoding/json"

	"fmt"
	"github.com/biangacila/luvungula-go/global"
	"github.com/biangacila/luvungula-go/io"
	"github.com/robbert229/jwt"
	"strings"
	"time"
)

var User interface{}

type ServiceLogin struct {
	Username  string
	Password  string
	DB_NAME   string
	TableName string

	FieldPassword string
	FieldUsername string
	HashPassword  bool
	/* Working variable */
	token    string
	user     interface{}
	hasLogin bool
}

func (obj *ServiceLogin) NewLogin() (bool, interface{}, string, string) {

	//todo let find out if our user exist base on the username
	user, okUser := obj.checkIfUsernameExist()
	if !okUser {
		return false, nil, "", "Invalid credentials"
	}
	obj.user = user

	//todo let compare passport
	fPass := strings.ToLower(obj.FieldPassword)
	vPass, _ := user[fPass]
	password := fmt.Sprintf("%v", vPass)
	if obj.HashPassword {
		if password != GetMd5(obj.Password) {
			return false, nil, "", "Invalid credentials"
		}
	} else {
		if password != obj.Password {
			return false, nil, "", "Invalid credentials"
		}
	}

	user[fPass] = ""
	obj.user = user

	//todo let create our token
	token := obj.createToken()

	//todo let send out our success result

	return true, user, token, ""
}
func (obj *ServiceLogin) createToken() string {
	secret := SignKey
	algorithm := jwt.HmacSha256(string(secret))
	claims := jwt.NewClaim()
	str, _ := json.Marshal(obj.user)
	my := make(map[string]interface{})
	json.Unmarshal(str, &my)
	for key, val := range my {
		value := fmt.Sprintf("%v", val)
		claims.Set(key, value)
	}
	claims.SetTime("exp", time.Now().Add(8*time.Hour))
	dt, hr := global.GetDateAndTimeString()
	claims.Set("date", dt)
	claims.Set("time", hr)
	token, err := algorithm.Encode(claims)
	tokenString := fmt.Sprintf("%s", token)
	if algorithm.Validate(token) != nil {
		panic(err)
	}
	TokensList[tokenString] = obj.user
	return tokenString
}

func (obj *ServiceLogin) checkIfUsernameExist() (map[string]interface{}, bool) {
	var ls *[]interface{}
	qry := fmt.Sprintf("select * from %v.%v where %v='%v' ",
		obj.DB_NAME, obj.TableName, obj.FieldUsername, obj.FieldUsername)
	res := io.RunQueryCass(qry)
	global.ConvertQueryResultToAnyStruct(ls, res)
	my := global.ConvertInterfaceIntoMapArray(ls)

	if len(my) > 0 {
		return my[0], true
	}
	return nil, false
}
