package file_hash

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
)

func ComputeFileHash(filePath string) (hashResult string, err error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	hashFn := sha256.New()
	hashFn.Write([]byte(data))
	hashResult = fmt.Sprintf("%x", hashFn.Sum(nil))
	return
}
