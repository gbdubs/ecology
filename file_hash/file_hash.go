package file_hash

import (
	"crypto/sha1"
	"fmt"
	"github.com/gbdubs/ecology/output"
	"io/ioutil"
)

func ComputeFileHash(filePath string, o *output.Output) (hashResult string, err error) {
	o.Info("Computing Hash Of %s", filePath).Indent()
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		o.Error(err)
		return
	}
	hashFn := sha1.New()
	hashFn.Write([]byte(data))
	hashResult = fmt.Sprintf("%x", hashFn.Sum(nil))
	return
}
