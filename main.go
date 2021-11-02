package main

import (
	"fmt"
	"kua"
	"net/http"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	//从命令行获取配置文件的目录路径--cp ./config/
	pflag.String("cp","./config/","config.yaml path")
	viper.BindPFlags(pflag.CommandLine)

	//读取配置文件并反序列化到全局结构体变量
	kua.GetConfigYaml(viper.GetString("cp"))	

	fmt.Println(kua.ConfigStore.Port,kua.ConfigStore.LogPath,kua.ConfigStore.LogLevel)

	r := kua.New()

	r.GET("/verify", func(c *kua.Context) {
		c.JSON(http.StatusOK, kua.GenerateCaptchaHandler())
	})

	r.POST("/login", func(c *kua.Context) {
		mes := kua.H{"username": c.PostForm("username"), "password": c.PostForm("password")}
		ver := kua.V{"capId": c.PostForm("capId"), "bas64": c.PostForm("bas64")}
		c.JSON(http.StatusOK, kua.FindMongo(mes, ver))
	})

	r.POST("/register", func(c *kua.Context) {
		mes := kua.H{"username": c.PostForm("username"), "password": c.PostForm("password")}
		c.JSON(http.StatusOK, kua.InsertMongo(mes))
	})

	r.POST("/logout/:username", func(c *kua.Context) {
		user := kua.H{"username": c.Param("username"), "auth": c.PostForm("auth")}
		// buf, _ := json.Marshal(user)
		// fmt.Printf(string(buf))
		c.JSON(http.StatusOK, kua.DeleteMongo(user))
	})

	err := r.Run()
	if err != nil {
		kua.WriteErrToLogFile(4,err)
	}
}
