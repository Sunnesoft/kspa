package graph_shortest_paths

import (
	"io/ioutil"
	"os"
)

func LoadText(fn string) ([]byte, error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return ioutil.ReadAll(file)
}

func WriteText(fn string, data []byte) error {
	return os.WriteFile(fn, data, 0666)
}
