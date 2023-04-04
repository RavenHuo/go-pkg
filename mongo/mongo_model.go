/**
 * @Author raven
 * @Description
 * @Date 2021/9/14
 **/
package mongo

// model 接口
type ModelConfig struct {
	DBName  string
	ColName string
}

type IModel interface {
	GetConfig() ModelConfig
}
