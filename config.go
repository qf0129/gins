package ginz

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/sirupsen/logrus"
)

type Configuration struct {
	DbEngine   string
	SqliteFile string

	DbHost     string
	DbPort     uint
	DbUser     string
	DbPsd      string
	DbDatabase string

	PrimaryKey      string
	DefaultPageSize int

	AppHost    string
	AppPort    uint
	AppMode    string
	AppTimeout uint
	LogLevel   string

	// 加密算法密钥
	SecretKey string

	// 令牌过期时间，单位秒
	TokenExpiredTime uint

	Custom map[string]any
}

func (c *Configuration) Set(k string, v any) {
	c.Custom[k] = v
}

func (c *Configuration) Get(k string) any {
	return c.Custom[k]
}

func (c *Configuration) GetStr(k string) (r string) {
	if val := c.Custom[k]; val != nil {
		r, _ = val.(string)
	}
	return
}

func (c *Configuration) GetInt(k string) (r int) {
	if val := c.Custom[k]; val != nil {
		r, _ = val.(int)
	}
	return
}

func (c *Configuration) GetInt64(k string) (r int64) {
	if val := c.Custom[k]; val != nil {
		r, _ = val.(int64)
	}
	return
}

func (c *Configuration) GetUint(k string) (r uint) {
	if val := c.Custom[k]; val != nil {
		r, _ = val.(uint)
	}
	return
}

func (c *Configuration) GetBool(k string) (r bool) {
	if val := c.Custom[k]; val != nil {
		r, _ = val.(bool)
	}
	return
}

func (c *Configuration) GetFloat64(k string) (r float64) {
	if val := c.Custom[k]; val != nil {
		r, _ = val.(float64)
	}
	return
}

func (c *Configuration) Remove(k string) {
	delete(c.Custom, k)
}

var Config = &Configuration{
	DbEngine:   "sqlite",
	SqliteFile: "sqlite.db",

	DbHost:     "127.0.0.1",
	DbPort:     3306,
	DbUser:     "root",
	DbPsd:      "root",
	DbDatabase: "test",

	PrimaryKey:      "id",
	DefaultPageSize: 10,

	AppHost:    "",
	AppPort:    8080,
	AppMode:    "debug",
	AppTimeout: 60,
	LogLevel:   "debug",

	SecretKey:        "Abcd@123",
	TokenExpiredTime: 7200,

	Custom: make(map[string]any),
}

func LoadConfigFile() error {
	// 读取json文件
	data, err := os.ReadFile("config.json")
	if err != nil {
		logrus.Warn("Can not find config.json")
		return err
	}

	// json数据解析到配置
	err = json.Unmarshal(data, &Config)
	if err != nil {
		fmt.Println("json unmarshal failed, err:", err)
		return err
	}

	// json数据解析到map
	var m map[string]any
	err = json.Unmarshal(data, &m)
	if err != nil {
		fmt.Println("json unmarshal failed, err:", err)
		return err
	}

	// 加载json配置
	for k, v := range m {
		// 判断不包含的key，放到自定义配置里
		_, ok := reflect.TypeOf(*Config).FieldByName(k)
		if !ok {
			Config.Set(k, v)
		}
	}
	return nil
}
