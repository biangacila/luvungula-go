package email

import (
	"crypto/tls"
	"fmt"
	"github.com/go-mail/mail"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type EmailHost struct {
	Username    string
	Sender      string
	Receiver    string
	Subject     string
	Body        string
	ReplyTo     string
	FromEmail   string
	FromCompany string

	HasAttached  string
	Attaches     string
	AttachedType string

	SMTP_PORT     string
	SMTP_SERVER   string
	SMTP_USERNAME string
	SMTP_PASSWORD string
}

const DIR_TEMP_ATTACHED_DOWNLOAD = "downloads"

func (obj *EmailHost) Send() {
	port, _ := strconv.Atoi(obj.SMTP_PORT)
	d := mail.NewDialer(obj.SMTP_SERVER, port, obj.SMTP_USERNAME, obj.SMTP_PASSWORD)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	m := mail.NewMessage(mail.SetCharset("ISO-8859-1"), mail.SetEncoding(mail.Base64))
	fromName := fmt.Sprintf("%s<%s>", obj.FromCompany, obj.FromEmail)
	fmt.Println("fromName :)--> ", fromName)
	m.SetHeaders(map[string][]string{
		"From":    {fromName},
		"To":      {obj.Receiver},
		"Subject": {obj.Subject},
	})
	m.SetBody("text/html", obj.Body)
	CreateFolderUploadIfNotExist(DIR_TEMP_ATTACHED_DOWNLOAD)
	var lsitdelete []string
	attachedList := strings.Split(obj.Attaches, ";")
	if obj.HasAttached == "yes" {
		if obj.AttachedType == "local" {
			for _, filename := range attachedList {
				m.Attach(filename)
			}
		}
		if obj.AttachedType == "external" {
			for _, fileUrl := range attachedList {
				mybite, filename := GetHttpFileContent2(fileUrl)
				attFile := fmt.Sprintf("./%s/%s", DIR_TEMP_ATTACHED_DOWNLOAD, filename)
				ioutil.WriteFile(attFile, mybite, 0644)
				m.Attach(attFile)
				timer := time.NewTimer(1 * time.Second)
				lsitdelete = append(lsitdelete, attFile)
				<-timer.C
			}
		}
	}

	err := d.DialAndSend(m)
	fmt.Println("@-=> EMAIL SEND REPORT : ", err)

	timer2 := time.NewTimer(time.Second * 60)
	go func() {
		<-timer2.C
		fmt.Println("Timer 2 expired")
		for _, fname := range lsitdelete {
			var err = os.Remove(fname)
			CheckError(err)
		}
	}()

}
