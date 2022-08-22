/**
 * @Author raven
 * @Description
 * @Date 2022/7/28
 **/
package redigo

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

var ctx = context.Background()
var pool, _ = getRedisPool()

func TestNewPool(t *testing.T) {
	option := &Options{
		Addr: "127.0.0.1:6379",
	}
	pool, err := NewPool(option)
	if err != nil {
		fmt.Printf("NewPool err:%s \n", err)
		return
	}
	defer pool.ClosePool()
	redisConn, err := pool.GetConn()
	if err != nil {
		fmt.Printf("GetConn err:%s \n", err)
		return
	}
	defer redisConn.Close()
}

func getRedisPool() (*Pool, error) {
	option := &Options{
		Addr: "127.0.0.1:6379",
	}
	pool, err := NewPool(option)
	if err != nil {
		fmt.Printf("NewPool err:%s \n", err)
		return nil, err
	}
	return pool, nil
}

func TestPing(t *testing.T) {
	pool, err := getRedisPool()
	if err != nil {
		fmt.Printf("NewPool err:%s \n", err)
		return
	}
	defer pool.ClosePool()
	redisConn, err := pool.GetConn()
	if err != nil {
		fmt.Printf("GetConn err:%s \n", err)
		return
	}
	defer redisConn.Close()
	cmder := redisConn.Ping(ctx)
	fmt.Println("ping success ------------")
	fmt.Printf("cmder val:%s err:%s \n", cmder.Val(), cmder.Err())
}

func TestGetSet(t *testing.T) {
	pool, err := getRedisPool()
	if err != nil {
		fmt.Printf("NewPool err:%s \n", err)
		return
	}
	defer pool.ClosePool()
	redisConn, err := pool.GetConn()
	if err != nil {
		fmt.Printf("GetConn err:%s \n", err)
		return
	}
	defer redisConn.Close()
	setCmder := redisConn.Set(ctx, "name", "daenerys")
	if setCmder.Err() != nil {
		fmt.Printf("set err :%s \n", setCmder.Err())
		return
	}
	fmt.Println("set success ------------")
	fmt.Printf("setCmder val:%s err:%s \n", setCmder.Val(), setCmder.Err())

	getCmder := redisConn.Get(ctx, "name")
	if getCmder.Err() != nil {
		fmt.Printf("get err :%s \n", getCmder.Err())
		return
	}
	fmt.Println("get success ------------")
	fmt.Printf("getCmder val:%s err:%s \n", getCmder.Val(), getCmder.Err())

	if getCmder.Val() != "daenerys" {
		fmt.Println("redisConn get not equal set ")
	}
}

func TestDel(t *testing.T) {
	redisConn, err := pool.GetConn()
	if err != nil {
		fmt.Printf("GetConn err:%s \n", err)
		return
	}
	defer redisConn.Close()
	setCmder := redisConn.Set(ctx, "name", "hello")
	if setCmder.Err() != nil {
		fmt.Printf("set err :%s \n", setCmder.Err())
		return
	}
	fmt.Println("set success ------------")
	fmt.Printf("setCmder val:%s err:%s \n", setCmder.Val(), setCmder.Err())

	delCmder := redisConn.Del(ctx, "name")
	if delCmder.Err() != nil {
		fmt.Printf("del err :%s \n", delCmder.Err())
		return
	}
	fmt.Println("del success ------------")
	fmt.Printf("delCmder val:%+v err:%s \n", delCmder.Val(), delCmder.Err())

	delNotExistCmder := redisConn.Del(ctx, "daenerys")
	if delCmder.Err() != nil {
		fmt.Printf("del err :%s \n", delNotExistCmder.Err())
		return
	}

	fmt.Println("delNotExist success ------------")
	fmt.Printf("delNotExistCmder val:%+v err:%s \n", delNotExistCmder.Val(), delNotExistCmder.Err())
}

func TestExist(t *testing.T) {
	redisConn, err := pool.GetConn()
	if err != nil {
		fmt.Printf("GetConn err:%s \n", err)
		return
	}
	defer redisConn.Close()
	setCmder := redisConn.Set(ctx, "name", "hello")
	if setCmder.Err() != nil {
		fmt.Printf("set err :%s \n", setCmder.Err())
		return
	}
	fmt.Println("set success ------------")
	fmt.Printf("setCmder val:%s err:%s \n", setCmder.Val(), setCmder.Err())

	existCmder := redisConn.Exists(ctx,"name")
	if existCmder.Err() != nil {
		fmt.Printf("exist err :%s \n", existCmder.Err())
		return
	}
	fmt.Printf("exist success val:%+v \n", existCmder.Val())

	delCmder := redisConn.Del(ctx, "name")
	if delCmder.Err() != nil {
		fmt.Printf("del err :%s \n", delCmder.Err())
		return
	}
	fmt.Println("del success ------------")
	fmt.Printf("delCmder val:%+v err:%s \n", delCmder.Val(), delCmder.Err())

	notExistCmder := redisConn.Exists(ctx,	"name")
	if notExistCmder.Err() != nil {
		fmt.Printf("exist err :%s \n", notExistCmder.Err())
		return
	}
	fmt.Printf("exist success val:%+v \n", notExistCmder.Val())

	if notExistCmder.Val() != 0 {
		t.Fatal("exist err")
	}

}
func TestTtl(t *testing.T) {
	redisConn, err := pool.GetConn()
	if err != nil {
		fmt.Printf("GetConn err:%s \n", err)
		return
	}
	defer redisConn.Close()
	setCmder := redisConn.SetEX(ctx, "name", "hello", 10)
	if setCmder.Err() != nil {
		t.Fatalf("setnx err :%s \n", setCmder.Err())
		return
	}
	fmt.Println("set success ------------")
	fmt.Printf("setCmder val:%+v err:%s \n", setCmder.Val(), setCmder.Err())

	durationCmder := redisConn.TTL(ctx,"name")
	if durationCmder.Err() != nil {
		t.Fatalf("ttl err:%s", durationCmder.Err())
		return	
	}
	fmt.Printf("ttl success val:%+v", durationCmder.Val())

}

func TestRefer(t *testing.T) {
	for i := 0; i < 10; i++ {
		go printRefer()
	}
	printRefer()
	time.Sleep(10 * time.Second)
}

func printRefer() {
	var skip = 1
	var pc, _, _, ok = runtime.Caller(skip)
	if !ok {
		fmt.Println("not ok")
		return
	}
	var f = runtime.FuncForPC(pc)
	var fname = f.Name()
	fmt.Println(fname)
}
