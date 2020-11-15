package internal                                           

import (                                               
    "fmt"
    "gopkg.in/yaml.v2"
    "io/ioutil"
)

type StrategyOpt struct {
	DeviceCycle *int `yaml:"DeviceCycle"`
	ControlCycle *int `yaml:"ControlCycle"`
}
type SocketOpt struct{
	Port *int `yaml:"Port"`
}

type HttpOpt struct{
	Port *int `yaml:"Port"`
	BasicAuth *BasicAuthOpt `yaml:"BasicAuth"`
}

type BasicAuthOpt struct{
	User *string `yaml:"User"`
	Password *string `yaml:"Password"`
}

//Yamlとして読み取る構造体の全体像
type YamlConfig struct{
	LogPath *string `yaml:"LogPath"`
	SocketServer *SocketOpt `yaml:"SocketServer"`
	HttpServer *HttpOpt `yaml:"HttpServer"`
	Strategy *StrategyOpt `yaml:"Strategy"`
}

//Yamlをパースして構造体を返す
//本来ならば参照を返すと効率的だが、起動時の一回しか呼ばれないので可読性の高い値返しにしている
func ParseYaml(path string) (YamlConfig, error) {
	data,err := ioutil.ReadFile(path)
	var conf YamlConfig
    if err!=nil{
		return conf,fmt.Errorf("ParseYaml() - Problem to open the file:",err.Error())
	}
	//厳密に構造体を読み取れる
	//定義されていないKeyを指定することはできない
	//Keyを省略することはできてしまうので後々検証が必要
	err = yaml.UnmarshalStrict(data,&conf)
	if err!=nil{
		return conf,fmt.Errorf("ParseYaml() - Problem to interpret yaml:",err.Error())
	}
	//fmt.Printf("(ParseYaml)>>%#v\n",conf)
	Verify(&conf)
	return conf,nil
}

func Verify(conf *YamlConfig){
	if conf.LogPath == nil {
		conf.LogPath = new(string)
		*conf.LogPath = "/tmp/iot_gateway.log"
	}
	if conf.SocketServer == nil{
		conf.SocketServer = new(SocketOpt)
	}
	if conf.SocketServer.Port == nil{
		conf.SocketServer.Port = new(int)
		*conf.SocketServer.Port = 8081
	}
	if conf.HttpServer == nil{
		conf.HttpServer = new(HttpOpt)
	}
	if conf.HttpServer.Port == nil{
		conf.HttpServer.Port = new(int)
		*conf.HttpServer.Port = 8080
	}
	if conf.HttpServer.BasicAuth == nil{
		conf.HttpServer.BasicAuth = new(BasicAuthOpt)
	}
	if conf.HttpServer.BasicAuth.User == nil{
		conf.HttpServer.BasicAuth.User = new(string)
		*conf.HttpServer.BasicAuth.User = "someone"
	}
	if conf.HttpServer.BasicAuth.Password == nil{
		conf.HttpServer.BasicAuth.Password = new(string)
		*conf.HttpServer.BasicAuth.Password = "somepassword"
	}
	if conf.Strategy == nil{
		conf.Strategy = new(StrategyOpt)
	}
	if conf.Strategy.DeviceCycle == nil{
		conf.Strategy.DeviceCycle = new(int)
		*conf.Strategy.DeviceCycle = 10
	}
	if conf.Strategy.ControlCycle == nil{
		conf.Strategy.ControlCycle = new(int)
		*conf.Strategy.ControlCycle = 90
	}
}
