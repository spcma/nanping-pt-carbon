package infrastructure

import (
	"app/internal/module/carbonreport/domain"
	"app/internal/module/ipfs/rpc"
	"context"
	"os"
	"strings"
)

type npfsFileRepository struct {
	client  *rpc.LApiStub
	session string
}

// NewFileRepository 创建文件仓储
func NewFileRepository(client *rpc.LApiStub, session string) domain.FileRepository {
	return &npfsFileRepository{
		client:  client,
		session: session,
	}
}

func (r *npfsFileRepository) GetFile(ctx context.Context, path string) (*domain.File, error) {
	_, err := r.client.FilesStat(r.session, path)
	if err != nil {
		if isDirNotExist(err) {
			return nil, domain.ErrDirNotFound
		}
		return nil, err
	}

	// 如果是文件，返回文件信息
	return &domain.File{
		Path:  path,
		Name:  path[strings.LastIndex(path, "/")+1:],
		IsDir: false,
	}, nil
}

func (r *npfsFileRepository) ListDirectory(ctx context.Context, path string) ([]domain.File, error) {
	links, err := r.client.FilesLs(r.session, path)
	if err != nil {
		return nil, err
	}

	files := make([]domain.File, 0, len(links))
	for _, link := range links {
		file := domain.File{
			Path:  path + "/" + link.Name,
			Name:  link.Name,
			Size:  int64(link.Size),
			Hash:  link.Hash,
			IsDir: link.IsDir(),
		}
		files = append(files, file)
	}

	return files, nil
}

func (r *npfsFileRepository) CreateDirectory(ctx context.Context, path string, recursive bool) error {
	return r.client.FilesMkdir(r.session, path, recursive)
}

func (r *npfsFileRepository) DeleteFile(ctx context.Context, path string, recursive, force bool) error {
	return r.client.FilesRm(r.session, path, recursive, force)
}

func (r *npfsFileRepository) SaveContent(ctx context.Context, content []byte, dirPath, filename string) (domain.IPFSID, error) {
	// 打开临时文件
	fsid, err := r.client.MFOpenTempFile(r.session)
	if err != nil {
		return "", err
	}

	// 写入数据
	_, err = r.client.MFSetData(fsid, content, 0)
	if err != nil {
		return "", err
	}

	// 确保目录存在
	err = r.ensureDirExists(dirPath, true)
	if err != nil {
		return "", err
	}

	// 保存到 NPFS
	nodePath := dirPath + "/" + filename
	ipfsid, err := r.client.MFTemp2Files(fsid, nodePath)
	if err != nil {
		return "", err
	}

	return domain.IPFSID(ipfsid), nil
}

func (r *npfsFileRepository) SaveLocalFile(ctx context.Context, localPath, fsDir, filename string) (domain.IPFSID, error) {
	// 读取本地文件
	data, err := r.readFile(localPath)
	if err != nil {
		return "", err
	}

	// 打开临时文件
	fsid, err := r.client.MFOpenTempFile(r.session)
	if err != nil {
		return "", err
	}

	// 写入数据
	_, err = r.client.MFSetData(fsid, data, 0)
	if err != nil {
		return "", err
	}

	// 确保目录存在
	err = r.ensureDirExists(fsDir, true)
	if err != nil {
		return "", err
	}

	// 保存到 NPFS
	nodePath := fsDir + "/" + filename
	ipfsid, err := r.client.MFTemp2Files(fsid, nodePath)
	if err != nil {
		return "", err
	}

	return domain.IPFSID(ipfsid), nil
}

func (r *npfsFileRepository) ReadFile(ctx context.Context, filePath string) ([]byte, int64, error) {
	// 打开文件 URL
	fsid, err := r.client.MMOpenUrl(r.session, filePath)
	if err != nil {
		return nil, 0, err
	}
	defer r.client.MMClose(fsid)

	// 获取文件大小
	size, err := r.client.MFGetSize(fsid)
	if err != nil {
		return nil, 0, err
	}

	// 读取文件数据
	data, err := r.client.MFGetData(fsid, 0, int(size))
	if err != nil {
		return nil, 0, err
	}

	return data, size, nil
}

// ensureDirExists 确保目录存在
func (r *npfsFileRepository) ensureDirExists(path string, recursive bool) error {
	exists, err := r.checkDirExists(path)
	if err != nil {
		return err
	}

	if !exists {
		return r.client.FilesMkdir(r.session, path, recursive)
	}

	return nil
}

// checkDirExists 检查目录是否存在
func (r *npfsFileRepository) checkDirExists(path string) (bool, error) {
	_, err := r.client.FilesStat(r.session, path)
	if err != nil {
		if isDirNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// isDirNotExist 判断错误是否表示目录不存在
func isDirNotExist(err error) bool {
	return domain.IsDirNotExist(err)
}

// readFile 读取本地文件
func (r *npfsFileRepository) readFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}
