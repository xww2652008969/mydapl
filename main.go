package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Master struct {
	Author              string   `json:"Author"`
	Name                string   `json:"Name"`
	Punchline           string   `json:"Punchline"`
	Description         string   `json:"Description"`
	InternalName        string   `json:"InternalName"`
	ApplicableVersion   string   `json:"ApplicableVersion"`
	AssemblyVersion     string   `json:"AssemblyVersion"`
	RepoUrl             string   `json:"repoUrl"`
	IconUrl             string   `json:"IconUrl"`
	Changelog           string   `json:"Changelog"`
	Tags                []string `json:"Tags"`
	LoadPriority        int      `json:"LoadPriority"`
	DalamudApiLevel     int      `json:"DalamudApiLevel"`
	DownloadLinkInstall string   `json:"DownloadLinkInstall"`
}
type PluginConfig struct {
	Author            string   `json:"Author"`
	Name              string   `json:"Name"`
	InternalName      string   `json:"InternalName"`
	AssemblyVersion   string   `json:"AssemblyVersion"`
	Description       string   `json:"Description"`
	ApplicableVersion string   `json:"ApplicableVersion"`
	RepoUrl           string   `json:"RepoUrl"`
	Tags              []string `json:"Tags"`
	DalamudApiLevel   int      `json:"DalamudApiLevel"`
	LoadRequiredState int      `json:"LoadRequiredState"`
	LoadSync          bool     `json:"LoadSync"`
	CanUnloadAsync    bool     `json:"CanUnloadAsync"`
	LoadPriority      int      `json:"LoadPriority"`
	IconUrl           string   `json:"IconUrl"`
	Punchline         string   `json:"Punchline"`
	Changelog         string   `json:"Changelog"`
	AcceptsFeedback   bool     `json:"AcceptsFeedback"`
}

func main() {
	// 获取当前工作目录并处理错误
	root, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取工作目录失败: %v", err)
	}

	// 存储子目录列表
	var subDirs []string

	// 遍历目录树收集非隐藏子目录
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过隐藏目录（包括其子目录）
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}

			// 跳过根目录自身
			if path != root {
				subDirs = append(subDirs, path)
				return filepath.SkipDir // 仅处理直接子目录
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("目录遍历失败: %v", err)
	}

	// 读取主配置文件
	masterPath := filepath.Join(root, "pluginmaster.json")
	masterData, err := os.ReadFile(masterPath)
	if err != nil {
		log.Fatalf("读取主配置失败: %v", err)
	}

	// 解析主配置
	var plugins []Master
	if err := json.Unmarshal(masterData, &plugins); err != nil {
		log.Fatalf("解析主配置失败: %v", err)
	}

	// 处理每个子目录
	for _, dirPath := range subDirs {
		// 从完整路径中提取目录名
		dirName := filepath.Base(dirPath)
		configFile := filepath.Join(dirPath, dirName+".json")

		// 读取插件配置
		configData, err := os.ReadFile(configFile)
		if err != nil {
			log.Fatalf("读取插件配置 %s 失败: %v", configFile, err)
		}

		// 解析插件配置
		var cfg PluginConfig
		if err := json.Unmarshal(configData, &cfg); err != nil {
			log.Fatalf("解析插件配置 %s 失败: %v", configFile, err)
		}
		// 更新主配置信息
		found := false
		for i := range plugins {
			if plugins[i].Name == cfg.Name {
				plugins[i].AssemblyVersion = cfg.AssemblyVersion
				plugins[i].DalamudApiLevel = cfg.DalamudApiLevel
				found = true
				break
			}
		}
		if !found {
			log.Printf("警告: 插件 %s 未在主配置中找到", cfg.Name)
		}
	}

	// 生成带格式的JSON
	updatedData, err := json.MarshalIndent(plugins, "", "  ")
	if err != nil {
		log.Fatalf("JSON序列化失败: %v", err)
	}

	// 写入更新后的配置
	if err := os.WriteFile(masterPath, updatedData, 0644); err != nil {
		log.Fatalf("写入更新配置失败: %v", err)
	}

	fmt.Println("配置文件更新成功")
}
