package infrastructure

import (
	rpc2 "app/internal/module/ipfs/rpc"
	"os"
)

// CreateFsClient 创建文件系统客户端
func CreateFsClient() (*rpc2.LApiStub, string, error) {
	strPPT, err := rpc2.GetLocalPassport(4080, 24)
	if err != nil {
		return nil, "", err
	}

	client := rpc2.InitLApiStubByUrl("127.0.0.1:4080")

	loginReply, err := client.LoginWithPPT(strPPT)
	if err != nil {
		return nil, "", err
	}

	return client, loginReply.Sid, nil
}

// ReadLocalFile 读取本地文件
func ReadLocalFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}
