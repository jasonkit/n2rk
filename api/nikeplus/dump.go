package nikeplus

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

func Export(activities []*Activity, filePath string) {
	var buf, outBuf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(activities); err != nil {
		fmt.Printf("Failed to encode: %v\n", err)
		return
	}

	w := gzip.NewWriter(&outBuf)
	w.Write(buf.Bytes())
	w.Close()

	if err := ioutil.WriteFile(filePath, outBuf.Bytes(), 0666); err != nil {
		fmt.Printf("Failed to write: %v", err)
		return
	}
}

func Import(filePath string) []*Activity {
	var activities []*Activity
	var fp *os.File
	var err error

	if fp, err = os.Open(filePath); err != nil {
		fmt.Printf("Failed to read: %v\n", err)
		return nil
	}

	gunzip, err := gzip.NewReader(fp)
	if err != nil {
		fmt.Printf("Failed to decompress: %v\n", err)
	}

	dec := gob.NewDecoder(gunzip)

	if err := dec.Decode(&activities); err != nil {
		fmt.Printf("Failed to decode: %v\n", err)
		return nil
	}

	return activities
}
