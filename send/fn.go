package send

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Holiday struct {
	Day     string
	Date    string
	HDay    string
	Comment string
}

var PublicHoliday map[string]Holiday

func StringToTime(date string, timein string) time.Time {
	/*myformat := date + "T" + timein + "+00:00"
	t1, e := time.Parse(
		time.RFC3339,
		myformat)*/

	mf := fmt.Sprintf("%s 00:00:00", date)
	t1, e := dateparse.ParseAny(mf)

	if e != nil {
		log.Println("ERROR DATE PARSE > ", date, mf, " > ", e)

	}
	return t1
}

func FindNextDateFromStartNumberDays(startDate string, numDays float64) string {

	currentDate := startDate
	var counter float64
	for {
		inDate := StringToTime(currentDate, "00:00:00")
		nextDate := strings.Split(inDate.Add(time.Duration(24*1)*time.Hour).String(), " ")[0]
		if IsThisWorkingDay(nextDate, false, false) {
			counter++
		}
		currentDate = nextDate
		if counter == numDays {
			break
		}
	}
	return currentDate
}
func IsThisWorkingDay(dateIn string, WorkOnSaturday bool, WorkOnSunday bool) bool {

	booPub := IsThisPublicHoliday(dateIn)
	if booPub { // this is the public holiday
		return false
	}
	weekDays := []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}
	var workingDays []string
	for _, dayin := range weekDays {
		if dayin == "saturday" {
			if WorkOnSaturday {
				workingDays = append(workingDays, strings.ToLower(dayin))
			}
		} else if dayin == "sunday" {
			if WorkOnSunday {
				workingDays = append(workingDays, dayin)
			}
		} else {
			workingDays = append(workingDays, dayin)
		}
	}
	//todo
	targetDate := StringToTime(dateIn, "00:00:00")
	day := targetDate.Weekday().String()
	boo := IsInArray(workingDays, strings.ToLower(day))
	return boo
}

func AddDayToDate(startDate string, numberDay int) string {
	date := StringToTime(startDate, "00:00:00")
	nextDate := strings.Split(date.Add(time.Duration(24*numberDay)*time.Hour).String(), " ")[0]
	return nextDate
}

func SubstructDayToDate(startDate string, numberDay int) string {
	date := StringToTime(startDate, "00:00:00")
	nextDate := strings.Split(date.Add(time.Duration(-24*numberDay)*time.Hour).String(), " ")[0]
	return nextDate
}

func fnWitchDataType(myInterface interface{}) string {
	switch myInterface.(type) {
	case int:
		return "int"
	case int64:
		return "int"
	case int32:
		return "int"
	case float64:
		return "float"
	case float32:
		return "float"
	case string:
		return "string"
	case bool:
		return "bool"
	case map[string]float64:
		return "map"
	case map[string]string:
		return "map"
	case map[string]interface{}:
		return "map"
	default:
		return "none"
	}
}

func RecoverMe(module string) {
	if r := recover(); r != nil {
		//fmt.Println("-->> RECOVERD FROM "+module+" > , ",r)
	}
}

