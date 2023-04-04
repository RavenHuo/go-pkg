/**
 * @Author raven
 * @Description
 * @Date 2021/9/15
 **/
package mongo

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type condition struct {
	mode       mgo.Mode
	context    *mongoContext
	model      IModel
	query      bson.M
	selector   interface{}
	sortFields []string
	skip       int
	limit      int
	retries    int
	startTime  *time.Time
	dbName     string
	colName    string
}

func (p *condition) Query(query bson.M) *condition {
	p.query = query
	return p
}

func (p *condition) Select(selector interface{}) *condition {
	p.selector = selector
	return p
}

func (p *condition) Sort(fields ...string) *condition {
	p.sortFields = fields
	return p
}

func (p *condition) Skip(n int) *condition {
	p.skip = n
	return p
}

func (p *condition) Limit(n int) *condition {
	p.limit = n
	return p
}

func (p *condition) Retry(n int) *condition {
	p.retries = n
	return p
}

func (p *condition) Mode(n mgo.Mode) *condition {
	p.mode = n
	return p
}

func (p *condition) List(result interface{}) (err error) {
	var q *mgo.Query
	q, err = p.buildQueryWithRetry()
	if err == nil {
		err = p.findAllWithRetry(q, result)
	}
	//日志
	if p.startTime != nil {
		p.context.shouldLog(err, "#mongo# find list from %s.%s with key `%s` in %.3fms, err:%s", p.dbName, p.colName,
			keyText(p.query), float32(time.Since(*p.startTime))/float32(time.Millisecond), errStr(err))
	}
	return
}

func (p *condition) FindOne(result interface{}) (err error) {
	var q *mgo.Query
	q, err = p.buildQueryWithRetry()
	if err == nil {
		err = p.findOneWithRetry(q, result)
	}
	//日志
	if p.startTime != nil {
		p.context.shouldLog(err, "#mongo# find one from %s.%s with key `%s` in %.3fms, err:%s", p.dbName, p.colName,
			keyText(p.query), float32(time.Since(*p.startTime))/float32(time.Millisecond), errStr(err))
	}
	return
}

func (p *condition) FindById(objectId bson.ObjectId, result interface{}) (err error) {
	p.query = bson.M{"_id": objectId}
	return p.FindOne(result)
}

func (p *condition) Count() (count int, err error) {
	var q *mgo.Query
	q, err = p.buildQueryWithRetry()
	if err == nil {
		retry := p.retries
		for retry > 0 {
			count, err = q.Count()
			if err != nil {
				if !shouldRetry(err) {
					break
				}
				retry--
				p.context.Logger.Errorf("count %T time-%d error: %v", p.model, p.retries-retry, err)
			} else {
				break
			}
		}
	}
	//日志
	if p.startTime != nil {
		p.context.shouldLog(err, "#mongo# count from %s.%s with key `%s` in %.3fms, err:%s", p.dbName, p.colName,
			keyText(p.query), float32(time.Since(*p.startTime))/float32(time.Millisecond), errStr(err))
	}
	return
}

// 迭代器 获取
func (p *condition) FindIter(batch int) (*mgo.Iter, error) {
	q, err := p.buildQueryWithRetry()
	if q == nil {
		return nil, err
	} else {
		return q.Batch(batch).Iter(), err
	}
}

func (p *condition) LoopQuery(obj interface{}, cb func()) {
	q, err := p.buildQueryWithRetry()
	if err == nil {
		iter := q.Iter()
		defer iter.Close()
		if iter.Err() != nil {
			for iter.Next(obj) {
				cb()
			}
		}
	}
}

func (p *condition) UpdateWithBson(update bson.M) error {

	start := time.Now()
	col, err := p.buildCollectionWithRetry()
	if err != nil {
		return err
	}
	err = col.Update(p.query, update)
	p.context.shouldLog(err, "#mongo# update %s.%s with key `%s` in %.3fms, err:%v", col.Database.Name,
		col.Name, keyText(p.query), float32(time.Since(start))/float32(time.Millisecond), errStr(err))
	return err
}

