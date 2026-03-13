package domain

import (
	"app/internal/rpc"
	"os"
	"strings"
)

// CreateFsClient 创建文件系统客户端
func CreateFsClient() (*rpc.LApiStub, string, error) {
	strPPT, err := rpc.GetLocalPassport(4080, 24)
	if err != nil {
		return nil, "", err
	}

	client := rpc.InitLApiStubByUrl("127.0.0.1:4080")

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

// IsDirNotExist 判断错误是否表示目录不存在
func IsDirNotExist(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "no link named") ||
		strings.Contains(err.Error(), "no linked named")
}
