package upload

import (
	"UploadAndDownload/orm"
	"UploadAndDownload/utils/config"
	r "UploadAndDownload/utils/gin"
	"UploadAndDownload/utils/log"
	"bytes"
	"io"

	k "UploadAndDownload/utils/kafka"

	"encoding/json"
	"os"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
)

func init() {
	os.MkdirAll(config.Param.TrackfileSavepath, os.ModePerm)
	r.Router.POST("/trajectoryFile", post)
}

func post(c *gin.Context) {
	reader, err := c.Request.MultipartReader()
	if err != nil {
		log.ErrorLogger.Println(err)
		return
	}
	tem := &orm.TrackMeta{}
	//保存元数据和轨迹文件
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if part.FileName() == "" {
			//该部分part为视屏元数据
			b := &bytes.Buffer{}
			io.Copy(b, part)
			err = json.Unmarshal(b.Bytes(), tem)
			if err != nil {
				log.ErrorLogger.Println(err)
				return
			}
			if part.Close() != nil {
				log.ErrorLogger.Println(err)
			}
			continue
		}
		//该部分为文件
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, part)
		if err != nil {
			log.ErrorLogger.Println(err)
			return
		}
		tem.FileData = buf.Bytes()
		if part.Close() != nil {
			log.ErrorLogger.Println(err)
		}
	}
	//fmt.Println(tem)
	//生成kafka消息
	msg := NewKfkMessage(tem, "receive_trackfile")
	//发布kafka消息
	k.SyncProducer.SendMessage(msg)
	c.JSON(200, gin.H{
		"msg": "操作成功",
	})
}

func NewKfkMessage(v interface{}, topic string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	jsonByte, _ := json.Marshal(v)
	msg.Value = sarama.StringEncoder(string(jsonByte))
	return msg
}
