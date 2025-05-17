package detections

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spoofimei/syscanary/internal/global"
)

var (
	hashMap map[string]string
)

func FileIntegrityDetect(output chan string, stop chan *struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	output <- "debug:integrity:integrity monitoring started"

	for {
		for _, path := range global.Config.Modules.Integrity.Paths {
			out, err := process(path)
			if err != nil {
				output <- fmt.Sprintf("error:integrity:%s", err.Error())
			}
			if out != "" {
				output <- fmt.Sprintf("info:integrity:%s", out)
			}
		}

		if len(stop) > 0 {
			break
		}
		time.Sleep(time.Duration(global.Config.Modules.Integrity.Interval) * time.Second)
	}
}

func process(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	var files []string
	if info.IsDir() {
		filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				files = append(files, path)
			}
			return nil
		})
	}

	if len(hashMap) == 0 {
		hashMap, err = calculateHashes(files)
		if err != nil {
			return "", err
		}
		return "", nil
	}

	var out string
	newHashMap, err := calculateHashes(files)
	if err != nil {
		return "", err
	}
	for refFile, refHash := range hashMap {
		hash, ok := newHashMap[refFile]

		if !ok {
			out = fmt.Sprintf("file '%s' has been deleted", refFile)
			break
		}
		if hash != refHash {
			out = fmt.Sprintf("file '%s' has been modified", refFile)
		}
	}
	for newFile, _ := range newHashMap {
		if _, ok := hashMap[newFile]; !ok {
			out = fmt.Sprintf("file '%s' has been created", newFile)
		}
	}
	hashMap = newHashMap
	return out, nil
}

func calculateHashes(files []string) (map[string]string, error) {
	tmpHashMap := make(map[string]string)
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		hashBytes := sha256.Sum256(content)
		hash := hex.EncodeToString(hashBytes[:])
		tmpHashMap[file] = hash
	}
	return tmpHashMap, nil
}
