package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config 全局配置信息变量（指针）
var Config = new(AppConfig)

type AppConfig struct {
	// 应用基础配置
	Host         string `mapstructure:"host"`
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	Port         int    `mapstructure:"port"`
	StartTime    string `mapstructure:"start_time"`
	MachineID    int64  `mapstructure:"machine_id"`
	*LogConfig   `mapstructure:"log"`
	*MysqlConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	// 日志配置
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MysqlConfig struct {
	// MySQL 配置
	Host           string `mapstructure:"host"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	DbName         string `mapstructure:"db_name"`
	Port           int    `mapstructure:"port"`
	MaxOpenConnect int    `mapstructure:"max_open_connect"`
	MaxIdleConnect int    `mapstructure:"max_idle_connect"`
}

type RedisConfig struct {
	// Redis 配置
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DbName   int    `mapstructure:"db_name"`
	PoolSize int    `mapstructure:"pool_size"`
}

// Init 读取配置文件并加载到全局 Config。
func Init(filename string) error {

	//方式1:直接指定配置文件的名称及后缀 和文件路径 （相对路径或绝对路径）
	//viper.SetConfigFile("./conf/config.yaml") //指定配置文件名和路径 相对路径

	//方式2:指定配置文件的文件名称 和配置文件的位置 ，文件名不需要带后缀，配置文件路径可配置多个 viper可以自动查找
	//viper.SetConfigName("config") // 配置文件名称 不用带后缀
	//viper.AddConfigPath(".")      //配置文件路径 可以填多个
	//viper.AddConfigPath("./conf") //配置文件路径 可以填多个

	//viper.SetConfigType("yaml")   //扩展名  专用于从远程获取配置信息时指定配置文件的后缀名  不是远程可以不填

	//方式3:通过传进来的文件名称路径参数来配置
	viper.SetConfigFile(filename) //指定配置文件名和路径 相对路径

	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {             // 处理读取配置文件的错误
		fmt.Printf("viper.ReadInConfig() err: %v\n", err)
		return err
	}
	//把读取到的配置信息 反序列化存到Config
	if err := viper.Unmarshal(Config); err != nil {
		fmt.Printf("viper.Unmarshal() err: %v\n", err)
	}
	viper.WatchConfig() //自动检测配置文件变化
	//配置变化后的回调函数
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件修改了：", e.Name)
		if err := viper.Unmarshal(Config); err != nil {
			fmt.Printf("viper.Unmarshal() err: %v\n", err)
		}
	})
	return nil
}
