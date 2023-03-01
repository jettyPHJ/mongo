package set_test

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	mapset "github.com/deckarep/golang-set"
)

var path string = "/home/jette/192.168.1.163/cloud_trackfile/testSimulator/tem/"

func Test01(t *testing.T) {
	fileNames := GetAllFile(path)
	set := mapset.NewSetFromSlice(fileNames)
	fmt.Println(set)
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

func Test02(t *testing.T) {
	fileName := "0a5aed69-390e-40e6-bf46-cb4bf340f878.csv"
	file, err := os.Open(path + fileName)
	if err != nil {
		fmt.Println(err)
	}
	csvReader := csv.NewReader(file)
	results, err2 := csvReader.ReadAll()
	if err2 != nil {
		fmt.Println(err)
	}
	fmt.Println(results)
	for _, v := range results {
		fmt.Println(v)
	}
}
