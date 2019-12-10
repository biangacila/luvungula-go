package files

import (
	"bytes"
	"encoding/base64"
	"github.com/biangacila/luvungula-go/dropbox"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type ServiceUpload struct {
	Org           string
	TableRef      string
	Ref           string
	Category      string
	Type          string
	Name          string
	Username      string
	Filename      string
	dropboxRef    string
	Link          string
	Base64String  string
	docType       string
	inRef         string
	fileExtension string

	size     int
	revision string
	path     string

	DROPBOX_ROOT string
	DROPBOX_CLIENTID string
	DROPBOX_SECRET string
	DROPBOX_TOKEN string
	DOWNLOAD_SERVICE_LINK string
}

 const TEMP_DIR="tmp-file-upload-dropbox"

func (obj *ServiceUpload) New() string {
	obj.docType = obj.Category
	//dt, hr := GetDateAndTimeString()
	obj.inRef = fmt.Sprintf("%v", obj.Ref)
	fmt.Println("Service upload start for -> ", obj.inRef)
	obj.SaveToDropbox()
	fmt.Println("Service upload end for -> ", obj.inRef, " > ", obj.revision)

	return obj.dropboxRef
}
func (obj *ServiceUpload) getRef() string {
	ref := fmt.Sprintf("%v/%v/%v/%v.%s", obj.Org, obj.TableRef, obj.Category, obj.inRef, obj.fileExtension)
	return ref
}
func (obj *ServiceUpload) SaveToDropbox() (string, int, string) {
	//todo create file into the temp directory
	src := obj.CreateFileIntoTempDir(TEMP_DIR, obj.Filename, obj.Base64String)
	busket := obj.DROPBOX_ROOT
	fmt.Println("::-->> ", src, " > ", obj.Filename, " > ", obj.fileExtension)
	myFilter := obj.getRef()
	dst := busket + "" + obj.Org + "/" + myFilter
	bdb := dropbox.HubDropbox{}
	//todo authentication to dropbox
	bdb.DROPBOX_CLIENTID = obj.DROPBOX_CLIENTID
	bdb.DROPBOX_SECRET = obj.DROPBOX_SECRET
	bdb.DROPBOX_TOKEN = obj.DROPBOX_TOKEN

	bdb.Company = obj.Org
	bdb.UploadFilename = dst
	bdb.UploadBusket = busket
	bdb.UploadSrc = src
	entry := bdb.Upload(src, dst)

	obj.revision = entry.Revision
	obj.size = entry.Size
	obj.path = entry.Path

	obj.Link = obj.DOWNLOAD_SERVICE_LINK + obj.revision
	obj.dropboxRef = obj.revision

	return entry.Revision, entry.Size, entry.Path
}

func (obj *ServiceUpload) CreateFileIntoTempDir(dirName, filename, Base64String string) string {
	obj.createFolderUploadIfNotExist(dirName)
	p1 := strings.Split(Base64String, ";")[0]
	str := fmt.Sprintf("%s;base64,", p1)
	base64String := strings.Replace(Base64String, str, "", 1)
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64String))
	buff := bytes.Buffer{}
	_, err := buff.ReadFrom(reader)
	if err != nil {
		log.Println("buff.ReadFrom err -->> ", err, " > ", base64String)
		return ""
	}
	delFileDirName := fmt.Sprintf("%s/%s", TEMP_DIR, filename)
	fmt.Println("delFileDirName -->> ", delFileDirName)
	err = ioutil.WriteFile(delFileDirName, buff.Bytes(), 0644)
	fmt.Println("ioutil.WriteFile err -->> ", err)

	arr := strings.Split(obj.Filename, ".")
	extr := arr[(len(arr) - 1)]
	obj.fileExtension = extr

	return delFileDirName
}

func(obj *ServiceUpload) createFolderUploadIfNotExist(path string) {
	_ = os.MkdirAll(path, os.ModePerm)
}

