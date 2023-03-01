package orm

type GetTrajectoryFile struct {
	StartTime   int    `json:"startTime" desc:"开始时间" binding:"required"`
	EndTime     int    `json:"endTime" desc:"结束时间" binding:"required"`
	RoadNumber  string `json:"roadNumber" desc:"道路编码" binding:"required"`
	ObjectClass string `json:"objectClass" desc:"对象类型"`
}
