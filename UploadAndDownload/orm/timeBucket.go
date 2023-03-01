package orm

type TimeBucket struct {
	AggregationStartTime int         `json:"aggregationstarttime" bson:"aggregationstarttime"`
	AggregationEndTime   int         `json:"aggregationendtime" bson:"aggregationendtime" `
	Count                int         `json:"count" bson:"count" `
	Size                 int         `json:"size" bson:"size" `
	Data                 []TrackMeta `json:"data" bson:"data" `
}
