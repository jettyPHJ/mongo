package set

import (
	"fmt"
	"io/ioutil"

	mapset "github.com/deckarep/golang-set"
)

var Path string = "/home/jette/192.168.1.163/cloud_trackfile/testSimulator/tem/"
var FileNameSet mapset.Set

func init() {
	fileNames := GetAllFile(Path)
	FileNameSet = mapset.NewSetFromSlice(fileNames)
	FileNameSet.Remove(nil)
}

func GetAllFile(pathname string) []interface{} {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println(err)
	}
	fileNames := make([]interface{}, 1)
	for _, fi := range rd {
		if !fi.IsDir() {
			fileNames = append(fileNames, fi.Name())
		}
	}
	return fileNames
}
