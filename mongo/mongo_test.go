/**
 * @Author raven
 * @Description
 * @Date 2021/9/14
 **/
package mongo

import (
	context2 "context"
	"github.com/sirupsen/logrus"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id         bson.ObjectId `bson:"_id" json:"_id"`                 // 主键
	UserId     int           `bson:"user_id" json:"user_id"`         // 用户id
	Name       string        `bson:"name" json:"name"`               // 用户名
	Phone      int64         `bson:"phone" json:"phone"`             // 手机
	CreateDate *time.Time    `bson:"create_date" json:"create_date"` // 创建时间
	CreateTime int64         `bson:"create_time" json:"create_time"` // 创建时间戳
}

func (User) GetConfig() ModelConfig {
	return ModelConfig{
		DBName:  "test",
		ColName: "user",
	}
}

var userModel User

func TestNewMongoContext(t *testing.T) {
	mongoConfig := &Options{
		Retries:      2,
		PoolMax:      200,
		PoolSize:     100,
		PrintNoClose: true,
	}
	NewMongo(mongoConfig, logrus.New())

	context := NewMongoContext(context2.Background())
	defer context.Close()
	var user User
	err := context.WithModel(&userModel).Query(bson.M{"user_id": 1}).FindOne(&user)
	if err != nil {
		logrus.Infof("err=%s", err.Error())
	} else {
		logrus.Infof("user=%v", user)
	}

}