func LoadPublicHoliday() map[string]Holiday {
	PublicHoliday = make(map[string]Holiday)
	PublicHoliday["01-01"] = Holiday{Day: "Monday", Date: "January 01", HDay: "New Year's day"}
	PublicHoliday["03-21"] = Holiday{Day: "Wednesday", Date: "March 21", HDay: "Human Rights day"}
	PublicHoliday["03-30"] = Holiday{Day: "Friday", Date: "March 30", HDay: "Good Friday"}
	PublicHoliday["04-02"] = Holiday{Day: "Monday", Date: "April 02", HDay: "Familly"}
	PublicHoliday["04-27"] = Holiday{Day: "Friday", Date: "April 27", HDay: "National day"}
	PublicHoliday["05-01"] = Holiday{Day: "Tuesday", Date: "May 01", HDay: "Labour day"}
	PublicHoliday["06-16"] = Holiday{Day: "Saturday", Date: "June 16", HDay: "Youth day"}
	PublicHoliday["08-09"] = Holiday{Day: "Thurday", Date: "August 16", HDay: "national women day"}
	PublicHoliday["09-24"] = Holiday{Day: "Monday", Date: "September 24", HDay: "Heritage day"}
	PublicHoliday["12-16"] = Holiday{Day: "Sunday", Date: "December 16", HDay: "Day of Reconciliation"}
	PublicHoliday["12-16"] = Holiday{Day: "Sunday", Date: "December 16", HDay: "Day of Reconciliation"}
	PublicHoliday["12-17"] = Holiday{Day: "Monday", Date: "December 17", HDay: "Public Holiday as reconciliation day falls on a sunday 2018"}
	PublicHoliday["12-25"] = Holiday{Day: "Tuesday", Date: "December 25", HDay: "Christmas day"}
	PublicHoliday["12-26"] = Holiday{Day: "Wednesday", Date: "December 26", HDay: "Day of Good Will"}
	return PublicHoliday
}

func DateValidationFormat_YYYY_MM_DD(dateIn string) (bool, string) {
	arr := strings.Split(dateIn, "-")
	if len(arr) != 3 {
		return false, fmt.Sprintf("missing date part %v/%v", len(arr), 3)
	}

	//todo validate year
	strYear := arr[0]
	intYear, _ := strconv.Atoi(strYear)
	if len(strYear) != 4 {
		return false, fmt.Sprintf("year character must be 4 by got %v", len(strYear))
	}
	if intYear < 2018 {
		return false, fmt.Sprintf("year  must be equal or great then 2018 by got %v", intYear)
	}

	//todo validate month
	strMonth := arr[1]
	intMonth, _ := strconv.Atoi(strMonth)
	if len(strMonth) != 2 {
		return false, fmt.Sprintf("month character must be 2 by got %v", len(strMonth))
	}
	if intMonth < 1 || intMonth > 12 {
		return false, fmt.Sprintf("month  must be between 1 to 12  by got %v", intMonth)
	}

	//todo validate day
	strDay := arr[2]
	intDay, _ := strconv.Atoi(strDay)
	if len(strDay) != 2 {
		return false, fmt.Sprintf("day character must be 2 by got %v", len(strMonth))
	}
	if intDay < 1 || intDay > 31 {
		return false, fmt.Sprintf("day  must be between 1 to 31  by got %v", intDay)
	}

	return true, ""
}

func Extract9LastDigitNumberFromString(numIn, prefix string) string {
	inputFmt := numIn[len(numIn)-9 : len(numIn)]
	return prefix + inputFmt
}

func DisplayObject(title string, obj interface{}) {
	str, err := json.Marshal(obj)
	dis := fmt.Sprintf("->[ %s ] => %s", title, string(str))
	if err != nil {
		log.Println("Error DisplayObject > ", err, title)
	}
	log.Println(dis)
}
func ConvertDateShortNameToDigit(str string) string {
	mymap := make(map[string]string)
	mymap["Jan"] = "01"
	mymap["Feb"] = "02"
	mymap["Mar"] = "03"
	mymap["Apr"] = "04"
	mymap["May"] = "05"
	mymap["Jun"] = "06"
	mymap["Jul"] = "07"
	mymap["Aug"] = "08"
	mymap["Sep"] = "09"
	mymap["Oct"] = "10"
	mymap["Nov"] = "11"
	mymap["Dec"] = "12"

	val, _ := mymap[str]
	return val
}

