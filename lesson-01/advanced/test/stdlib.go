package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"time"
)

func main() {
	// tickerDemo()
	// a := 97
	// v := string(a)
	// v2 := rune(a)
	// fmt.Println(v, v2)
	// hashDemo()
	hash2Demo()
}

func tickerDemo() {
	fmt.Println("=== tickerDemo ===")
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	for i := 0; i < 3; i++ {
		fmt.Println(<-ticker.C, "第", i, "次触发")
	}
	fmt.Println("Ticker 演示结束")
}

func hashDemo() {
	fmt.Println("=== xxx ===")
	data := "hello world"
	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(data))
	fmt.Printf("SHA1(\"%s\")=%x\n", data, sha1Hash.Sum(nil))

	sha256Hash := sha256.New()
	sha256Hash.Write([]byte(data))
	fmt.Printf("SHA256(\"%s\")=%x\n", data, sha256Hash.Sum(nil))

	sha512Hash := sha512.New()
	sha512Hash.Write([]byte(data))
	fmt.Printf("SHA512(\"%s\")=%x\n", data, sha512Hash.Sum(nil))
}

func hash2Demo() {
	fmt.Println("=== 单次写入、多次写入 ===")
	streamHash := sha256.New()
	for _, data := range []string{"hello", " ", "world"} {
		streamHash.Write([]byte(data))
	}
	fmt.Printf("streamHash.Sum(nil): %x\n", streamHash.Sum(nil))

	oneTimeHash := sha256.New()
	oneTimeHash.Write([]byte("hello world"))
	fmt.Printf("oneTimeHash.Sum(nil): %x\n", oneTimeHash.Sum(nil))
}
