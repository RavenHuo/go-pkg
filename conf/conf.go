package conf

import (
	"fmt"
	"github.com/RavenHuo/go-kit/encode"
	"github.com/RavenHuo/go-kit/encode/ini"
	"github.com/RavenHuo/go-kit/encode/json"
	"github.com/RavenHuo/go-kit/encode/toml"
	"github.com/RavenHuo/go-kit/encode/xml"
	"github.com/RavenHuo/go-kit/encode/yaml"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

// Configuration
// 模仿 sync.Map overrider 存所有的数据
// cacheMap 是缓存，当查不到的时候才读取override
type Configuration struct {
	mu            *sync.RWMutex          // 读写锁，修改 override 的时候使用
	cacheMap      *sync.Map              // 用于缓存 xx.xx的数据，减少遍历 override map
	override      map[string]interface{} // key是字符串，value是 interface，也是map[string]interface{}
	keyDelimiter  string                 // key的分隔符，默认是.
	watchHandlers []WatchHandler         // 修改key的时候的处理方法
}

// Init init configuration
func Init(opts ...Opt) *Configuration {
	options := defaultOptions()
	for _, o := range opts {
		o(options)
	}

	conf := &Configuration{}
	dataMap := make(map[string]interface{})
	// 解析器
	encoder := parseEncoder(options.path)

	fileContent, err := ioutil.ReadFile(options.path)
	if err != nil {
		pwd, _ := os.Getwd()
		panic(fmt.Sprintf("read config file curPath:%s failed :%s", pwd, err))
	}
	err = encoder.Decode(fileContent, &dataMap)
	if err != nil {
		panic(fmt.Sprintf("encode failed, fileContent:%s, encoder:%s, err:%s", fileContent, encoder.Name(), err))
	}
	conf.override = dataMap
	conf.mu = &sync.RWMutex{}
	conf.cacheMap = &sync.Map{}
	conf.keyDelimiter = options.keyDelimiter
	conf.watchHandlers = options.watchHandlers
	return conf
}

func (conf *Configuration) Get(key string) interface{} {
	// 先查询缓存map
	val, ok := conf.cacheMap.Load(key)
	if ok {
		return val
	}

	// 不存在的时候
	conf.mu.RLock()
	defer conf.mu.RUnlock()
	paths := strings.Split(key, conf.keyDelimiter)
	lastKey := paths[len(paths)-1]
	pMap := searchMap(conf.override, paths[:len(paths)-1])
	val = pMap[lastKey]

	// 写入缓存
	conf.cacheMap.Store(key, val)
	return val
}

func (conf *Configuration) Set(key string, value interface{}) {
	// 加写锁
	conf.mu.Lock()
	defer conf.mu.Unlock()
	paths := strings.Split(key, conf.keyDelimiter)
	lastKey := paths[len(paths)-1]
	pMap := searchMap(conf.override, paths[:len(paths)-1])
	pMap[lastKey] = value

	// 写入缓存
	conf.cacheMap.Store(key, value)
}

func parseEncoder(path string) encode.Encoder {
	var encoder encode.Encoder
	if strings.HasSuffix(path, ".yaml") {
		encoder = yaml.NewEncoder()
	} else if strings.HasSuffix(path, ".json") {
		encoder = json.NewEncoder()
	} else if strings.HasSuffix(path, ".ini") {
		encoder = ini.NewEncoder()
	} else if strings.HasSuffix(path, ".toml") {
		encoder = toml.NewEncoder()
	} else if strings.HasSuffix(path, ".xml") {
		encoder = xml.NewEncoder()
	} else {
		panic("not support encoder")
	}
	return encoder
}

func searchMap(m map[string]interface{}, path []string) map[string]interface{} {
	for _, k := range path {
		m2, ok := m[k]
		// 不存在
		if !ok {
			return make(map[string]interface{})
		}
		// 将value 强转成  map[string]interface{}
		m3, ok := m2.(map[string]interface{})
		if !ok {
			m3 = make(map[string]interface{})
			m[k] = m3
		}
		m = m3
	}
	return m
}
