package meta

import (
	mydb "Cloud-Storage-System-based-on-Golang/filestore-server/db"
	"sort"
)

//Filemeta:文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)

}

//UpdateFileMeta:新增/更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

//UpdateFileMetaDB:新增更新文件元信息到mysql中
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(
		fmeta.FileSha1, fmeta.FileName, fmeta.FileSize,
		fmeta.Location)
}

//GetFileMeta:通过sha1值获取文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

//GetFileMetaDB从mysql获取文件元信息
func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	tfile, err := mydb.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileSha1 : tfile.FileHash,
		FileName : tfile.FileName.String,
		FileSize : tfile.FileSize.Int64,
		Location : tfile.FileAddr.String,
	}
	return fmeta, nil
}

//获取批量的文件元信息列表
func GetLastFileMetas(count int) []FileMeta {
	// fMetaArray := make([]FileMeta, len(fileMetas))
	//修复数组元素添加新的bug
	//也可改为：fMetaArray := make([]FileMeta, 0)
	var fMetaArray []FileMeta
	for _, v := range fileMetas {
		fMetaArray = append(fMetaArray, v)
	}

	sort.Sort(ByUploadTime(fMetaArray))
	return fMetaArray[0:count]
}

//RemoveFileMeta:删除元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