func GetHttpFileContent2(myurl string) ([]byte, string) {
	defer RecoverMe("GetHttpFileContent2")
	myurltoken := strings.Split(myurl, "/")
	filename := myurltoken[(len(myurltoken) - 1)]
	// Just a simple GET request to the image URL
	// We get back a *Response, and an error
	res, err := http.Get(myurl)

	if err != nil {
		log.Fatalf("http.Get -> %v", err)
	}

	// We read all the bytes of the image
	// Types: data []byte
	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatalf("ioutil.ReadAll -> %v", err)
	}

	// You have to manually close the body, check docs
	// This is required if you want to use things like
	// Keep-Alive and other HTTP sorcery.Response body was
	res.Body.Close()

	// You can now save it to disk or whatever...
	if filename == "policy-schedule" {
		filename = "policy-schedule_" + uuid.New() + ".pdf"
	}
	if filename == "policy-wording" {
		filename = "policy-wording _" + uuid.New() + ".pdf"
	}
	ioutil.WriteFile("./attachement/"+filename, data, 0777)

	//log.Println("I saved your image buddy! > " + filename)

	return data, filename
}
func SecureRequestHTTPS(myUrl string) string {
	defer RecoverMe("SecureRequestHTTPS")
	//myUrl="https://maps.googleapis.com/maps/api/distancematrix/json?origins=14.614786,121.046587&destinations=14.610301,121.080233&mode=driving&language=en&departure_time=now&key=AIzaSyBqmizvy1TNIFdvM1whOXFJqb45Nl8hCIs"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Get(myUrl)
	if err != nil {
		log.Println(err)
		return "{}"
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println(fmt.Sprintf("Response code was %v; want 200", res.StatusCode))
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return "{}"
	}
	expected := []byte("Hello World")
	if bytes.Compare(expected, body) != 0 {
		//log.Println(fmt.Sprintf("Response body was '%v'; want '%v'", expected, string(body)))
	}
	if err != nil {
		return "{}"
	}
	return string(body)
}

func httpContentPost(myUrl string, mymap map[string]interface{}) string {

	data := url.Values{}
	for key, val := range mymap {
		strVal := fmt.Sprintf("%v", val)
		data.Set(key, strVal)
	}

	req, err := http.NewRequest("POST", myUrl, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "{}"
	}
	//req.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	//req.Header.Set("Connection", "Keep-Alive")

	var client http.Client
	resp, err := client.Do(req)

	if err != nil {
		return "{}"
	}
	defer resp.Body.Close()

	dataOut, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "{}"
	}

	return string(dataOut)

}

func GetCsvFileContentUrl2(filename string) string {
	dirname := filename
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		// dirname does not exist
		//dirname = "/imbani" + string(filepath.Separator) + "cassandra.conf"
	}

	d, err := os.Open(dirname)
	if err != nil {
		log.Println("ERROR opening file > ", err)
		return ""
	}
	defer d.Close()
	buf := bytes.NewBuffer(nil)
	f, err := os.Open(dirname)
	if err != nil {
		return ""
	}
	io.Copy(buf, f)
	f.Close()
	s := string(buf.Bytes())
	/*
		let build the data now
	*/
	d.Close()
	return s
}

func HttpContentPost(myUrl string, mymap map[string]interface{}) string {

	data := url.Values{}
	for key, val := range mymap {
		strVal := fmt.Sprintf("%v", val)
		data.Set(key, strVal)
	}

	strData, _ := json.Marshal(mymap)
	var jsonStr = []byte(string(strData))

	req, err := http.NewRequest("POST", myUrl, bytes.NewBuffer(jsonStr))

	if err != nil {
		fmt.Println("Error NewRequest  ", err)
		return "{}"
	}
	//req.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	//req.Header.Set("Connection", "Keep-Alive")

	var client http.Client
	resp, err := client.Do(req)

	if err != nil {
		return "{}"
	}
	defer resp.Body.Close()

	dataOut, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "{}"
	}

	return string(dataOut)

}

func WriteNewLineToLogFile(line, filename string) {

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	if _, err = f.WriteString(line + "\n"); err != nil {
		panic(err)
	}
	defer f.Close()
}

