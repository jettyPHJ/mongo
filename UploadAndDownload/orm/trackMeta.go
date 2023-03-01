package orm

import (
	"UploadAndDownload/utils/log"
	"bytes"
	"encoding/csv"
)

type TrackMeta struct {
	FactoryID     string `json:"factoryID" bson:"factoryID" binding:"required"`
	StartTime     int    `json:"startTime"  bson:"startTime" binding:"required"`
	EndTime       int    `json:"endTime"  bson:"endTime" binding:"required"`
	RoadNumber    string `json:"roadNumber"  bson:"roadNumber" binding:"required"`
	TraceFileName string `json:"traceFileName"  bson:"traceFileName" binding:"required"`
	TraceID       string `json:"traceID"  bson:"traceID" binding:"required"`
	ObjectClass   string `json:"objectClass"  bson:"objectClass" binding:"required"`
	LicensePlate  string `json:"licensePlate"  bson:"licensePlate" binding:"required"`
	TraceNumber   int    `json:"traceNumber"  bson:"traceNumber" binding:"required"`
	FileData      []byte `json:"fileData"  bson:"fileData"`
}

func (meta *TrackMeta) TrackMetaToCsv() ([][]string, error) {
	defer func() {
		if recover() != nil {
			log.ErrorLogger.Panicln(recover())
		}
	}()
	buffreader := bytes.NewReader(meta.FileData)
	csvReader := csv.NewReader(buffreader)
	_csvData, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	re := [][]string{}
	for _, v := range _csvData {
		tem := []string{meta.TraceID, meta.RoadNumber, v[0], meta.ObjectClass, v[1], v[2], v[3], v[4], v[5], v[6], meta.LicensePlate}
		re = append(re, tem)
	}
	return re, nil
}
