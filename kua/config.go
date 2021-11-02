package kua

import "github.com/spf13/viper"

type ConfigInfo struct {
	Port string
	LogPath string
	LogLevel int
}

var ConfigStore ConfigInfo

func GetConfigYaml(filePath string){
	v := viper.New()
	
	//设置配置文件的文件名,文件类型,和路径
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(filePath)

	if err:= v.ReadInConfig();err != nil {
		WriteErrToLogFile(1,err)
	}
	
	parseYaml(v)

}

func parseYaml(v *viper.Viper){
	if err := v.Unmarshal(&ConfigStore) ; err != nil{
		WriteErrToLogFile(1,err)
	}
}
