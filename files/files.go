package files

import (
	"fmt"
	"gocliscaff/common"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/gookit/color"
	"github.com/spf13/viper"
)

type File struct {
	Path           string
	ByteSize       int64
	PrettyByteSize string
}

func ReadDirRecursively(dirPath string) ([]File, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var results []File
	for _, fileInfo := range files {
		filePath := filepath.Join(dirPath, fileInfo.Name())
		if fileInfo.IsDir() {
			subFiles, err := ReadDirRecursively(filePath)
			if err != nil {
				return nil, err
			}
			results = append(results, subFiles...)
		} else {
			if fileInfo.Size() >= (int64(viper.GetInt("minfilesize")) * 100000) {
				foundFile := File{}
				foundFile.Path = filePath
				foundFile.ByteSize = fileInfo.Size()
				foundFile.PrettyByteSize = common.PrettyBytes(fileInfo.Size())

				results = append(results, foundFile)
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].ByteSize > results[j].ByteSize
	})

	return results, nil
}

func PrintResults(files []File) {
	fmt.Println()
	common.PrintColor("darkturquoise", "background", fmt.Sprintf("Largest files in: %s", viper.GetString("path")))
	fmt.Println()

	spacing := make(map[string]int)
	highWaterMark := 0

	for _, file := range files {
		if len(file.Path) > highWaterMark {
			highWaterMark = len(file.Path)
		}

		spacing[file.Path] = len(file.Path)
	}

	for _, file := range files {
		padding := strconv.Itoa(highWaterMark + 2)

		if file.ByteSize >= int64(viper.GetInt("highlight")*1000000) {
			color.HEXStyle("000", common.AllHex["yellow1"]).Printf("%-"+padding+"s %10s\n", file.Path, file.PrettyByteSize)
		} else {
			color.HEXStyle(common.AllHex["steelblue2"]).Printf("%-"+padding+"s %10s\n", file.Path, file.PrettyByteSize)
		}
	}
}
