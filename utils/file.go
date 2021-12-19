package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func GetAllRules(pocDir string,) ([]string,error){
	var files []string
	if _,err:=os.Stat(pocDir);os.IsNotExist(err){
		return nil,err
	}
	asbPath, _ := filepath.Abs(pocDir)
	err := filepath.Walk(asbPath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path,".yml")||strings.HasSuffix(path,".yaml"){
			files = append(files, path)
			Info("Loading --- "+path)
		}
		return nil
	})
	return files,err
}

func ReadingLines(filename string) []string {
	var result []string
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return result
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val := scanner.Text()
		if val == "" {
			continue
		}
		if strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {}else {val="http://"+val }

		result = append(result, val)
	}

	if err := scanner.Err(); err != nil {
		return result
	}
	return result
}