func ConvertBase64ToString(data string) string {
	sDesc, _ := b64.StdEncoding.DecodeString(data)
	//fmt.Println(string(sDesc))
	return string(sDesc)
}

func ConvertStringToBase64(data string) string {
	sEnc := b64.StdEncoding.EncodeToString([]byte(data))
	//fmt.Println(sEnc)
	return sEnc
}

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetCsvFileContentUrl(filename string) ([]string, error) {
	dirname := filename
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		// dirname does not exist
		//dirname = "/imbani" + string(filepath.Separator) + "cassandra.conf"
	}

	d, err := os.Open(dirname)
	if err != nil {

		return []string{}, err
	}
	defer d.Close()
	buf := bytes.NewBuffer(nil)
	f, err := os.Open(dirname)
	if err != nil {
		return []string{}, err
	}
	io.Copy(buf, f)
	f.Close()
	s := string(buf.Bytes())
	/*
		let build the data now
	*/
	upList := []string{}
	lines := strings.Split(s, "\n")
	for _, oneline := range lines {
		upList = append(upList, oneline)
	}
	d.Close()
	return upList, nil
}

func MAP(obj interface{}) map[string]string {
	mymap := make(map[string]string)
	s := structs.New(obj)
	f := s.Fields()
	for _, row := range f {
		mymap[row.Name()] = row.Kind().String()
	}
	return mymap
}

func CreateFolderUploadIfNotExist(path string) {
	_ = os.MkdirAll(path, os.ModePerm)
}

func GetHttpFileContent(myurl string, dirName string) ([]byte, string) {
	myurltoken := strings.Split(myurl, "/")
	filename := myurltoken[(len(myurltoken) - 1)]
	// Just a simple GET request to the image URL
	// We get back a *Response, and an error
	res, err := http.Get(myurl)

	if err != nil {
		log.Fatalf("http.Get -> %v", err)
	}

	// We read all the bytes of the image
	// Types: data []byte
	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatalf("ioutil.ReadAll -> %v", err)
	}

	// You have to manually close the body, check docs
	// This is required if you want to use things like
	// Keep-Alive and other HTTP sorcery.
	res.Body.Close()

	// You can now save it to disk or whatever...
	if filename == "policy-schedule" {
		filename = "policy-schedule_" + uuid.New() + ".pdf"
	}
	if filename == "policy-wording" {
		filename = "policy-wording _" + uuid.New() + ".pdf"
	}
	ioutil.WriteFile(dirName+"/"+filename, data, 0777)

	//log.Println("I saved your image buddy! > " + filename)

	return data, filename
}
func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		//os.Exit(0)
	}
}

func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func GetDateAndTimeString() (string, string) {
	mydate := time.Now()
	arr := strings.Split(fmt.Sprintln(mydate.Format("2006-01-02 15:04:05")), " ")
	date := arr[0]
	time := arr[1]
	return strings.TrimSpace(date), strings.TrimSpace(time)
}

var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!")

func NewPassword(length int) string {
	return rand_char(length, StdChars)
}

func rand_char(length int, chars []byte) string {
	new_pword := make([]byte, length)
	random_data := make([]byte, length+(length/4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, random_data); err != nil {
			panic(err)
		}
		for _, c := range random_data {
			if c >= maxrb {
				continue
			}
			new_pword[i] = chars[c%clen]
			i++
			if i == length {
				return string(new_pword)
			}
		}
	}
	panic("unreachable")
}

func GetMd5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	out := fmt.Sprintf("%x", h.Sum(nil))
	//fmt.Println("GetMd5 > ",out," > ",str)
	return out
}

func CleanToLower(str string) string {
	str = strings.Trim(str, " ")
	//str = strings.ToLower(str)
	return str
}

