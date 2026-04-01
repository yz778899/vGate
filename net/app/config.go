package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type config struct {
	// GateConfig
	// LogConfig
	Gate struct {
		WsPath        string `yaml:"WsPath"`        //路径Config
		WsPort        int    `yaml:"WsPort"`        //端口
		SecretKey     string `yaml:"secretKey"`     //安全密钥
		HeartbeatTime int    `yaml:"HeartbeatTime"` //心跳时间
		ReadOverTime  int    `yaml:"ReadOverTime"`  //读超时

	} `yaml:"gate"`
	Log struct {
		Open bool   `yaml:"open"`
		Path string `yaml:"path"`
	} `yaml:"log"`
}

func getConfig(path string) *config {

	cwd, _ := os.Getwd()
	if folder, err := FindFolderUpward(cwd, "config"); err == nil {
		path = folder + "\\" + path
	}

	// 读取 YAML 文件
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("读取文件失败: %v", err)
	}

	//反序列化到结构体
	var config config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("解析 YAML 失败: %v", err)
	} else {
		//fmt.Printf("config = %#v", config)
	}
	return &config
}

// FindFolderUpward 向上查找文件夹
func FindFolderUpward(startPath, targetFolder string) (string, error) {
	currentPath, err := filepath.Abs(startPath)
	if err != nil {
		return "", err
	}

	// 逐级向上查找
	for {
		// 检查当前路径下是否存在目标文件夹
		targetPath := filepath.Join(currentPath, targetFolder)
		if info, err := os.Stat(targetPath); err == nil && info.IsDir() {
			return targetPath, nil
		}

		// 到达根目录，停止查找
		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			break
		}
		currentPath = parentPath
	}

	return "", fmt.Errorf("folder %s not found", targetFolder)
}
