package tapbag

import (
	"os"
	"path/filepath"
	"strings"
)

type FileManager struct {
	SrcPath string
	Types map[string]string
	OpenableTypes []string
	DownloadQueue map[string]*FileSystem
}

func NewFileManager(srcPath string) *FileManager {
	fm := &FileManager{
		SrcPath: srcPath,
		Types: map[string]string{
			// archive
			".zip": "archive",
			".gz":  "archive",
			".tar": "archive",
			".xz":  "archive",
			// audio
			".mp3": "audio",
			// doc
			".doc": "doc",
			// image
			".png": "image",
			".jpg": "image",
			".svg": "image",
			// pdf
			".pdf": "pdf",
			// ppt
			".ppt": "ppt",
			// psd
			".psd": "psd",
			// text
			".txt": "text",
			// coder
			".py":   "text",
			".cpp":  "text",
			".h":    "text",
			".hpp":  "text",
			".c":    "text",
			".go":   "text",
			".java": "text",
			".js":   "text",
			".html": "text",
			".css":  "text",
			".vue":  "text",
			".md": "text",
			// video
			".mp4": "video",
			".avi": "video",
			// xls
			".xls": "xls",
		},
		OpenableTypes: []string{"folder", "audio", "image", "text", "video"},
		DownloadQueue: make(map[string]*FileSystem),
	}
	return fm
}

func (this *FileManager) getType(info os.FileInfo) string {
	if info.IsDir() {
		return "folder"
	}
	extName := strings.ToLower(filepath.Ext(info.Name()))
	if t, ok := this.Types[extName]; ok {
		return t
	}
	return "file"
}

func (this *FileManager) fileOpenable(info os.FileInfo) bool {
	ftype := this.getType(info)
	for _, it := range this.OpenableTypes {
		if ftype == it {
			return true
		}
	}
	return false
}
