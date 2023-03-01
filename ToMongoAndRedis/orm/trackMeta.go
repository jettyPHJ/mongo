package orm

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

// type FileData struct {
// 	TimeStamp string
// 	ImageX    string
// 	ImageY    string
// 	Width     string
// 	Height    string
// 	WorldX    string
// 	WorldY    string
// }
