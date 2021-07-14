package main

import (
	"fmt"
	"sync_tool/config"
)

func main()  {
	var pathConfig config.PathConf
	pathConfig.GetPathConf()
	fmt.Println("path:" + pathConfig.BaseModel)

	// 保持控制台不关闭
	_,_ = fmt.Scanf("a")
}
