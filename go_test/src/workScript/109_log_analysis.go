package fileManage

import (
	"fmt"
	"io/ioutil"
	"strings"
	"regexp"
	"os"
	"github.com/mholt/archiver"
	"log"
	"bufio"
)

// 公共变量定义
var srcPath string = `C:\Users\Administrator\Downloads`
var desPath string = `C:\Users\Administrator\Desktop\109_log`

//迁移 从日志服务器下载下来的日志文件文件
func moveTarget109LogFile() []string {
	//提取目标文件的文件名 [*text_reocord*]
	textRecordsFiles := findFuzzyFile(srcPath, "text_reocord")
	if len(textRecordsFiles) == 0 {
		textRecordsFiles = findFuzzyFile(desPath, "text_reocord")
	}

	var desFileNames []string
	for _, fileName := range textRecordsFiles {
		desFileNames = append(desFileNames, fileName)
		err := os.Rename(srcPath + `\` + fileName, desPath + `\` + fileName) //使用Rename 来实现移动文件
		if err != nil {
			fmt.Println("[ERROR]从Download 目录移动文件失败")
			//fmt.Printf("file %s move failed.\n srcPath:%s, desPath:%s.\n err is:%s\n", fileName, srcPath, desPath, err)
		} else {
			//fmt.Printf("移动文件:%s!", fileName)
		}
	}
	return desFileNames
}

func getTargetDirectoryFileNamesPath(directoryName string) []string {
	var  fileNames []string
	if len(directoryName) == 0 {
		fmt.Printf("[ERROR] 目录字符为空!")
		return fileNames
	}

	files, err := ioutil.ReadDir(directoryName)
	//fmt.Printf("显示目录：%s 文件，文件数目:%d\n", directoryName, len(files))
	if err != nil {
		fmt.Println(err)
	} else {
		for _, filename := range files {
			if filename.IsDir() {
				//fmt.Printf("目录名:%s\n", filename.Name())
			} else {
				fileNames = append(fileNames, filename.Name())
				//fmt.Printf("文件名:%s\n", filename.Name())
			}
		}
	}

	return fileNames
}

// findFuzzyFile 模糊查找文件
func findFuzzyFile(inputFilename string, regexName string) []string {
	if len(inputFilename) == 0 {
		fmt.Println("[ERROR] 文件名为空！")
	}
	var hitFiles []string
	filePaths := getTargetDirectoryFileNamesPath(inputFilename)
	for _, path := range filePaths {
		if strings.Contains(path, regexName) {
			hitFiles = append(hitFiles, path)
		}
	}

	return hitFiles
}

func getLogFileTime(inputFileName string) string {
	//匹配时间 xxxx-xx-xx
	reg := regexp.MustCompile(`[0-9]{4}-[0-9]{2}-[0-9]{2}`)
	return reg.FindString(inputFileName)
}

//处理 109组日志
func Handle109Log() {
	desFileNames := moveTarget109LogFile()

	//埋点，新建目录用于汇总 所有日期的日志文件
	var LogFor109Dir string = desPath + `\` + "109_log"
	os.Mkdir(LogFor109Dir, 077)

	//迁移文件
	for _, fileName := range desFileNames {
		//提取日志文件中的日期值，形如 xxxx-xx-xx
		logTime := getLogFileTime(fileName)

		// 创建日期目录
		logTimePath := desPath + `\` + logTime
		err := os.Mkdir(logTimePath, 077)
		if err != nil {
			fmt.Println(err)
		}

		//复制 *.tar.gz 文件到 对应日期目录中
		os.Link(desPath + `\` + fileName, desPath + `\` + logTime + `\` + fileName)

		//压缩文件 *.tar.gz 并删除 日期目录下的*.tar.gz文件
		archiver.TarGz.Open(desPath + `\` + logTime + `\` + fileName, desPath + `\` + logTime)
		os.Remove(desPath + `\` + logTime + `\` + fileName)

		//最终归档的每日日志文件名
		var finalLogFile[2] string
		finalLogFile[0] = desPath + `\` + logTime + `\` + `text_reocords.log.` + logTime + `_00~11`
		finalLogFile[1] = desPath + `\` + logTime + `\` + `text_reocords.log.` + logTime + `_12~23`

		//组织文件内容，每12个文件合并成一个文件
		logFiles := getTargetDirectoryFileNamesPath(desPath + `\` + logTime)
		var indexFileName int = 0
		var logFilesContents[] string
		for i, fileName := range logFiles {
			fmt.Printf("index:%d, target file name:%s\n", i, fileName)

			filesContents, err := ReadLineByLineFromLocalFile(desPath + `\` + logTime + `\` + fileName)
			if err == true {
				for _, content := range filesContents {
					logFilesContents = append(logFilesContents, content)
				}
			} else {
				fmt.Printf("read file:%s failed\n", fileName)
			}


			if i == 11 || i == 23 {
				retSaveFile := SaveToLocalFile(finalLogFile[indexFileName] , logFilesContents)
				if retSaveFile != true {
					fmt.Println("save multi file to one file failed.")
				}
				indexFileName++
				logFilesContents = logFilesContents[:0]
			}
		}

		//赋值每日汇总的文件到 ../109_log 目录
		os.Link(finalLogFile[0], LogFor109Dir + `\` + `text_reocords.log.` + logTime + `_00~11`)
		os.Link(finalLogFile[1], LogFor109Dir + `\` + `text_reocords.log.` + logTime + `_12~23`)
	}
}

func ReadLineByLineFromLocalFile(filename string) ([]string, bool) {
	var fileContents []string
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return fileContents, false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileContents = append(fileContents, scanner.Text())
		fileContents = append(fileContents, "\n")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return fileContents, true
}

func SaveToLocalFile(outPutFilename string, content []string) bool {
	// open output file
	fo, err := os.Create(outPutFilename)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	// write a chunk
	i := 0
	var buf string
	for ; i < len(content); i++ {
		buf = content[i]
		if _, err := fo.Write([]byte(buf)); err != nil {
			panic(err)
		}
	}

	return true
}