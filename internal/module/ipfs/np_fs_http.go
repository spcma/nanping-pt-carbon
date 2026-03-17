package ipfs

import (
	"app/internal/module/ipfs/rpc"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// NpFsHTTPHandler NPFS HTTP 处理器
type NpFsHTTPHandler struct {
	client *rpc.LApiStub
	sid    string
}

// NewNpFsHTTPHandler 创建 HTTP 处理器
func NewNpFsHTTPHandler(client *rpc.LApiStub, curSid string) (*NpFsHTTPHandler, error) {
	return &NpFsHTTPHandler{
		client: client,
		sid:    curSid,
	}, nil
}

// Close 关闭连接
func (h *NpFsHTTPHandler) Close() {
	if h.client != nil {
		h.client.Logout(h.sid, "")
	}
}

// HTTPResponse 统一响应结构
type HTTPResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, HTTPResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, HTTPResponse{
		Code:    code,
		Message: message,
	})
}

// CheckDirExist 检查目录是否存在
func (h *NpFsHTTPHandler) CheckDirExist(path string) bool {
	_, err := h.client.FilesStat(h.sid, path)
	if err != nil {
		return !strings.Contains(err.Error(), "no link named")
	}
	return true
}

// CreateDir 创建目录
func (h *NpFsHTTPHandler) CreateDir(path string, recursive bool) error {
	return h.client.FilesMkdir(h.sid, path, recursive)
}

// ListDir 列出目录
func (h *NpFsHTTPHandler) ListDir(path string) ([]rpc.LsLink, error) {
	return h.client.FilesLs(h.sid, path)
}

// DeleteFile 删除文件
func (h *NpFsHTTPHandler) DeleteFile(path string, recursive, force bool) error {
	return h.client.FilesRm(h.sid, path, recursive, force)
}

// ReadFile 读取文件
func (h *NpFsHTTPHandler) ReadFile(filePath string) ([]byte, int64, error) {
	fsid, err := h.client.MMOpenUrl(h.sid, filePath)
	if err != nil {
		return nil, 0, err
	}
	defer h.client.MMClose(fsid)

	size, err := h.client.MFGetSize(fsid)
	if err != nil {
		return nil, 0, err
	}

	data, err := h.client.MFGetData(fsid, 0, int(size))
	if err != nil {
		return nil, 0, err
	}

	return data, size, nil
}

// SaveContent 保存内容到文件
func (h *NpFsHTTPHandler) SaveContent(content, fsDir, filename string) (string, error) {
	// 打开临时文件
	fsid, err := h.client.MFOpenTempFile(h.sid)
	if err != nil {
		return "", err
	}

	// 写入数据
	_, err = h.client.MFSetData(fsid, []byte(content), 0)
	if err != nil {
		return "", err
	}

	// 确保目录存在
	err = h.EnsureDirExists(fsDir, true)
	if err != nil {
		return "", err
	}

	// 保存到 NPFS
	nodePath := fsDir + "/" + filename
	ipfsid, err := h.client.MFTemp2Files(fsid, nodePath)
	if err != nil {
		return "", err
	}

	return ipfsid, nil
}

// SaveLocalFile 保存本地文件到 NPFS
func (h *NpFsHTTPHandler) SaveLocalFile(localPath, fsDir, filename string) (string, error) {
	// 打开临时文件
	fsid, err := h.client.MFOpenTempFile(h.sid)
	if err != nil {
		return "", err
	}

	// 读取本地文件
	data, err := os.ReadFile(localPath)
	if err != nil {
		return "", err
	}

	// 写入数据
	_, err = h.client.MFSetData(fsid, data, 0)
	if err != nil {
		return "", err
	}

	// 确保目录存在
	err = h.EnsureDirExists(fsDir, true)
	if err != nil {
		return "", err
	}

	// 保存到 NPFS
	nodePath := fsDir + "/" + filename
	ipfsid, err := h.client.MFTemp2Files(fsid, nodePath)
	if err != nil {
		return "", err
	}

	return ipfsid, nil
}

// EnsureDirExists 确保目录存在
func (h *NpFsHTTPHandler) EnsureDirExists(path string, recursive bool) error {
	exists := h.CheckDirExist(path)
	if !exists {
		return h.CreateDir(path, recursive)
	}
	return nil
}

// SetupRouter 设置路由
func (h *NpFsHTTPHandler) SetupRouter() *gin.Engine {
	router := gin.Default()

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		Success(c, gin.H{"status": "ok"})
	})

	// API 路由组
	api := router.Group("/api/v1")
	{
		// 目录操作
		api.POST("/dir/check", h.checkDirHandler)
		api.POST("/dir/create", h.createDirHandler)
		api.GET("/dir/list", h.listDirHandler)
		api.DELETE("/dir/delete", h.deleteDirHandler)

		// 文件操作
		api.GET("/file/read", h.readFileHandler)
		api.POST("/file/save", h.saveFileHandler)
		api.POST("/file/upload", h.uploadFileHandler)
		api.GET("/file/download", h.downloadFileHandler)
		api.DELETE("/file/delete", h.deleteFileHandler)
		api.GET("/calc", h.CalcOnceHandler)
	}

	return router
}

// checkDirHandler 检查目录 handler
func (h *NpFsHTTPHandler) checkDirHandler(c *gin.Context) {
	var param struct {
		Path string `json:"path" form:"path" binding:"required"`
	}

	if err := c.ShouldBind(&param); err != nil {
		Error(c, 400, "参数错误："+err.Error())
		return
	}

	exists := h.CheckDirExist(param.Path)
	Success(c, gin.H{
		"path":   param.Path,
		"exists": exists,
	})
}

