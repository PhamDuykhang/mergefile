package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

func main() {
	var targetDir = flag.String("target", "", "the folder where you place your file to megere")
	var resultDir = flag.String("result", "", "the folder where you put your file")
	flag.Parse()
	files, err := ioutil.ReadDir(fmt.Sprintf("./%s", *targetDir))
	if err != nil {
		log.Fatal(err)
	}
	fileMap := make(map[string][]string)
	for _, f := range files {
		name := f.Name()[0:strings.LastIndex(f.Name(), "-")]
		ext := f.Name()[strings.LastIndex(f.Name(), "-")+1:]
		l, ok := fileMap[name]
		if ok {
			l = append(l, ext)
			fileMap[name] = l
		} else {
			fileMap[name] = []string{ext}
		}

	}
	var wg sync.WaitGroup
	for k, v := range fileMap {
		wg.Add(1)
		go mergeFile(&wg, *resultDir, k, v)
	}
	wg.Wait()

}
func mergeFile(wg *sync.WaitGroup, dir string, fileName string, page []string) {
	defer wg.Done()
	f, err := os.OpenFile(fmt.Sprintf("./%s/%s.txt", dir, fileName), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		logrus.Error(err)

	}
	for _, fileElement := range page {
		b, err := ioutil.ReadFile(fmt.Sprintf("./data/1/%s-%s", fileName, fileElement))
		if err != nil {
			logrus.Warn(err)
		}
		f.WriteString(string(b))
	}
	defer f.Close()
}
