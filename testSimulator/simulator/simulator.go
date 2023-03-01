package simulator

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	l "testSimulator/log"
	"testSimulator/myhttp"
	"testSimulator/set"
	"time"

	uuid "github.com/satori/go.uuid"
)

var filePath string = set.Path
var postPath string = "http://localhost:8080/trajectoryFile"
var control_goroutinNum = make(chan int, 1024)

type Simulators struct {
	Sum       int //模拟数量
	Frequency int //模拟频率
}

func PrintGoroutinNum() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("协程出错：", err)
		}
	}()
	goroutinNum := 0
	go func() {
		for {
			time.Sleep(time.Second)
			go func(tem int) {
				l.ErrorLogger.Println("协程数量为：", tem)
			}(goroutinNum)
		}
	}()
	for {
		goroutinNum = goroutinNum + <-control_goroutinNum
	}
}

func NewSimulators(sum, frequency int) *Simulators {
	return &Simulators{
		Sum:       sum,
		Frequency: frequency,
	}
}

func (s *Simulators) Send(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("协程出错：", err)
		}
	}()
	for i := 0; i < s.Sum; i++ {
		go oneSimulator(s.Frequency, ctx)
	}
}

func oneSimulator(frequency int, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second / time.Duration(frequency)):
			go sendFileAndMeta()
		}
	}
}

func sendFileAndMeta() {
	control_goroutinNum <- 1
	fileName, meta := renameFileAndGetMeta()
	if fileName == "" {
		return
	}
	params := map[string]string{
		"Meta": meta,
	}
	myhttp.Post(postPath, params, filePath+fileName)
	if !set.FileNameSet.Add(fileName) { //将改名后的id放入set
		l.ErrorLogger.Println("set维护出错!!")
	}
	control_goroutinNum <- -1
}

type trackMeta struct {
	FactoryID     string `json:"factoryID" desc:"轨迹文件id" binding:"required"`
	StartTime     int    `json:"startTime" desc:"轨迹文件id" binding:"required"`
	EndTime       int    `json:"endTime" desc:"轨迹文件id" binding:"required"`
	RoadNumber    string `json:"roadNumber" desc:"轨迹文件id" binding:"required"`
	TraceFileName string `json:"traceFileName" desc:"轨迹文件id" binding:"required"`
	TraceID       string `json:"traceID" desc:"轨迹文件id" binding:"required"`
	ObjectClass   string `json:"objectClass" desc:"轨迹文件id" binding:"required"`
	LicensePlate  string `json:"licensePlate" desc:"轨迹文件id" binding:"required"`
	TraceNumber   int    `json:"traceNumber" desc:"轨迹文件id" binding:"required"`
}

var TypeOfObjectClass []string = []string{"car", "two-wheeler", "person", "bus"}

func renameFileAndGetMeta() (fileName, Meta string) {
	var popElem interface{}
	for {
		popElem = set.FileNameSet.Pop()
		if popElem == nil {
			time.Sleep(time.Second)
			continue
		} else {
			break
		}
	}

	oldName := popElem.(string)
	newFileId := uuid.NewV4().String()
	newName := newFileId + ".csv"
	if err := os.Rename(filePath+oldName, filePath+newName); err != nil {
		panic(err)
	}

	fileMeta := trackMeta{
		FactoryID:     "factoryID_" + fmt.Sprint(generateRandom(1000)),
		StartTime:     int(time.Now().UnixMilli()) - generateRandom(300000),
		EndTime:       int(time.Now().UnixMilli()),
		RoadNumber:    "roadNumber_" + fmt.Sprint(generateRandom(100)),
		TraceFileName: newName,
		TraceID:       newFileId,
		ObjectClass:   TypeOfObjectClass[generateRandom(len(TypeOfObjectClass))],
		LicensePlate:  "Unkonw",
		TraceNumber:   generateRandom(100),
	}
	bytes, _ := json.Marshal(&fileMeta)
	return newName, string(bytes)

}

func generateRandom(max int) int {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(max)
	return r
}