func (p *condition) Update(update IModel) error {

	start := time.Now()
	col, err := p.buildCollectionWithRetry()
	if err != nil {
		return err
	}
	err = col.Update(p.query, update)
	p.context.shouldLog(err, "#mongo# update %s.%s with key `%s` in %.3fms, err:%v", col.Database.Name,
		col.Name, keyText(p.query), float32(time.Since(start))/float32(time.Millisecond), errStr(err))
	return err
}

func (p *condition) Upsert(update IModel) (*mgo.ChangeInfo, error) {

	start := time.Now()
	col, err := p.buildCollectionWithRetry()
	if err != nil {
		return nil, err
	}
	ci, err := col.Upsert(p.query, update)

	updated := 0
	upserted := 0
	if ci != nil {
		updated = ci.Updated
	}
	if ci != nil && ci.UpsertedId != nil {
		upserted = 1
	}
	p.context.shouldLog(err, "#mongo# upsert %s.%s (update %d, upsert %d) in %.3fms, err:%s", col.Database.Name, col.Name,
		updated, upserted, float32(time.Since(start))/float32(time.Millisecond), errStr(err))

	return ci, err
}

func (p *condition) Remove(where bson.M) error {

	start := time.Now()
	col, err := p.buildCollectionWithRetry()
	if err != nil {
		return err
	}
	err = col.Remove(where)
	p.context.shouldLog(err, "#mongo# remove from %s.%s with key %s in %.3fms, err:%s", col.Database.Name, col.Name,
		keyText(where), float32(time.Since(start))/float32(time.Millisecond), errStr(err))
	return err
}

func (p *condition) RemoveAll(where bson.M) error {

	start := time.Now()
	col, err := p.buildCollectionWithRetry()
	if err != nil {
		return err
	}
	_, err = col.RemoveAll(where)
	p.context.shouldLog(err, "#mongo# remove all from %s.%s with key %s in %.3fms, err:%s", col.Database.Name, col.Name,
		keyText(where), float32(time.Since(start))/float32(time.Millisecond), errStr(err))
	return err
}

func (p *condition) buildQuery(q *mgo.Query) {
	if q == nil {
		return
	}
	if p.selector != nil {
		q = q.Select(p.selector)
	}
	if len(p.sortFields) > 0 {
		q = q.Sort(p.sortFields...)
	}
	if p.skip > 0 {
		q = q.Skip(p.skip)
	}
	if p.limit > 0 {
		q = q.Limit(p.limit)
	}
}

func (p *condition) buildQueryWithRetry() (q *mgo.Query, err error) {
	retry := p.retries
	for retry > 0 {
		q, err = p.context.Find(p.model, p.query)
		if err != nil {
			if !shouldRetry(err) {
				break
			}
			retry--
			p.context.Logger.Errorf("connect %T time-%d error: %v", p.model, p.retries-retry, err)
		} else {
			break
		}
	}
	p.buildQuery(q)
	return
}

func (p *condition) findAllWithRetry(q *mgo.Query, result interface{}) (err error) {
	retries := p.retries
	retry := p.retries
	for retry > 0 {
		err = q.All(result)
		if err != nil {
			if !shouldRetry(err) {
				break
			}
			retry--
			p.context.Logger.Errorf("query all %v time-%d error: %v", q, retries-retry, err)
		} else {
			break
		}
	}
	return
}

func (p *condition) findOneWithRetry(q *mgo.Query, result interface{}) (err error) {
	retries := p.retries
	retry := p.retries
	for retry > 0 {
		err = q.One(result)
		if err != nil {
			if !shouldRetry(err) {
				break
			}
			retry--
			p.context.Logger.Errorf("query one %v time-%d error: %v", q, retries-retry, err)
		} else {
			break
		}
	}
	return
}

func (p *condition) buildCollectionWithRetry() (c *mgo.Collection, err error) {
	retry := p.retries
	for retry > 0 {
		c, err = p.context.GetCollectionWithMode(p.model, p.mode)
		if err != nil {
			if !shouldRetry(err) {
				break
			}
			retry--
			p.context.Logger.Errorf("connect %T time-%d error: %v", p.model, p.retries-retry, err)
		} else {
			break
		}
	}
	return
}

func shouldRetry(err error) bool {
	// retry is done on the lower layer when ErrMaxConnectionLimited happens
	return err != nil && !IsErrNotFound(err) && !IsErrMaxConnectionLimited(err)
}
