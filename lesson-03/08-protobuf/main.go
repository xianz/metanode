package main

/**
 */

import (
	"fmt"
	"log"
	"os"
	pbtypes "protobuf/generated"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/encoding/proto"
)

func main() {
	r := gin.Default()

	// ProtoBuf 路由
	api := r.Group("proto")
	{
		api.GET("/save", SaveProto)
		api.GET("/load", LoadProto)
	}

	r.Run(":8080")
}

var protoFilePath = "person.pb"

func SaveProto(c *gin.Context) {
	pb := &pbtypes.Person{
		Name:  "John",
		Id:    1,
		Email: "john@example.com",
	}
	data, err := proto.Marshal(pb)
	if err != nil {
		log.Fatal("marshal error:", err)
	}
	if err := os.WriteFile(protoFilePath, data, 0644); err != nil {
		log.Fatal("write file error:", err)
	}
	fmt.Printf("序列化后的二进制长度: %d\n", len(data))
	c.String(200, string(data))
}

func LoadProto(c *gin.Context) {
	data, err := os.ReadFile(protoFilePath)
	if err != nil {
		log.Fatal("read file error:", err)
	}
	pb := &pbtypes.Person{}
	if err := proto.Unmarshal(data, pb); err != nil {
		log.Fatal("unmarshal error:", err)
	}
	fmt.Printf("反序列化后的数据: %#v\n", pb)
	c.String(200, fmt.Sprintf("%v", pb))
}