// createDirHandler 创建目录 handler
func (h *NpFsHTTPHandler) createDirHandler(c *gin.Context) {
	var param struct {
		Path      string `json:"path" form:"path" binding:"required"`
		Recursive bool   `json:"recursive" form:"recursive"`
	}

	if err := c.ShouldBind(&param); err != nil {
		Error(c, 400, "参数错误："+err.Error())
		return
	}

	if param.Recursive {
		param.Recursive = true
	}

	err := h.CreateDir(param.Path, param.Recursive)
	if err != nil {
		Error(c, 500, "创建失败："+err.Error())
		return
	}

	Success(c, gin.H{
		"path":    param.Path,
		"created": true,
	})
}

// listDirHandler 列出目录 handler
func (h *NpFsHTTPHandler) listDirHandler(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		Error(c, 400, "缺少参数：path")
		return
	}

	links, err := h.ListDir(path)
	if err != nil {
		Error(c, 500, "列出失败："+err.Error())
		return
	}

	// 转换为 JSON 友好的格式
	result := make([]gin.H, 0, len(links))
	for _, link := range links {
		result = append(result, gin.H{
			"name": link.Name,
			"type": func() string {
				if link.IsDir() {
					return "directory"
				}
				return "file"
			}(),
			"size": link.Size,
		})
	}

	Success(c, gin.H{
		"path":  path,
		"files": result,
	})
}

// deleteDirHandler 删除目录 handler
func (h *NpFsHTTPHandler) deleteDirHandler(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		Error(c, 400, "缺少参数：path")
		return
	}

	recursive := c.Query("recursive") == "true"
	force := c.Query("force") == "true"

	err := h.DeleteFile(path, recursive, force)
	if err != nil {
		Error(c, 500, "删除失败："+err.Error())
		return
	}

	Success(c, gin.H{
		"path":    path,
		"deleted": true,
	})
}

// readFileHandler 读取文件 handler
func (h *NpFsHTTPHandler) readFileHandler(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		Error(c, 400, "缺少参数：path")
		return
	}

	data, size, err := h.ReadFile(path)
	if err != nil {
		Error(c, 500, "读取失败："+err.Error())
		return
	}

	Success(c, gin.H{
		"path": path,
		"size": size,
		"content": func() string {
			// 如果是文本文件，返回内容
			if size < 10000 {
				return string(data)
			}
			return "(文件太大，仅显示前 10000 字节)"
		}(),
	})
}

// saveFileHandler 保存文件 handler（从请求体保存）
func (h *NpFsHTTPHandler) saveFileHandler(c *gin.Context) {
	var param struct {
		Content  string `json:"content" binding:"required"`
		Dir      string `json:"dir" binding:"required"`
		Filename string `json:"filename" binding:"required"`
	}

	if err := c.ShouldBindJSON(&param); err != nil {
		Error(c, 400, "参数错误："+err.Error())
		return
	}

	ipfsid, err := h.SaveContent(param.Content, param.Dir, param.Filename)
	if err != nil {
		Error(c, 500, "保存失败："+err.Error())
		return
	}

	Success(c, gin.H{
		"ipfsid": ipfsid,
		"path":   param.Dir + "/" + param.Filename,
	})
}

// uploadFileHandler 上传文件 handler（multipart/form-data）
func (h *NpFsHTTPHandler) uploadFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		Error(c, 400, "缺少文件参数")
		return
	}

	dir := c.PostForm("dir")
	if dir == "" {
		Error(c, 400, "缺少参数：dir")
		return
	}

	filename := c.PostForm("filename")
	if filename == "" {
		filename = file.Filename
	}

	// 保存上传的文件到临时位置
	tmpPath := filepath.Join(os.TempDir(), file.Filename)
	if err := c.SaveUploadedFile(file, tmpPath); err != nil {
		Error(c, 500, "保存临时文件失败："+err.Error())
		return
	}
	defer os.Remove(tmpPath)

	// 保存到 NPFS
	ipfsid, err := h.SaveLocalFile(tmpPath, dir, filename)
	if err != nil {
		Error(c, 500, "上传失败："+err.Error())
		return
	}

	Success(c, gin.H{
		"ipfsid": ipfsid,
		"path":   dir + "/" + filename,
	})
}

// downloadFileHandler 下载文件 handler
func (h *NpFsHTTPHandler) downloadFileHandler(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		Error(c, 400, "缺少参数：path")
		return
	}

	filename := c.Query("filename")
	if filename == "" {
		filename = filepath.Base(path)
	}

	data, _, err := h.ReadFile(path)
	if err != nil {
		Error(c, 500, "读取失败："+err.Error())
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", len(data)))

	c.Data(http.StatusOK, "application/octet-stream", data)
}

// deleteFileHandler 删除文件 handler
func (h *NpFsHTTPHandler) deleteFileHandler(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		Error(c, 400, "缺少参数：path")
		return
	}

	force := c.Query("force") == "true"

	err := h.DeleteFile(path, false, force)
	if err != nil {
		Error(c, 500, "删除失败："+err.Error())
		return
	}

	Success(c, gin.H{
		"path":    path,
		"deleted": true,
	})
}

// CalcOnceHandler 单次计算处理
func (h *NpFsHTTPHandler) CalcOnceHandler(c *gin.Context) {
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

				records, err := parseFile(link.Name)
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
