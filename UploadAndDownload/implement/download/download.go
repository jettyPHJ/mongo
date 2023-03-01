package download

import (
	"UploadAndDownload/orm"
	r "UploadAndDownload/utils/gin"
	"UploadAndDownload/utils/log"
	"UploadAndDownload/utils/mgo"
	rds "UploadAndDownload/utils/redis"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var Coll *mongo.Collection

func init() {
	Coll = mgo.Database.Collection("trackFileMeta2")
	r.Router.GET("/trajectoryFile", get)
}

func get(c *gin.Context) {
	fmt.Println("路由到get函数时间戳:", time.Now())
	var req orm.GetTrajectoryFile
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		fmt.Println(err)
		c.JSON(601, gin.H{
			"msg": "缺少参数，请检查参数",
		})
		return
	}
	fmt.Println("请求参数：", req)
	result, err2 := FindFileDataByReq(req)
	if err2 != nil {
		c.JSON(602, gin.H{
			"msg": err2.Error(),
		})
		return
	}
	// fileName := "./temFiles/" + uuid.NewV4().String() + ".csv"
	// err3 := os.WriteFile(fileName, result, 0666)
	// if err3 != nil {
	// 	c.JSON(603, gin.H{
	// 		"msg": err3.Error(),
	// 	})
	// 	return
	// }
	if result == nil {
		c.JSON(200, gin.H{
			"msg": "未查到相关数据",
		})
		return
	}
	// c.File(fileName)
	// if err := os.Remove(fileName); err != nil {
	// 	log.ErrorLogger.Println(err)
	// }
	if _, err4 := c.Writer.Write(result); err4 != nil {
		log.ErrorLogger.Println(err4)
	}
}

func FindFileDataByReq(req orm.GetTrajectoryFile) ([]byte, error) {
	key := fmt.Sprint(req.StartTime) + "_" + fmt.Sprint(req.EndTime) + "_" + req.RoadNumber
	if req.ObjectClass != "" {
		key = key + "_" + req.ObjectClass
	}
	val, err := rds.Client.Get(key).Bytes()
	if err == nil {
		rds.Client.Expire(key, time.Minute*5) //更新redis中key时间
		return val, nil
	}
	//redis中没找到，到mongo中找
	_filter := []bson.M{}
	_filter = append(_filter, bson.M{"aggregationStartTime": bson.M{"$lte": req.EndTime}})
	_filter = append(_filter, bson.M{"aggregationEndTime": bson.M{"$gte": req.StartTime}})
	filter := bson.M{}
	filter["$and"] = _filter

	_filter2 := []bson.M{}
	_filter2 = append(_filter2, bson.M{"startTime": bson.M{"$lte": req.EndTime}})
	_filter2 = append(_filter2, bson.M{"endTime": bson.M{"$gte": req.StartTime}})
	_filter2 = append(_filter2, bson.M{"roadNumber": req.RoadNumber})
	if req.ObjectClass != "" {
		_filter2 = append(_filter2, bson.M{"objectClass": req.ObjectClass})
	}
	filter2 := bson.M{}
	filter2["$and"] = _filter2
	matchStage1 := bson.D{
		{Key: "$match", Value: filter},
	}
	unwindStage := bson.D{
		{Key: "$unwind", Value: "$data"},
	}
	groupStage1 := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "traceID", Value: "$data.traceID"},
			{Key: "factoryID", Value: "$data.factoryID"},
			{Key: "endTime", Value: "$data.endTime"},
			{Key: "roadNumber", Value: "$data.roadNumber"},
			{Key: "objectClass", Value: "$data.objectClass"},
			{Key: "licensePlate", Value: "$data.licensePlate"},
			{Key: "traceNumber", Value: "$data.traceNumber"},
			{Key: "startTime", Value: "$data.startTime"},
			{Key: "fileData", Value: "$data.fileData"},
		}},
	}
	groupStage2 := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "data", Value: 0},
			{Key: "_id", Value: 0},
		}},
	}
	matchStage2 := bson.D{
		{Key: "$match", Value: filter2},
	}
	pipeline := mongo.Pipeline{matchStage1, unwindStage, groupStage1, groupStage2, matchStage2}

	fmt.Println("开始执行管道命令时间戳:", time.Now())
	cursor, err2 := Coll.Aggregate(context.TODO(), pipeline)
	if err2 != nil {
		fmt.Println(err2)
		return nil, err2
	}
	defer cursor.Close(context.TODO())
	fmt.Println("执行管道命令结束时间戳:", time.Now())
	fmt.Println("查询结果游标长度：", cursor.RemainingBatchLength())
	//准备更新redis新的key-value
	// var results []bson.M
	// if err := cursor.All(context.TODO(), &results); err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(results)
	result := &bytes.Buffer{}
	csvWriter := csv.NewWriter(result)
	var wg sync.WaitGroup
	var Mutex sync.Mutex
	csvWriter.Write([]string{"traceID", "roadNumber", "time", "objectClass", "imageX", "imageY", "width", "height", "wordX", "wordY", "licensePlate"})
	for {
		if cursor.TryNext(context.TODO()) {
			wg.Add(1)
			go func(cursor mongo.Cursor) {
				defer func() {
					wg.Done()
					if err := recover(); err != nil {
						log.ErrorLogger.Panicln(err)
					}
				}()
				buff := &orm.TrackMeta{}
				if err := cursor.Decode(buff); err != nil {
					panic("error1:" + err.Error())
				}
				csvData, err := buff.TrackMetaToCsv()
				if err != nil {
					panic("error2:" + err.Error())
				}
				Mutex.Lock()
				if err := csvWriter.WriteAll(csvData); err != nil {
					Mutex.Unlock()
					panic("error3:" + err.Error())
				}
				Mutex.Unlock()
			}(*cursor)
			continue
		}
		if err := cursor.Err(); err != nil {
			fmt.Println(err)
		}
		if cursor.ID() == 0 {
			break
		}
	}
	wg.Wait()
	fmt.Println("解码和转csv完成时间戳:", time.Now())
	//已获取结果，更新redis，返回
	if result.Bytes() != nil {
		rds.Client.Set(key, result.Bytes(), time.Minute*5)
	}
	return result.Bytes(), nil
}