func HttpContentPostJson(myUrl string, mymap map[string]interface{}) string {
	data := url.Values{}
	for key, val := range mymap {
		strVal := fmt.Sprintf("%v", val)
		data.Set(key, strVal)
	}
	strData, _ := json.Marshal(mymap)
	var jsonStr = []byte(string(strData))
	req, err := http.NewRequest("POST", myUrl, bytes.NewBuffer(jsonStr))

	if err != nil {
		fmt.Println("Error NewRequest  ", err)
		return "{}"
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	var client http.Client
	resp, err := client.Do(req)

	if err != nil {
		return "{}"
	}
	defer resp.Body.Close()
	dataOut, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "{}"
	}
	return string(dataOut)
}

func PublishToReact(w http.ResponseWriter, r *http.Request, obj interface{}, htmlcode int) {
	type ReactResponse struct {
		Data interface{}
	}
	myresp := ReactResponse{}
	myresp.Data = obj

	w.WriteHeader(htmlcode)
	myw, _ := json.Marshal(obj)
	w.Write(myw)
}

func GetPostedDataMapAndString(r *http.Request) (map[string]interface{}, string) {
	mymap := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&mymap)
	if err != nil {
		//log.Println("ERROR CREDIT NEW > ",err)
		emp := make(map[string]interface{})
		return emp, "{}"
	}
	defer r.Body.Close()

	strP, _ := json.Marshal(mymap)
	strJ := string(strP)
	//fmt.Println("GetPostedDataMapAndString ******> ",strJ)
	return mymap, strJ
}
func IsInArray(strArray []string, val string) bool {
	boo := false
	for _, item := range strArray {
		if item == val {
			boo = true
			break
		}
	}
	return boo
}

func IsThisPublicHoliday(dateIn string) bool {
	arr := strings.Split(dateIn, "-")
	str := fmt.Sprintf("%v-%v", arr[1], arr[2])
	pub := LoadPublicHoliday()
	_, ok := pub[str]
	return ok
}
func LoadPublicHoliday1() map[string]Holiday {
	PublicHoliday = make(map[string]Holiday)
	PublicHoliday["01-01"] = Holiday{Day: "Tuesday", Date: "January 01", HDay: "New Year's day"}
	PublicHoliday["03-21"] = Holiday{Day: "Thursday", Date: "March 21", HDay: "Human Rights day"}
	PublicHoliday["04-19"] = Holiday{Day: "Friday", Date: "April 19", HDay: "Good Friday"}
	PublicHoliday["04-22"] = Holiday{Day: "Monday", Date: "April 22", HDay: "Family"}
	PublicHoliday["04-27"] = Holiday{Day: "Saturday", Date: "April 27", HDay: "National day"}
	PublicHoliday["05-01"] = Holiday{Day: "Wednesday", Date: "May 01", HDay: "Labour day"}
	PublicHoliday["05-08"] = Holiday{Day: "Wednesday", Date: "May 08", HDay: "Election day"}
	PublicHoliday["06-16"] = Holiday{Day: "Sunday", Date: "June 16", HDay: "Youth day"}
	PublicHoliday["06-17"] = Holiday{Day: "Monday", Date: "June 17", HDay: "Youth day Holiday"}
	PublicHoliday["08-09"] = Holiday{Day: "Friday", Date: "August 09", HDay: "national women day"}
	PublicHoliday["09-24"] = Holiday{Day: "Tuesday", Date: "September 24", HDay: "Heritage day"}
	PublicHoliday["12-16"] = Holiday{Day: "Monday", Date: "December 16", HDay: "Day of Reconciliation"}
	//PublicHoliday["12-17"] = Holiday{Day: "Monday", Date: "December 17", HDay: "Public Holiday as reconciliation day falls on a sunday 2018"}
	PublicHoliday["12-25"] = Holiday{Day: "Wednesday", Date: "December 25", HDay: "Christmas day"}
	PublicHoliday["12-26"] = Holiday{Day: "Thursday", Date: "December 26", HDay: "Day of Good Will"}
	return PublicHoliday
}

