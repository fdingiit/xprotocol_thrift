package zoneclient

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
)

func gunzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	// nolint
	defer r.Close()

	return ioutil.ReadAll(r)
}

func readFile(f string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Clean(f))
}

func writeFile(fileName string, fileData []byte) error {
	if isExist, _ := existFile(fileName); !isExist {
		dir, err := filepath.Abs(filepath.Dir(fileName))
		if err != nil {
			return err
		}
		if err = os.MkdirAll(dir, 0750); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(fileName, fileData, 0644)
}

func existFile(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func merrorfmt(errs []error) string {
	var b strings.Builder
	for i, err := range errs {
		_, _ = b.WriteString("#")
		_, _ = b.WriteString(strconv.Itoa(i))
		_, _ = b.WriteString(err.Error())
	}
	return b.String()
}

func merror(err error, errs ...error) error {
	noerr := true
	if err == nil {
		for _, err := range errs {
			if err != nil {
				noerr = false
			}
		}
	} else {
		noerr = false
	}

	if noerr {
		return nil
	}

	merr := multierror.Append(err, errs...)
	merr.ErrorFormat = merrorfmt
	return merr
}

func isPressFlow(flowType string) bool {
	if strings.TrimSpace(flowType) == "PRESS" {
		return true
	}
	return false
}
