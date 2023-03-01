package orm

type TimeBucket struct {
	AggregationStartTime int         `bson:"aggregationStartTime"`
	AggregationEndTime   int         `bson:"aggregationEndTime"`
	Count                int         `bson:"count"`
	Size                 int         `bson:"size"`
	Data                 []TrackMeta `bson:"data"`
}

func (bucket *TimeBucket) AddTraMeta(data TrackMeta) {
	if bucket.Count == 0 {
		bucket.AggregationStartTime = data.StartTime
		bucket.AggregationEndTime = data.EndTime
		bucket.Count++
		bucket.Size = 550 + len(data.FileData)
		bucket.Data = append(bucket.Data, data)
		return
	}
	if data.StartTime < bucket.AggregationStartTime {
		bucket.AggregationStartTime = data.StartTime
	}
	if data.EndTime > bucket.AggregationEndTime {
		bucket.AggregationEndTime = data.EndTime
	}
	bucket.Size = bucket.Size + 500 + len(data.FileData)
	bucket.Count++
	bucket.Data = append(bucket.Data, data)
}
