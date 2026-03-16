package main

import (
	"flag"
	"fmt"
	"os"
)

func GetArgs() {
	args := os.Args
	// args[0] 是程序本身路径
	// args[1:] 才是实际参数

	fmt.Println(args)

	// go run main.go -port=8080 debug
	//	os.Args[0] = main
	//	os.Args[1] = -port=8080
	//	os.Args[2] = debug
}

func GetArgsByFlag() {
	port := flag.Int("port", 8080, "服务端口")
	env := flag.String("env", "dev", "运行环境")

	flag.Parse()

	fmt.Println(*port, *env)

	// go run main.go -port=9000 -env=prod
	// 	9000 prod

	//	解析完 flag 后获取“非 flag 参数”
	args := flag.Args() // 剩余参数
	// go run main.go -env=prod task1 task2
	//	flag.Args() == []string{"task1", "task2"}
	fmt.Println(args)
}
