package main

import (
	"app/rpc"
	"fmt"
)

func main() {
	//以下是不需要认证身份的示例
	//节点的地址端口是127.0.0.1:4800
	stub := rpc.InitLApiStubByUrl("127.0.0.1:4800")

	var ver string
	err := stub.GetVarObj(&ver, "", rpc.ApiVarVersion)
	if err != nil {
		panic(err)
	}

	fmt.Println("cur ver:", ver)
}
