package utils

import (
	"os"
	"io/ioutil"
	"path"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type MyFile struct {
	Path  string
	Name  string
	Dir   string
	Ext   string
	IsDir bool
}

func FileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func FileNotExist(path string) bool {
	return !FileExist(path)
}

func FileInfo(filepath string) (os.FileInfo, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	file.Close()
	return fi, nil
}

func FileReadAllBytes(path string) []byte {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		panic(err)
	}
	return fd
}

func FileReadAllString(path string) string {
	return string(FileReadAllBytes(path))
}

func FileReadLineString(path string) []string {
	s := FileReadAllString(path)
	return strings.Split(s,"\n")
}

func FileWriteBytes(path string, bytes []byte) {
	ioutil.WriteFile(path, bytes, 0777)
}

func FileWriteString(path string, str string) {
	err := ioutil.WriteFile(path, []byte(str), 0777)
	if err != nil {
		fmt.Printf("Write file failed!%s", err)
	}
}

func FileSizeFriendly(fileSize int64) string {
	if fileSize > 1024 * 1024 * 1024 {
		val := float64(fileSize) / (1024 * 1024 * 1024.0)
		return fmt.Sprintf("%.2f GB", val)
	} else if (fileSize > 1024 * 1024) {
		val := float64(fileSize) / (1024 * 1024.0)
		return fmt.Sprintf("%.2f MB", val)
	} else if (fileSize > 1024) {
		val := float64(fileSize) / (1024.0)
		return fmt.Sprintf("%.2f KB", val)
	} else {
		return fmt.Sprintf("%d Byte", fileSize)
	}

	return ""
}

/* 仅包含文件，不包含目录 */
func FilesAtDir(dir string) []MyFile {
	result := make([]MyFile, 0)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			result = append(result, MyFile{
				Path:path,
				Name:info.Name(),
				Ext:filepath.Ext(path),
				Dir:strings.Replace(path, info.Name(), "", -1),
				IsDir:false,
			})
		}
		return nil
	})

	return result
}

/* 包含所有文件：目录和文件 */
func FilesAllAtDir(dir string) []MyFile {
	result := make([]MyFile, 0)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			result = append(result, MyFile{
				Path:path,
				Name:info.Name(),
				Ext:filepath.Ext(path),
				Dir:strings.Replace(path, info.Name(), "", -1),
				IsDir:false,
			})
		} else {
			result = append(result, MyFile{
				Path:path,
				Name:info.Name(),
				Ext:"",
				Dir:strings.Replace(path, info.Name(), "", -1),
				IsDir:true,
			})
		}
		return nil
	})

	return result
}

func Copy_folder(source string, dest string) (err error) {

	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			err = Copy_folder(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err = Copy_file(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

func Copy_file(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}

func ExeFullPath() string {
	dir, _ := filepath.Abs(os.Args[0])
	return dir
}

func ExeDir() string {
	return path.Dir(ExeFullPath())
}

func ExeDirAppend(relativePath string) string {
	return path.Join(ExeDir(), relativePath)
}

func ExePath() string {
	dir, _ := os.Getwd()
	return dir
}

func MakeDir(dirPath string) string {
	dirnameTMP := fmt.Sprintf(dirPath)
	err := os.MkdirAll(dirnameTMP, 0755)
	if err != nil {
		return ""
	}else {
		return dirnameTMP
	}
}

