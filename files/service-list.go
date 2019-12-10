package files

import (
	"encoding/json"
	"github.com/biangacila/luvungula-go/dropbox"
	"fmt"
	"strings"
)

type SearchResult struct {
	Id       string
	Name     string
	Path     string
	Revision string
	Size     int
	Tag      string
	Modified string
	ParentId string


}

type ServiceList struct {
	Company   string
	TableRef  string
	Category  string
	SearchKey string

	DROPBOX_CLIENTID string
	DROPBOX_SECRET string
	DROPBOX_TOKEN string
}

func (obj *ServiceList) List() []SearchResult {
	// Dropbox
	//Companies
	//ubu
	//easipath
	prefix := fmt.Sprintf("%v--%v", obj.Company, obj.Category)
	//todo
	bdb := dropbox.HubDropbox{}
	bdb.DROPBOX_CLIENTID = obj.DROPBOX_CLIENTID
	bdb.DROPBOX_SECRET = obj.DROPBOX_SECRET
	bdb.DROPBOX_TOKEN = obj.DROPBOX_TOKEN
	bdb.SearchQuery = prefix
	//bdb.SearchQuery = obj.SearchKey
	ls := bdb.SearchByPrefix()
	output := []SearchResult{}
	for _, row := range ls {
		o := SearchResult{}
		str, _ := json.Marshal(row)
		json.Unmarshal(str, &o)
		if obj.SearchKey != "" {
			if strings.Contains(o.Name, obj.SearchKey) {
				o.Name = obj.cleanName(o)
				output = append(output, o)
			}
		} else {
			o.Name = obj.cleanName(o)
			output = append(output, o)
		}
	}

	/*slice.Sort(output[:], func(i, j int) bool {
		return output[i].Modified < output[j].Modified
	})*/

	return output
}
func (obj *ServiceList) ListSearch() {
	//todo
}
func (obj *ServiceList) cleanName(o SearchResult) string {
	prefix := fmt.Sprintf("%v--%v", obj.Company, obj.Category)
	l := strings.Replace(o.Name, prefix, "", 100)
	arr := strings.Split(l, ".")
	ext := arr[(len(arr) - 1)]
	rExt := fmt.Sprintf(".%v", ext)
	l = strings.Replace(o.Name, rExt, "", 100)
	return l
}
