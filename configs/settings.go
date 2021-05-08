package configs

import (
	"github.com/Unknwon/goconfig"
	"github.com/dongxiaoyi/toolBox/internal"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)
/*
以下为工具集的配置
*/

type ActionConfig struct {
	Filename string
	Sections *goconfig.ConfigFile
}

func NewActionConfig(filename string) *ActionConfig {
	logger := NewLogger(true, true, true, true, false, "console")

	if !strings.HasPrefix(filename, "/") {
		filename = filepath.Join(internal.AbsPath(), filename)
	}

	// 手动读文件，修正（单section的多行文本配置转单行）后使用goconfig.LoadFromData()
	// collection所有section
	configContent, err := readFile(filename)
	if err != nil {
		logger.Error(err)
		os.Exit(3)
	}
	configContentFix := configFix(configContent, logger)

	config, err := goconfig.LoadFromData([]byte(configContentFix))
	if err != nil {
		logger.Error(err)
		os.Exit(4)
	}

	return &ActionConfig{
		Filename: filename,
		Sections: config,
	}

}

func (c *ActionConfig) LoadConfigBySectionName(sectionName string) (map[string]string, error) {
	return c.Sections.GetSection(sectionName)
}

func readFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fd, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(fd), nil
}

// 修正配置格式
func configFix(content string, logger *zap.SugaredLogger) string {
	r, _ := regexp.Compile("\"\"\"")
	subMatchIndex := r.FindAllStringSubmatchIndex(content, -1)
	subMatchIndexLength := len(subMatchIndex)
	if subMatchIndexLength & 0x1 != 0 {
		logger.Error("配置文件多行匹配失败!")
		os.Exit(5)
	}

	newMatchIndex := make([][]int, 0)

	for i := 2; i <= len(subMatchIndex); i+=2 {
		indexList := []int{subMatchIndex[i-2][0], subMatchIndex[i-1][1]}
		newMatchIndex = append(newMatchIndex, indexList)
	}

	// 原配置文件与待替换配置索引差集
	oldContentSplitIndex := make([][]int, 1)
	oldContentSplitIndex[0] = []int{0}

	for i := 0; i < len(newMatchIndex); i++ {
		for j := 0; j < 2; j++ {
			if j == 0 {
				oldContentSplitIndex[i] = append(oldContentSplitIndex[i], newMatchIndex[i][j])
			} else {
				oldContentSplitIndex = append(oldContentSplitIndex, []int{newMatchIndex[i][j]})
			}
		}
	}
	oldContentSplitIndex[len(oldContentSplitIndex)-1] = append(oldContentSplitIndex[len(oldContentSplitIndex)-1], len(content))

	// 原配置文件与待替换配置差集string
	oldContentSplitString := make([]string, 0)
	for i := 0; i < len(oldContentSplitIndex); i++ {
		oldContentSplitString = append(oldContentSplitString, content[oldContentSplitIndex[i][0]:oldContentSplitIndex[i][1]])
	}

	// 生成待替换配置："""多行变单行
	fixContent := make([]string, 0)
	for i := 0; i < len(newMatchIndex); i++ {
		str := content[newMatchIndex[i][0]:newMatchIndex[i][1]]
		str = strings.Trim(str, "\"\"\"")
		str = strings.Replace(str, "\n", " ", -1)
		fixContent = append(fixContent, str)
	}

	// 拼接为常规ini格式配置文本
	step := 0
	for i := 0; i < len(fixContent); i++ {
		oldContentSplitString = append(oldContentSplitString[:i+step+1], append([]string{fixContent[i]}, oldContentSplitString[i+step+1:]...)...)
		step += 1
	}

	return strings.Join(oldContentSplitString, "")
}