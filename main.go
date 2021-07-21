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
	var originPathAll []config.DirInfo
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
				if originPathAll[o].Path == targetPathAll[i].Dirs[j].Path {
					inOrigin = true
					break
				}
			}
			if !inOrigin {
				// 不在源目录中，删除
				redundantPath := targetPathAll[i].BaseDir + targetPathAll[i].Dirs[j].Path
				fmt.Println(chalk.Red, dateNow + "，删除目录：" + redundantPath, chalk.ResetColor)
				_ = os.RemoveAll(redundantPath)
			}
		}
	}
	// 循环源目录，一一和目标目录进行比对，如果有文件名和修改时间相同，则不用复制，其他视情况需要复制就复制
	for originIndex := range originPathAll{
		originFile := originPathAll[originIndex]
		for targetIndex := range targetPathAll {
			// 源目录，是否存在于目标目录中
			isExistTargetPath := false
			updatedTheSame := true
			for targetPathIndex := range targetPathAll[targetIndex].Dirs {
				if originFile.Path == targetPathAll[targetIndex].Dirs[targetPathIndex].Path {
					isExistTargetPath = true
					if originFile.UpdatedAt != targetPathAll[targetIndex].Dirs[targetPathIndex].UpdatedAt {
						updatedTheSame = false
					}
					break
				}
			}
			targetPath := targetPathAll[targetIndex].BaseDir + originFile.Path
			originAbsolutePath := pathConfig.BaseModel + originFile.Path
			if !isExistTargetPath {
				if originFile.IsDir {
					fmt.Println(chalk.Green, dateNow + "，创建文件夹：" + targetPath, chalk.ResetColor)
					_ = os.Mkdir(targetPath, 777)
				} else {
					fmt.Println(chalk.Green, dateNow + "，新建文件：" + targetPath, chalk.ResetColor)
					copyFile(originAbsolutePath, targetPath)
					// 改变目标目录的修改时间，因为是以比对修改时间是否一致来判断 2 个文件是否一样
					_ = os.Chtimes(targetPath, originFile.UpdatedAt, originFile.UpdatedAt)
				}
			} else {
				// 存在相同目录，比较修改时间
				if !originFile.IsDir {
					// 文件夹不去管，只管文件
					if !updatedTheSame {
						// 更新时间不同，再复制
						fmt.Println(chalk.Yellow, dateNow + "，覆盖文件：" + targetPath, chalk.ResetColor)
						copyFile(originAbsolutePath, targetPath)
						_ = os.Chtimes(targetPath, originFile.UpdatedAt, originFile.UpdatedAt)
					}
				}
			}
		}
	}
}

func GetAllFiles(dir string) []config.DirInfo {
	// 获取一个目录下的所有文件
	// 在这里不仅要拿到文件名，还要拿到是否是目录、修改时间等信息，方便后面做处理
	var originPathAll [] config.DirInfo
	fileInfoList, _ := ioutil.ReadDir(dir)
	for i := range fileInfoList {
		var dirInfo config.DirInfo
		dirInfo.Path = dir + "\\" + fileInfoList[i].Name()
		dirInfo.IsDir = fileInfoList[i].IsDir()
		dirInfo.UpdatedAt = fileInfoList[i].ModTime()
		originPathAll = append(originPathAll, dirInfo)
		if fileInfoList[i].IsDir() {
			originPathAll = append(GetAllFiles(dirInfo.Path), originPathAll...)
		}
	}
	return originPathAll
}

func AbPathToRelativePath(dirs []config.DirInfo, baseDir string) []config.DirInfo {
	// 绝对地址，变成相对地址，且将文件夹放到前面
	var result []config.DirInfo
	for i := range dirs {
		dirs[i].Path = strings.Replace(dirs[i].Path, baseDir, "", 1)
		if dirs[i].IsDir {
			// 如果是文件夹，放在最前面
			result = append([]config.DirInfo{dirs[i]}, result...)
		} else {
			// 如果是文件的话，放到最后面
			result = append(result, []config.DirInfo{dirs[i]}...)
		}
	}
	return result
}

func GetAllFileWithRelativePath(dir string) []config.DirInfo {
	// 以相对地址的方式获得所有目录
	dirs := GetAllFiles(dir)
	dirs = AbPathToRelativePath(dirs, dir)
	return dirs
}

/**
测试方法
 */
func copyFile(baseDir string, targetDir string) {
	input, _ := ioutil.ReadFile(baseDir)
	_ = ioutil.WriteFile(targetDir, input, 0644)
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
