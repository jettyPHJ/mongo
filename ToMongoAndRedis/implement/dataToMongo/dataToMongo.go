package datatomongo

import (
	"ToMongoAndRedis/modules/kfk_moudle"
	"ToMongoAndRedis/orm"
	"ToMongoAndRedis/utils/log"
	m "ToMongoAndRedis/utils/mgo"
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"go.mongodb.org/mongo-driver/mongo"
)

var collectionName string = "trackFileMeta2"
var Coll *mongo.Collection

func init() {
	Coll = m.Database.Collection(collectionName)
	partitions := kfk_moudle.NewPartitions("receive_trackfile")
	go partitions.ConsumeAllParts(act, 3)
}

func act(partConsumer sarama.PartitionConsumer) {
	var Mutex sync.Mutex
	bucketTrackMeta := &orm.TimeBucket{}
	toDog := make(chan *orm.TimeBucket, 10)
	go dog(toDog, &Mutex)
	for m := range partConsumer.Messages() {
		temTrackMeta := orm.TrackMeta{}
		if err := json.Unmarshal(m.Value, &temTrackMeta); err != nil {
			log.ErrorLogger.Println(err)
		}
		if bucketTrackMeta.Size+len(temTrackMeta.FileData) < 16000000 {
			Mutex.Lock()
			bucketTrackMeta.AddTraMeta(temTrackMeta)
			Mutex.Unlock()
			toDog <- bucketTrackMeta //喂狗
			if bucketTrackMeta.Size >= 14000000 {
				go timeBucketToMgo(*bucketTrackMeta)
				bucketTrackMeta = &orm.TimeBucket{}
			}
		} else { //最后一个文件过大，应放入下一个桶
			Mutex.Lock()
			go timeBucketToMgo(*bucketTrackMeta)
			bucketTrackMeta = &orm.TimeBucket{}
			bucketTrackMeta.AddTraMeta(temTrackMeta)
			Mutex.Unlock()
			toDog <- bucketTrackMeta //喂狗
		}
	}
}

func dog(bucket <-chan *orm.TimeBucket, Mutex *sync.Mutex) {
	timeCount := 0
	temBucket := &orm.TimeBucket{}
	for {
		if timeCount == 5 && temBucket.Count != 0 {
			Mutex.Lock()
			timeBucketToMgo(*temBucket)
			timeCount = 0
			*temBucket = orm.TimeBucket{}
			Mutex.Unlock()
		}
		select {
		case <-time.After(time.Second):
			timeCount++
			timeCount = timeCount % 6
		case tem := <-bucket:
			timeCount = 0
			temBucket = tem
		}
	}
}

// func timeBucketToMgo(bucket orm.TimeBucket) {
// 	filter := bson.M{"size": bson.M{"$lt": 12000000}}
// 	opts := options.Update().SetUpsert(true)
// 	update := bson.M{}
// 	update["$inc"] = bson.M{"count": bucket.Count, "size": bucket.Size}
// 	update["$push"] = bson.M{"data": bson.M{"$each": bucket.Data}}
// 	update["$min"] = bson.M{"aggregationstarttime": bucket.AggregationStartTime}
// 	update["$max"] = bson.M{"aggregationendtime": bucket.AggregationEndTime}
// 	c, _ := context.WithTimeout(context.Background(), time.Second*5)
// 	_, err := Coll.UpdateOne(c, filter, update, opts)
// 	if err != nil {
// 		log.ErrorLogger.Println(err)
// 		return
// 	}
// }
func timeBucketToMgo(bucket orm.TimeBucket) {
	c, _ := context.WithTimeout(context.Background(), time.Second*3)
	_, err := Coll.InsertOne(c, bucket)
	if err != nil {
		log.ErrorLogger.Println(err)
		return
	}
}
