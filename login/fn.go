package login

import (
	"crypto/md5"
	"fmt"
	"io"
)

func GetMd5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	out := fmt.Sprintf("%x", h.Sum(nil))
	//fmt.Println("GetMd5 > ",out," > ",str)
	return out
}
