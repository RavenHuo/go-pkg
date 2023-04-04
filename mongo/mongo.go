/**
 * @Author raven
 * @Description
 * @Date 2021/9/13
 **/
package mongo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime/debug"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoContext struct {
	clients      []*Client
	Timeout      time.Duration
	closed       bool
	noAutoClose  bool
	Mode         mgo.Mode
	Logger       *logrus.Logger
	printNoClose bool
	retryTime    int
	ctx          context.Context
}

func NewMongoContext(ctx context.Context) *mongoContext {
	return &mongoContext{ctx: ctx, Logger: logrus.New()}
}

func (mongo *mongoContext) WithModel(model IModel) *condition {
	queryOp := new(condition)
	now := time.Now()
	queryOp.startTime = &now
	queryOp.context = mongo
	queryOp.model = model
	queryOp.retries = 1
	queryOp.mode = mgo.Nearest

	modelConfig := model.(IModel).GetConfig()
	queryOp.colName = modelConfig.ColName
	queryOp.dbName = modelConfig.DBName
	return queryOp
}

func isWarnErr(err error) bool {
	return err == mgo.ErrNotFound || mgo.IsDup(err)
}

func errStr(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("%v", err)
}

func keyText(m bson.M) string {
	var buffer bytes.Buffer
	for k, v := range m {
		buffer.WriteString(k)
		buffer.WriteString("=")
		buffer.WriteString(fmt.Sprintf("%v", v))
	}
	return buffer.String()
}

func (mongo *mongoContext) shouldLog(err error, format string, params ...interface{}) {

	if err == nil {
		return
	}

	if isWarnErr(err) {
		mongo.Logger.Warnf(format, params)
	} else {
		mongo.Logger.Errorf(format, params)
	}

}

func (mongo *mongoContext) GetSession(dbName string, mode mgo.Mode) (*mgo.Session, error) {
	if mongo.clients == nil {
		mongo.clients = make([]*Client, 0, 10)
	}
	var client, err = clientPool.getClient(dbName, mode)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("error to get mongo client")
	}

	//检测没有释放的行为
	var stack []byte
	if clientPool.printNoClose {
		stack = debug.Stack()
	}
	go func() {
		timeout := mongo.Timeout
		if timeout == 0 {
			timeout = time.Second * 10
		}
		var ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()
		select {
		case <-client.closeCh:
		case <-ctx.Done():
			// only for detect
			mongo.Logger.Warn("mongo is not closed!")
			if stack != nil {
				mongo.Logger.Error(string(stack))
			}
		}
	}()

	mongo.clients = append(mongo.clients, client)
	return client.session, nil
}

func (mongo *mongoContext) Close() {
	// debug("Close")
	for _, client := range mongo.clients {
		client.Close()
	}
	mongo.clients = mongo.clients[0:0]
	mongo.closed = true
}

func (mongo *mongoContext) IsReadError(err error) bool {
	return err != nil && err != mgo.ErrNotFound
}

func (mongo *mongoContext) autoClose() {
	if !mongo.noAutoClose {
		mongo.Close()
	}
}

func (mongo *mongoContext) Find(obj interface{}, query bson.M) (*mgo.Query, error) {
	return mongo.FindWithMode(obj, query, mgo.Nearest)
}
func (mongo *mongoContext) FindWithMode(obj interface{}, query bson.M, mode mgo.Mode) (*mgo.Query, error) {
	start := time.Now()
	col, err := mongo.GetCollectionWithMode(obj.(IModel), mode)
	if err != nil {
		return nil, err
	}
	q := col.Find(query)
	mongo.Logger.Infof("#mongo# find from %s.%s with key `%s` in %.3fms", col.Database.Name, col.Name,
		keyText(query), float32(time.Since(start))/float32(time.Millisecond))

	return q, nil
}

func (mongo *mongoContext) GetCollectionWithMode(model IModel, mode mgo.Mode) (*mgo.Collection, error) {
	config := model.GetConfig()
	session, err := mongo.GetSession(config.DBName, mode)
	if err != nil {
		return nil, err
	}
	db := session.DB(config.DBName)
	col := db.C(config.ColName)
	return col, err
}
