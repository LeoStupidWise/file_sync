package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/ttacon/chalk"
	"io/ioutil"
	"os"
	"strings"
	"sync_tool/config"
	"time"
)

func main()  {
	// 保持控制台不关闭
	//_,_ = fmt.Scanf("a")
	fmt.Println(chalk.Green, GetTimeNow() + "，开始同步监听", chalk.ResetColor)
	var pathConfig config.PathConf
	pathConfig.GetPathConf()
	c := cron.New()
	_,_ = c.AddFunc(pathConfig.Cron, doCopy)
	c.Start()
	select {}
}

func GetTimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func doCopy()  {
	// golang 中必须使用下面这个时间才行
	dateNow := GetTimeNow()
	var pathConfig config.PathConf
	pathConfig.GetPathConf()
	// 同步策略
	// 按照下面的格式获取源目录下所有文件（没有 . 后缀的就是文件夹）
	// 目录名 + 文件全名 + 文件修改时间
	// 比较源目录和目标目录的所有文件（把 A 目录的内容复制到 B 目录，A 就是源目录，B 就是目标目录）
	// 我们只需要去除掉没有发生变化的部分（按照上面规则得出相同字符串的文件），剩下的部分全部从源目录复制到目标目录即可
	// 注意，目标目录有，而源目录没有的文件，要予以在目标目录进行删除

	// 获取源目录的所有文件
	var originPathAll []string
	originPathAll = GetAllFileWithRelativePath(pathConfig.BaseModel)

	// 获取目标目录的所有文件
	// 这里保存目标目录的详情
	var targetPathAll [] config.TargetDir
	for i := range pathConfig.TargetFiles {
		tem := GetAllFileWithRelativePath(pathConfig.TargetFiles[i])
		var targetDir config.TargetDir
		targetDir.BaseDir = pathConfig.TargetFiles[i]
		if tem != nil {
			targetDir.Dirs = append(targetDir.Dirs, tem...)
		}
		targetPathAll = append(targetPathAll, targetDir)
	}
	// 循环所有目标目录，一一个源目录就行比对，不在源目录中的，则删除
	for i := range targetPathAll {
		for j := range targetPathAll[i].Dirs {
			var inOrigin = false
			for o := range originPathAll{
				if originPathAll[o] == targetPathAll[i].Dirs[j] {
					inOrigin = true
					break
				}
			}
			if !inOrigin {
				// 不在源目录中，删除
				redundantPath := targetPathAll[i].BaseDir + targetPathAll[i].Dirs[j]
				fmt.Println(chalk.Red, dateNow + "，删除目录：" + redundantPath, chalk.ResetColor)
				_ = os.RemoveAll(redundantPath)
			}
		}
	}
	// 循环源目录，一一和目标目录进行比对，如果有文件名和修改时间相同，则不用复制，其他视情况需要复制就复制
	for originIndex := range originPathAll{
		originPath := originPathAll[originIndex]
		for targetIndex := range targetPathAll {
			// 源目录，是否存在于目标目录中
			isExistTargetPath := false
			for targetPathIndex := range targetPathAll[targetIndex].Dirs {
				if originPath == targetPathAll[targetIndex].Dirs[targetPathIndex] {
					isExistTargetPath = true
				}
			}
			targetPath := targetPathAll[targetIndex].BaseDir + originPath
			if !isExistTargetPath {
				fmt.Println(chalk.Green, dateNow + "，新增目录：" + targetPath, chalk.ResetColor)
			} else {
				// 存在相同目录，比较修改时间
			}
		}
	}
}

func GetAllFiles(dir string) []string {
	// 获取一个目录下的所有文件
	var originPathAll [] string
	fileInfoList, _ := ioutil.ReadDir(dir)
	for i := range fileInfoList {
		dirNow := dir + "\\" + fileInfoList[i].Name()
		originPathAll = append(originPathAll, dirNow)
		if fileInfoList[i].IsDir() {
			originPathAll = append(GetAllFiles(dirNow), originPathAll...)
		}
	}
	return originPathAll
}

func AbPathToRelativePath(dirs []string, baseDir string) []string {
	// 绝对地址，变成相对地址
	for i := range dirs {
		dirs[i] = strings.Replace(dirs[i], baseDir, "", 1)
	}
	return dirs
}

func GetAllFileWithRelativePath(dir string) []string {
	// 以相对地址的方式获得所有目录
	dirs := GetAllFiles(dir)
	dirs = AbPathToRelativePath(dirs, dir)
	return dirs
}

/**
测试方法
 */
func copyFile(baseDir string, targetDir string) {
	fileInfoList, _ := ioutil.ReadDir(baseDir)
	for i := range fileInfoList {
		originDirNow := baseDir+"\\"+fileInfoList[i].Name()
		targetDirNow := targetDir+"\\"+fileInfoList[i].Name()
		copyFile(originDirNow, targetDirNow)
		// 如果是目录
		if fileInfoList[i].IsDir() {
			// 目标地址不存在这个目录
			targetFileExist, _ := PathExists(targetDirNow)
			if !targetFileExist {
				// 目标目录中不存在对应文件夹，先创建这个文件夹
				err := os.MkdirAll(targetDir, 777)
				if err != nil {
					fmt.Println( err)
				}
			}
		} else {
			//
		}
		//fmt.Println(reflect.TypeOf(fileInfoList[i]))
		//fmt.Println(fileInfoList[i].Name())
		//fmt.Println(fileInfoList[i].ModTime().Unix())
	}
	fmt.Println("=================" + targetDir)
}

/**
目录是否存在
 */
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
