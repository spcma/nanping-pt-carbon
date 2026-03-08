package main

import (
	"app/rpc"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

// ==================== 全局变量 ====================

var npFsClient *rpc.LApiStub
var npSessionId string

// ==================== 数据结构 ====================

// NpFsParam NPFS 操作参数结构
type NpFsParam struct {
	Path     string `json:"path" form:"path"`
	Filename string `json:"filename" form:"filename"`
	Content  string `json:"content" form:"content"`
}

// TempFileHandle 临时文件句柄
type TempFileHandle struct {
	fsid string
}

// ==================== 初始化相关 ====================

// InitNpFsService 初始化 NPFS 文件系统服务
// 获取本地通行证并登录，初始化全局客户端和会话 ID
func InitNpFsService() {
	client, curSid, err := createFsClient()
	if err != nil {
		panic(err)
	}

	npFsClient = client
	npSessionId = curSid
}

// CloseNpFsService 关闭 NPFS 文件系统服务连接
func CloseNpFsService() {
	if npFsClient != nil {
		npFsClient.Logout(npSessionId, "")
	}
}

// createFsClient 创建 FS 客户端连接
// 返回客户端实例、会话 ID 和错误信息
func createFsClient() (*rpc.LApiStub, string, error) {
	strPPT, err := rpc.GetLocalPassport(4800, 24)
	if err != nil {
		return nil, "", err
	}

	client := rpc.InitLApiStubByUrl("127.0.0.1:4800")

	loginReply, err := client.LoginWithPPT(strPPT)
	if err != nil {
		return nil, "", err
	}

	return client, loginReply.Sid, nil
}

// ==================== 目录操作相关 ====================

// CheckDirExists 检查 NPFS 目录是否存在
// path: 目录路径
// exists: 是否存在，err: 错误信息
func CheckDirExists(path string) (bool, error) {
	_, err := npFsClient.FilesStat(npSessionId, path)
	if err != nil {
		if isDirNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// CreateDirIfNotExists 如果目录不存在则创建
// path: 目录路径
// recursive: 是否递归创建父目录
// created: 是否执行了创建操作，err: 错误信息
func CreateDirIfNotExists(path string, recursive bool) (bool, error) {
	exists, err := CheckDirExists(path)
	if err != nil {
		log.Printf("err: %v", err)
		return false, err
	}

	if !exists {
		err = npFsClient.FilesMkdir(npSessionId, path, recursive)
		if err != nil {
			log.Printf("err: %v", err)
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// EnsureDirExists 确保目录存在，不存在则创建
// path: 目录路径
// recursive: 是否递归创建父目录
func EnsureDirExists(path string, recursive bool) error {
	_, err := CreateDirIfNotExists(path, recursive)
	return err
}

// isDirNotExist 判断错误是否表示目录不存在
func isDirNotExist(err error) bool {
	if err == nil {
		return false
	}
	return !strings.Contains(err.Error(), "no linked named")
}

// ListDirectory 列出目录内容
// path: 目录路径
// links: 目录中的链接列表，err: 错误信息
func ListDirectory(path string) ([]rpc.LsLink, error) {
	fmt.Println("ListDirectory", path, "npSessionId", npSessionId)
	links, err := npFsClient.FilesLs(npSessionId, path)
	if err != nil {
		return nil, err
	}
	fmt.Println("links", links)
	return links, nil
}

// DeleteFile 删除 NPFS 文件
// path: 文件路径
// recursive: 是否递归删除
// force: 是否强制删除
func DeleteFile(path string, recursive, force bool) error {
	err := npFsClient.FilesRm(npSessionId, path, recursive, force)
	if err != nil {
		return err
	}
	return nil
}

// ==================== 文件读取相关 ====================

// ReadFileFromNpfs 从 NPFS 读取文件数据
// filePath: NPFS 文件路径（如：/np_storage/1.jpg）
// data: 文件数据，size: 文件大小，err: 错误信息
func ReadFileFromNpfs(filePath string) ([]byte, int64, error) {
	// 打开文件 URL
	fsid, err := npFsClient.MMOpenUrl(npSessionId, filePath)
	if err != nil {
		return nil, 0, err
	}
	defer npFsClient.MMClose(fsid)

	// 获取文件大小
	size, err := npFsClient.MFGetSize(fsid)
	if err != nil {
		return nil, 0, err
	}

	// 读取文件数据
	data, err := npFsClient.MFGetData(fsid, 0, int(size))
	if err != nil {
		return nil, 0, err
	}

	return data, size, nil
}

// SaveFileToLocal 将 NPFS 文件保存到本地
// filePath: NPFS 文件路径
// localPath: 本地保存路径
func SaveFileToLocal(filePath, localPath string) error {
	data, _, err := ReadFileFromNpfs(filePath)
	if err != nil {
		return err
	}

	err = os.WriteFile(localPath, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// GetFileStream 获取文件流用于 HTTP 响应
// filePath: NPFS 文件路径
// filename: 下载文件名
// data: 文件数据，err: 错误信息
func GetFileStream(filePath, filename string) ([]byte, error) {
	fsid, err := npFsClient.MMOpenUrl(npSessionId, filePath)
	if err != nil {
		return nil, err
	}
	defer npFsClient.MMClose(fsid)

	bytess, err := npFsClient.MFGetData(fsid, 0, -1)
	if err != nil {
		return nil, err
	}

	return bytess, nil
}

// ==================== 临时文件操作相关 ====================

// OpenTempFile 打开 NPFS 临时文件
// fsid: 临时文件 ID，err: 错误信息
func OpenTempFile() (string, error) {
	fsid, err := npFsClient.MFOpenTempFile(npSessionId)
	if err != nil {
		return "", err
	}
	return fsid, nil
}

// WriteDataToTempFile 向临时文件写入数据
// fsid: 临时文件 ID
// data: 要写入的数据
// offset: 写入偏移量
// written: 写入的字节数，err: 错误信息
func WriteDataToTempFile(fsid string, data []byte, offset int) (int, error) {
	written, err := npFsClient.MFSetData(fsid, data, int64(offset))
	if err != nil {
		return 0, err
	}
	return written, nil
}

// SaveTempFileToNpfs 将临时文件保存到 NPFS 指定路径
// fsid: 临时文件 ID
// destPath: 目标路径（如：path/filename.txt）
// ipfsid: IPFS 文件 ID，err: 错误信息
func SaveTempFileToNpfs(fsid, destPath string) (string, error) {
	ipfsid, err := npFsClient.MFTemp2Files(fsid, destPath)
	if err != nil {
		return "", err
	}
	return ipfsid, nil
}

// ==================== 文件保存相关 ====================

// SaveLocalFileToNpfs 将本地文件保存到 NPFS
// localPath: 本地文件路径
// fsDir: NPFS 目录
// filename: 文件名
// ipfsid: IPFS 文件 ID，err: 错误信息
func SaveLocalFileToNpfs(localPath, fsDir, filename string) (string, error) {
	// 打开临时文件
	fsid, err := OpenTempFile()
	if err != nil {
		return "", err
	}

	// 读取本地文件
	fileData, err := ReadLocalFile(localPath)
	if err != nil {
		return "", err
	}

	// 写入数据到临时文件
	_, err = WriteDataToTempFile(fsid, fileData, 0)
	if err != nil {
		return "", err
	}

	// 确保目标目录存在
	err = EnsureDirExists(fsDir, true)
	if err != nil {
		return "", err
	}

	// 构建完整路径
	nodePath := fsDir + "/" + filename

	// 保存临时文件到 NPFS
	ipfsid, err := SaveTempFileToNpfs(fsid, nodePath)
	if err != nil {
		return "", err
	}

	return ipfsid, nil
}

// SaveContentToNpfs 将内容字符串保存到 NPFS
// content: 文件内容
// fsDir: NPFS 目录
// filename: 文件名
// ipfsid: IPFS 文件 ID，err: 错误信息
func SaveContentToNpfs(content, fsDir, filename string) (string, error) {
	// 打开临时文件
	fsid, err := OpenTempFile()
	if err != nil {
		return "", err
	}

	// 写入数据到临时文件
	_, err = WriteDataToTempFile(fsid, []byte(content), 0)
	if err != nil {
		return "", err
	}

	// 确保目标目录存在
	err = EnsureDirExists(fsDir, true)
	if err != nil {
		return "", err
	}

	// 构建完整路径
	nodePath := fsDir + "/" + filename

	// 保存临时文件到 NPFS
	ipfsid, err := SaveTempFileToNpfs(fsid, nodePath)
	if err != nil {
		return "", err
	}

	return ipfsid, nil
}

// ==================== 工具方法 ====================

// ReadLocalFile 读取本地文件
// path: 文件路径
// data: 文件内容，err: 错误信息
func ReadLocalFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ==================== 批量操作示例 ====================

// BatchCreateFilesExample 批量创建文件的示例（用于测试）
// count: 创建的文件数量
func BatchCreateFilesExample(count int) {
	wg := sync.WaitGroup{}
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func(index int) {
			defer wg.Done()

			// 打开临时文件
			fsid, err := OpenTempFile()
			if err != nil {
				return
			}

			// 写入数据
			_, err = WriteDataToTempFile(fsid, []byte(fmt.Sprintf("i:%d", index)), 0)
			if err != nil {
				return
			}

			// 检查并创建目录
			testDir := "test/"
			fmt.Println("testDir->", testDir)
			_, err = CreateDirIfNotExists(testDir, true)
			if err != nil {
				fmt.Println("err:", err.Error())
				return
			}

			// 构建文件路径
			nodePath := testDir + fmt.Sprintf("%d.txt", index)

			// 保存到 NPFS
			ipfsid, err := SaveTempFileToNpfs(fsid, nodePath)
			if err != nil {
				return
			}

			fmt.Printf("i:%d,filepath:%s\n", index, nodePath)
			fmt.Println("ipfsid", ipfsid)
		}(i)
	}

	wg.Wait()
}

// TestNpFsOperations 测试所有 NPFS 操作的示例函数
func TestNpFsOperations() {
	// 示例：读取文件
	data, size, err := ReadFileFromNpfs("/np_storage/1.jpg")
	if err != nil {
		panic(err)
	}
	fmt.Println("data:", string(data))
	fmt.Println("size:", size)

	// 示例：保存到本地
	err = SaveFileToLocal("/np_storage/1.jpg", "./hello.jpg")
	if err != nil {
		fmt.Println("保存失败:", err)
	}

	// 示例：保存本地文件到 NPFS
	ipfsid, err := SaveLocalFileToNpfs("./local.txt", "/my_files", "test.txt")
	if err != nil {
		fmt.Println("保存失败:", err)
	}
	fmt.Println("ipfsid:", ipfsid)

	// 示例：保存内容到 NPFS
	ipfsid2, err := SaveContentToNpfs("Hello World", "/my_files", "hello.txt")
	if err != nil {
		fmt.Println("保存失败:", err)
	}
	fmt.Println("ipfsid:", ipfsid2)

	// 示例：批量创建文件
	BatchCreateFilesExample(100)
}
