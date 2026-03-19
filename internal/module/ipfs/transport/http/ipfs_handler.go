package http

import (
	"app/internal/module/ipfs/application"
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// IpfsHandler NPFS HTTP 处理器
type IpfsHandler struct {
	appService *application.Service
}

// NewIpfsHandler 创建 HTTP 处理器
func NewIpfsHandler(appService *application.Service) (*IpfsHandler, error) {
	return &IpfsHandler{
		appService: appService,
	}, nil
}

// ListDir 列出目录
func (h *IpfsHandler) ListDir(c *gin.Context) {
	dir := c.Query("dir")
	h.appService.ListDir(platform_http.Ctx(c), dir)
}

// CreateDir 创建目录
func (h *IpfsHandler) CreateDir(c *gin.Context) {
	h.appService.CreateDir(platform_http.Ctx(c), "")
}

// DeleteFile 删除文件
func (h *IpfsHandler) DeleteFile(c *gin.Context) {
	var req application.DeleteFileRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return
	}

	//path string, recursive, force bool
	h.appService.DeleteFile(platform_http.Ctx(c), &req)
}

// UploadFileHandler 上传文件 handler（multipart/form-data）
func (h *IpfsHandler) UploadFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		return
	}

	dir := c.PostForm("dir")
	if dir == "" {
		return
	}

	filename := c.PostForm("filename")
	if filename == "" {
		filename = file.Filename
	}

	// 保存上传的文件到临时位置
	tmpPath := filepath.Join(os.TempDir(), file.Filename)
	if err := c.SaveUploadedFile(file, tmpPath); err != nil {
		return
	}
	defer os.Remove(tmpPath)

	// 保存到 NPFS
	ipfsid, err := h.SaveLocalFile(tmpPath, dir, filename)
	if err != nil {
		return
	}

	response.Success(c, gin.H{
		"ipfsid": ipfsid,
		"path":   dir + "/" + filename,
	})
}

// DownloadFileHandler 下载文件 handler
func (h *IpfsHandler) DownloadFileHandler(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		return
	}

	filename := c.Query("filename")
	if filename == "" {
		filename = filepath.Base(path)
	}

	data, _, err := h.ReadFile(path)
	if err != nil {
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", len(data)))

	c.Data(http.StatusOK, "application/octet-stream", data)
}

// CalcOnceHandler 单次计算处理
func (h *IpfsHandler) CalcOnceHandler(c *gin.Context) {
	dir := "/aibk/26/03/01/19BONNcTiriywcgCUSOv82Pj-i6/"

	links, err := h.ListDir(dir)
	if err != nil {
		Error(c, 500, "列出失败："+err.Error())
		return
	}

	for _, link := range links {
		if link.Type == 1 {
			//	是目录
		}
		if link.Type == 0 {
			//	是文件

			if strings.Contains(link.Name, "txt") {

				//err = SaveFileToLocal(dir+link.Name, link.Name)
				//if err != nil {
				//	return
				//}

				//data, size, err := h.ReadFile(dir + "gps_20260301003555.txt")
				//if err != nil {
				//	Error(c, 500, "读取失败："+err.Error())
				//	return
				//}

				//f, err := os.Create(link.Name)
				//if err != nil {
				//	return
				//}

				// 写入字节数组
				//n, err := f.Write(data)
				//if err != nil {
				//	fmt.Println("写入文件失败:", err)
				//	return
				//}
				//
				//fmt.Println("write byte: ", n)

				records, err := application.parseFile(link.Name)
				if err != nil {
					//f.Close()
					fmt.Errorf("%v", err.Error())
					return
				}

				//fmt.Println("filename", link.Name, "hash", link.Hash, "size:", size, "data: \n", string(data))
				fmt.Println(records)

				for i, record := range records {
					fmt.Println("index: ", i, "record:", record)
				}

				//f.Close()

				f, err := os.OpenFile(link.Name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
				if err != nil {
					fmt.Errorf("%v", err)
					return
				}

				writeString, err := f.WriteString("aaaaa")
				if err != nil {
					fmt.Errorf("%v", err)
					return
				}

				fmt.Println("write n: ", writeString)

				f.Close()

				//	如果文件已经存在在npfs目录中，那么同名文件不会覆盖，会写入失败，需要先删除后写入，可选择将文件读取到另外的目录中，然后进行删除写入操作
				//ipfsId, err := SaveLocalFileToNpfs(link.Name, "/tmpp/", link.Name)
				//if err != nil {
				//	return
				//}

				//fmt.Println("ipfsId:", ipfsId)

				return
			} else {
				fmt.Println("filename", link.Name, "hash", link.Hash)
			}
		}
	}

}

// IsDirNotExist 判断错误是否表示目录不存在
func IsDirNotExist(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "no link named") ||
		strings.Contains(err.Error(), "no linked named")
}
