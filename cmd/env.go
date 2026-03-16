package main

import (
	"fmt"
	"os"
)

func GetEnv() {
	//	检查环境变量是否存在
	env, ok := os.LookupEnv("APP_ENV")
	if !ok {
		// 不存在
		env = ""
	} else {
		// env 是具体值
		env = ""
	}

	fmt.Println("APP_ENV", env)
}

func GetAllEnv() {
	envs := os.Environ()
	for _, e := range envs {
		// e 的格式是 KEY=VALUE
		println(e)
	}
}

func OperationEnv() {
	//	设置环境变量
	err := os.Setenv("APP_ENV", "prod")
	if err != nil {
		return
	}

	//	删除环境变量
	err = os.Unsetenv("APP_ENV")
	if err != nil {
		return
	}
}