//  var wg sync.WaitGroup
// 	var Mutex sync.Mutex
// 	csvWriter.Write([]string{"traceID", "roadNumber", "time", "objectClass", "imageX", "imageY", "width", "height", "wordX", "wordY", "licensePlate"})
// 	for {
// 		if cursor.TryNext(context.TODO()) {
// 			wg.Add(1)
// 			go func(cursor *mongo.Cursor) {
// 				defer func() {
// 					wg.Done()
// 					if err := recover(); err != nil {
// 						log.ErrorLogger.Panicln(err)
// 					}
// 				}()
// 				buff := &orm.TrackMeta{}
// 				if err := cursor.Decode(buff); err != nil {
// 					panic("error1:" + err.Error())
// 				}
// 				csvData, err := buff.TrackMetaToCsv()
// 				if err != nil {
// 					panic("error2:" + err.Error())
// 				}
// 				Mutex.Lock()
// 				if err := csvWriter.WriteAll(csvData); err != nil {
// 					Mutex.Unlock()
// 					panic("error3:" + err.Error())
// 				}
// 				Mutex.Unlock()
// 			}(cursor)
// 			continue
// 		}
// 		if err := cursor.Err(); err != nil {
// 			fmt.Println(err)
// 		}
// 		if cursor.ID() == 0 {
// 			break
// 		}
// 	}
// 	wg.Wait()

// csvWriter.Write([]string{"traceID", "roadNumber", "time", "objectClass", "imageX", "imageY", "width", "height", "wordX", "wordY", "licensePlate"})
// 	buff := []orm.TrackMeta{}
// 	if err := cursor.All(context.TODO(), &buff); err != nil {
// 		log.ErrorLogger.Println(err)
// 	}
// 	for _, v := range buff {
// 		csvData, err := v.TrackMetaToCsv()
// 		if err != nil {
// 			log.ErrorLogger.Println(err)
// 		}
// 		if err := csvWriter.WriteAll(csvData); err != nil {
// 			log.ErrorLogger.Println(err)
// 		}
// 	}
