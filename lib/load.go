package lib

import (
	"bufio"
	"os"
	"path"
	"runtime"
)

func ScanFile(fileName string) (*bufio.Reader, *bufio.Scanner, error) {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))

	file, err := os.Open(path.Join(d, "../", fileName))
	if err != nil {
		return nil, nil, err
	}

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	return reader, scanner, nil
}
