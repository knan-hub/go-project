package setting

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
	Port      int    `mapstructure:"port"`

	*Self           `mapstructure:"self"`
	*LogConfig      `mapstructure:"log"`
	*VectorDbConfig `mapstructure:"verctorDb"`
	*RedisConfig    `mapstructure:"redis"`
}

type Self struct {
	INTERNAL_API_KEY string `mapstructure:"INTERNAL_API_KEY"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type VectorDbConfig struct {
	Url                         string `mapstructure:"url"`
	Key                         string `mapstructure:"key"`
	Username                    string `mapstructure:"username"`
	DatabaseName                string `mapstructure:"database_name"`
	SiteCollectionView          string `mapstructure:"site_collection_view"`
	KnowledgeBaseCollectionView string `mapstructure:"knowledge_base_collection_view"`
}

type RedisConfig struct {
	Host        string        `mapstructure:"host"`
	Port        string        `mapstructure:"port"`
	Password    string        `mapstructure:"password"`
	DefaultDb   int           `mapstructure:"defaultDb"`
	DialTimeout time.Duration `mapstructure:"dialTimeout"`
}

var Config = new(AppConfig)

func Init(filePath string) (err error) {
	// 方式1：直接指定配置文件路径（相对路径或者绝对路径）
	// 相对路径：相对执行的可执行文件的相对路径
	//viper.SetConfigFile("./conf/config.yaml")
	// 绝对路径：系统中实际的文件路径
	//viper.SetConfigFile("/Users/liwenzhou/Desktop/bluebell/conf/config.yaml")

	// 方式2：指定配置文件名和配置文件的位置，viper自行查找可用的配置文件
	// 配置文件名不需要带后缀
	// 配置文件位置可配置多个
	//viper.SetConfigName("config") // 指定配置文件名（不带后缀）
	//viper.AddConfigPath(".") // 指定查找配置文件的路径（这里使用相对路径）
	// viper.AddConfigPath("./conf")      // 指定查找配置文件的路径（这里使用相对路径）

	// 基本上是配合远程配置中心使用的，告诉viper当前的数据使用什么格式去解析
	//viper.SetConfigType("json")

	viper.SetConfigFile(filePath)

	// 读取配置信息
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		return
	}

	// 把读取到的配置信息反序列化到 Conf 变量中
	if err := viper.Unmarshal(Config); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	}

	viper.WatchConfig()

	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
		if err := viper.Unmarshal(Config); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		}
	})

	return
}
