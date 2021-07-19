package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync_tool/config"
)

func main()  {
	var pathConfig config.PathConf
	pathConfig.GetPathConf()
	//for i := range pathConfig.TargetFiles {
	//	//originPathAll = append(originPathAll, GetAllFiles(originPathAll, pathConfig.TargetFiles[i])...)
	//	GetAllFiles(originPathAll, pathConfig.TargetFiles[i])
	//}
	// 同步策略
	// 按照下面的格式获取源目录下所有文件（没有 . 后缀的就是文件夹）
	// 目录名 + 文件全名 + 文件修改时间
	// 比较源目录和目标目录的所有文件（把 A 目录的内容复制到 B 目录，A 就是源目录，B 就是目标目录）
	// 我们只需要去除掉没有发生变化的部分（按照上面规则得出相同字符串的文件），剩下的部分全部从源目录复制到目标目录即可
	// 注意，目标目录有，而源目录没有的文件，要予以在目标目录进行删除

	// 获取源目录的所有文件
	var originPathAll []string
	originPathAll = GetAllFiles(pathConfig.BaseModel)
	for i := range originPathAll {
		fmt.Println(originPathAll[i])
	}

	// 获取目标目录的所有文件
	// 这里想要搞一个二维数组，一直没搞出来
	fmt.Println("打印目标目录")
	targetPathLen := len(pathConfig.TargetFiles)
	var targetPathAll [][]string
	for i := range pathConfig.TargetFiles {
		tem := GetAllFiles( pathConfig.TargetFiles[i])
		if tem != nil {
			targetPathAll = append(targetPathAll[i], tem...)
		}
	}

	// 所有目录都是绝对地址，需要转换成相对地址，这样才能以源目录和目标目录 2 个不同的目录为根目录来比对文件的异同

	// 循环源目录，一一和目标目录进行比对，如果有文件名和修改时间相同，则不用复制，其他视情况需要复制就复制

	// 循环所有目标目录，一一个源目录就行比对，不在源目录中的，则删除

	// 保持控制台不关闭
	_,_ = fmt.Scanf("a")
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
