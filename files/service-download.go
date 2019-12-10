package files

import (
	"fmt"
	"github.com/biangacila/luvungula-go/dropbox"
)

type ServiceDownload struct {
	DROPBOX_CLIENTID string
	DROPBOX_SECRET string
	DROPBOX_TOKEN string
}

func (obj *ServiceDownload) DownloadByRev(myRev string) string {
	bdb := dropbox.HubDropbox{}
	bdb.DROPBOX_CLIENTID = obj.DROPBOX_CLIENTID
	bdb.DROPBOX_SECRET = obj.DROPBOX_SECRET
	bdb.DROPBOX_TOKEN = obj.DROPBOX_TOKEN
	bdb.DownloadRev = myRev
	fmt.Println("[DOWNLOAD REVISION REQUEST] > ", myRev)
	saveLink := bdb.Download()
	return saveLink
}
