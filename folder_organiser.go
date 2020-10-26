package folder_organiser

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"gopkg.in/djherbis/times.v1"
)

type FOFile struct {
	FullPath string
	CTime    time.Time
	FileName string
}

type FOFileMap struct {
	new      string
	old      string
	fileName string
}

// organise a folder by date. Will only organise files in the current dir
// dir is the directory you want to organise the files in
// layout is the go time.Format string e.g. 2006/01.
// folders will be created based on the split of the date e.g. 2006/01
func OrganiseByDate(dir string, layout string) {
	files := filesInDir(dir)
	cTimeFiles := getFileCreationDates(files, dir)
	cToN := currentToNewMap(cTimeFiles, layout, dir)
	organise(cToN)
}

func organise(fileMaps []FOFileMap) {
	for _, f := range fileMaps {

		_, existsErr := os.Stat(f.new)

		if os.IsNotExist(existsErr) {
			err := os.MkdirAll(f.new, 0775)

			if err != nil {
				log.Fatal(err)
			}
		}

		old := path.Join(f.old, f.fileName)
		new := path.Join(f.new, f.fileName)

		os.Rename(old, new)
		log.Println(old, " -> ", new)
	}
}

func filesInDir(dir string) []string {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		log.Fatal(err)
	}

	var filenames []string

	for _, file := range files {

		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}

		filenames = append(filenames, file.Name())
	}

	return filenames
}

func getFileCreationDates(files []string, dir string) []FOFile {
	foFiles := make([]FOFile, len(files))

	for i, file := range files {
		fullPath := path.Join(dir, file)
		cTime := getCreationDate(fullPath)

		foFiles[i] = FOFile{FullPath: dir, CTime: cTime, FileName: file}
	}

	return foFiles
}

func getCreationDate(file string) time.Time {
	stat, err := times.Stat(file)

	if err != nil {
		log.Fatal(err.Error())
	}

	if stat.HasBirthTime() {
		return stat.BirthTime()
	}

	return stat.ModTime()
}

func currentToNewMap(foFiles []FOFile, layout string, dir string) []FOFileMap {
	foFileMaps := make([]FOFileMap, len(foFiles))

	for i, file := range foFiles {
		formattedCTime := file.CTime.Format(layout)

		folders := []string{dir}
		folders = append(folders, strings.Split(formattedCTime, "/")...)

		newFullPath := path.Join(folders...)

		foFileMaps[i] = FOFileMap{old: file.FullPath, new: newFullPath, fileName: file.FileName}
	}

	return foFileMaps
}
