package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type typeFoL int
const (
	First typeFoL = 0
	Last typeFoL = 1
)

type ContentLookup struct {
	Path       string
	LineNumber int
}


// FIXME - Why in the world am I using globals for this. It's a 90 minute hack that's why

var contentMap map[ContentLookup]string
var numLines map[string]int
var whichFile map[int]string
var howManyFiles int

var commonHeader map[int]string
var commonFooter map[int]string

func main() {

	var files []string
	var loopInt int
	var templateHeader string
	var templateFooter string

	contentMap = make(map[ContentLookup]string)
	numLines = make(map[string]int)
	whichFile = make(map[int]string)

	commonHeader = make(map[int]string)
	commonFooter = make(map[int]string)

	startHere := os.Getenv("START_PATH")

	err := filepath.Walk(startHere, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".html" {
			if strings.Contains(path, "copies") == false {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	processFileList(files)

	fmt.Println("First Lines")
	compareNLines(45, First)
	fmt.Println(" ====== ")
	fmt.Println(" Last Lines")
	compareNLines(35, Last)
	fmt.Println(" ")

	//
	insertFixedLines()

	templateHeader = ""
	templateFooter = ""
	for loopInt = 0; loopInt < 50; loopInt++ {
		if thisHeaderLine, headerok := commonHeader[loopInt]; headerok {
			templateHeader = templateHeader + thisHeaderLine + "\n"
		}
		if thisFooterLine, footerok := commonFooter[loopInt]; footerok {
			templateFooter = thisFooterLine + "\n" + templateFooter
		}
	}

	errwh := ioutil.WriteFile("common_header.inc", []byte(templateHeader), 0644)
	if errwh != nil {
		fmt.Println("ERR - could not write header file")
	}
	errwf := ioutil.WriteFile("common_footer.inc", []byte(templateFooter), 0644)
	if errwf != nil {
		fmt.Println("ERR - could not write footer file")
	}
}

func findnumunique(xyz []uint32) int {
	var retval int
	var countthis map[uint32]int
	var tempint int

	countthis = make(map[uint32]int)

	tempint = 0;	
	for _, row := range xyz {
		tempint++
		countthis[row]++
	}

	retval = 0
	for k := range countthis {
		if k != 32489 {
			retval++
		}
	}

	return retval
}

func compareNLines(nZ int, firstorlast typeFoL) {
	var numuniques int
	var commonline string

	var tempstr string
	var outputstr string

	outputstr = ""
	hashitArray := make([]uint32, howManyFiles)

	for yy := 0; yy < nZ; yy++ {
		commonline = ""
		for zz := 0; zz < howManyFiles; zz++ {
			tempstr, _ = getLine(whichFile[zz], yy, firstorlast)
			hashitArray[zz] = hashit(strings.ToLower(tempstr))
		}
		numuniques = findnumunique(hashitArray)
		if numuniques == 1 {
			commonline, _ = getLine(whichFile[0], yy, firstorlast)
			if firstorlast == First {
				commonHeader[yy] = commonline
			}
			if firstorlast == Last {
				commonFooter[yy] = commonline
			}
		}

		tempstr = fmt.Sprintf("(%02d) %d = %s\n", yy, numuniques, commonline)
		if firstorlast == Last {
		outputstr = tempstr + outputstr
		} else {
		outputstr = outputstr + tempstr
		}
		//fmt.Printf("(%02d) %d = %s\n", yy, numuniques, commonline)
	}
	fmt.Printf("%s\n", outputstr)
}

func getLine(fname string, filepos int, whichdirection typeFoL) (retstr string, linenum int) {
	var looppos int

	looppos = filepos
	if whichdirection == Last {
		looppos = (numLines[fname] - 1) - filepos
	}

	retstr = contentMap[ContentLookup{fname, looppos}]
	linenum = looppos
	return retstr, linenum
}

func processFileList(fileList []string) {
	var xy int
	xy = 0
	for _, singleFileName := range fileList {
		processEachFile(singleFileName)
		whichFile[xy] = singleFileName
		xy++
	}
	howManyFiles = xy
}

func processEachFile(theFile string) {
	var xy int
	file, _ := os.Open(theFile)
	fscanner := bufio.NewScanner(file)
	xy = 0
	for fscanner.Scan() {
		contentMap[ContentLookup{theFile, xy}] = strings.TrimSpace(fscanner.Text())
		xy++
	}
	numLines[theFile] = xy
}

func hashit(s string) uint32 {
        h := fnv.New32a()
        h.Write([]byte(s))
        return h.Sum32()
}


