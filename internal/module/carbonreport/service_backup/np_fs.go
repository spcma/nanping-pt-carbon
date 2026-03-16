package service

import (
	rpc2 "app/internal/module/ipfs/rpc"
	"fmt"
	"os"
	"strings"
	"sync"
)

var _NpFs = new(NpFs)

type NpFs struct{}

func NpFsServ() *NpFs {
	return _NpFs
}

var fsClient *rpc2.LApiStub
var sid string

func InitNpFs() {
	client, curSid, err := stubClient()
	if err != nil {
		panic(err)
	}
	defer client.Logout(sid, "")

	fsClient = client
	sid = curSid
}

func CloseNpFs() {
	fsClient.Logout(sid, "")
}

// stubClient 获取client
func stubClient() (*rpc2.LApiStub, string, error) {
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

// fsDirExist 判断目录是否存在
func fsDirExist(err error) bool {
	return !strings.Contains(err.Error(), "no link named")
}

type NpfsParam struct {
	Path     string `json:"path" form:"path"`
	Filename string `json:"filename" form:"filename"`
	Content  string `json:"content" form:"content"`
}

func T() {

	path := ""

	//	检查目录是否存在
	_, err := fsClient.FilesStat(sid, path)
	if err != nil {
		if !fsDirExist(err) {
			//目录不存在，创建新的目录
			err = fsClient.FilesMkdir(sid, path, true)
			if err != nil {
				return
			}

			return
		}
	}

	//	删除文件
	err = fsClient.FilesRm(sid, path, true, true)
	if err != nil {
		return
	}

	//列出目录，这里列上一级
	links, err := fsClient.FilesLs(sid, path)
	if err != nil {
		return
	}

	_ = links

	//这里测试读取
	// files/ps 		/ps	=	files/ps
	fsid2, err := fsClient.MMOpenUrl(sid, "/np_storage"+"/1.jpg")
	if err != nil {
		panic(err)
	}
	defer fsClient.MMClose(fsid2)

	size, err := fsClient.MFGetSize(fsid2)
	if err != nil {
		return
	}
	data, err := fsClient.MFGetData(fsid2, 0, int(size))
	if err != nil {
		panic(err)
	}
	fmt.Println("data:", string(data))

	err = os.WriteFile("./hello.jpg", data, os.ModePerm)
	if err != nil {
	}

	mmOpenUrl, err := fsClient.MMOpenUrl(sid, path)
	if err != nil {
		return
	}

	bytesss, err := fsClient.MFGetData(mmOpenUrl, 0, -1)
	if err != nil {
		return
	}

	_ = bytesss

	//// 设置响应头
	//c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+url.PathEscape(uuid.NewEasy()+".png"))
	//c.Header("Content-Type", "application/octet-stream")
	//c.Header("Content-Length", fmt.Sprintf("%d", len(bytesss)))
	//
	//// 使用 io.Copy 将文件内容复制到响应写入器中
	//_, err = io.Copy(c.Writer, bytes.NewBuffer(bytesss))
	//if err != nil {
	//	ErrorLog(c, err)
	//	return
	//}

	//	检查目录是否存在
	_, err = fsClient.FilesStat(sid, path)
	if err != nil {
		if !fsDirExist(err) {
			//目录不存在，创建新的目录
			err = fsClient.FilesMkdir(sid, path, true)
			if err != nil {
				return
			}
		} else {
			return
		}
	}

	//	open npfs tmpfile
	fsid, err := fsClient.MFOpenTempFile(sid)
	if err != nil {
		return
	}

	content := ""
	//	set data
	_, err = fsClient.MFSetData(fsid, []byte(content), 0)
	if err != nil {
		return
	}

	filename := ""
	// save tmpfile to npfs
	ipfsid, err := fsClient.MFTemp2Files(fsid, path+"/"+filename+".txt")
	if err != nil {
		return
	}

	_ = ipfsid

	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			//	open npfs tmpfile
			fsid, err := fsClient.MFOpenTempFile(sid)
			if err != nil {
				return
			}

			//	set data
			_, err = fsClient.MFSetData(fsid, []byte(fmt.Sprintf("i:%d", i)), 0)
			if err != nil {
				return
			}

			//	check npfs dir

			ddd := "test/"
			fmt.Println("ddd->", ddd)
			_, err = fsClient.FilesStat(sid, "test/")
			if err != nil {
				if strings.Contains(err.Error(), "no link named") {
					//创建目录
					err = fsClient.FilesMkdir(sid, "test/", true)
					if err != nil {
						return
					}
				} else {
					fmt.Println("err:", err.Error())
					return
				}
			}

			//	storageDir
			nodePath := "test/" + fmt.Sprintf("%d.txt", i)

			// save tmpfile to npfs
			ipfsid, err := fsClient.MFTemp2Files(fsid, nodePath)
			if err != nil {
				return
			}

			fmt.Println(fmt.Sprintf("i:%d,filepath:%s", i, nodePath))
			fmt.Println("ipfsid", ipfsid)
		}()
	}
	wg.Wait()
}

// fsSave 保存文件到npfs
//
//	localPath 本地文件路径
//	fsDir npfs目录
//	filename 文件名
func fsSave(localPath, fsDir, filename string) (string, error) {

	//	open npfs tmpfile
	fsid, err := fsClient.MFOpenTempFile(sid)
	if err != nil {
		return "", err
	}

	//	read tmp local file
	readFile, err := ReadLocalFile(localPath)
	if err != nil {
		return "", err
	}

	//	set data
	_, err = fsClient.MFSetData(fsid, readFile, 0)
	if err != nil {
		return "", err
	}

	//	check npfs dir
	_, err = fsClient.FilesStat(sid, fsDir)
	if err != nil {
		if !fsDirExist(err) {
			//创建目录
			err = fsClient.FilesMkdir(sid, fsDir, true)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	nodePath := fsDir + "/" + filename

	// save tmpfile to npfs
	ipfsid, err := fsClient.MFTemp2Files(fsid, nodePath)
	if err != nil {
		return "", err
	}

	return ipfsid, nil
}

//  {host}/ipfs/QmWfVY9y3xjsixTgbd9AorQxH7VtMpzfx2HaWtsoUYecaX
