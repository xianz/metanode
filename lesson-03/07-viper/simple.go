package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
//	"log"
//	"net/http"
)

func main(){
viper.SetConfigName("config")
viper.SetConfigType("yaml")
viper.AddConfigPath(".")
err := viper.ReadInConfig()
if err != nil {
	panic("读取configuration失败")
}

r := gin.Default()
r.GET("/", func(c *gin.Context){
	port := viper.GetString("server.port")
	mode := viper.GetString("server.mode")
	c.JSON(200, gin.H{"message":"ok", "port":port, "mode":mode})
})
r.Run(":8080")
}
