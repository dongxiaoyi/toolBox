package libs

import (
	"go.uber.org/zap"
	"os"
	"fmt"
	"reflect"
	"strings"
	"github.com/dongxiaoyi/toolBox/configs"
	"github.com/spf13/cobra"
	// closestmatch要使用github最新的代码，不要使用go mod的版本，go mod版本的代码有bug
	"github.com/schollz/closestmatch"
)

type T struct {}

var CmdMod = &cobra.Command{
	Use:   "mod",
	Short: "module",
	Long: `module.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := configs.NewLogger(true, true, true, true,  false,"console")
		if len(args) < 2 {
			logger.Info("Args must >= 1!")
			fmt.Println("###################Error Exist##################")
			os.Exit(1)
		} else {
			// 此处根据不同的parser类型（file、str等）解析参数后分发具体mod操作逻辑
			modName := args[0] // 第一个参数必须是模块名称
			modActionName := args[1] // 第二个参数必须是模块的action名称
			// 第三个参数可选：from=str from=string，不填为从配置文件查找模块配置，其他为从具体文件查找模块配置
			bagSizes := []int{5}

			cmFilter := closestmatch.New(args, bagSizes)
			modFrom := cmFilter.Closest("from=")

			var modExecContent map[string]string  // 最终解析后传递给
			switch modFrom {
			case "from=str", "from=string":
				// 解析命令行
				modExecContent = cmdStr2Content(args[2:])
			default:
				// 默认为file类型，解析file配置内容(file类型的需要知道具体的action配置块)
				modExecContent = file2Content(modName, modActionName, strings.TrimLeft(modFrom, "from="), logger)
			}
			// 将content传递给具体的插件处理
			execAction(modName, modActionName, modExecContent)
		}
	},
}

func cmdStr2Content(cmdStr []string) map[string]string {
	// 命令行参数以换行符分割返回
	cmdStr = removeStr(cmdStr, "from=str")
	cmdStr = removeStr(cmdStr, "from=string")
	return map[string]string{"cmdContent": strings.Join(cmdStr, "\n")}
}

func removeStr(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func file2Content(modName, modActionName, filePath string, logger *zap.SugaredLogger) map[string]string {
	// 根据action配置块返回配置文本
	// modName.modActionName 为具体模块配置项section，且为小写
	sectionName := strings.ToLower(modName) + "." + strings.ToLower(modActionName)
	if filePath == "" {
		filePath = "configs/actions.ini"
	}
	config := configs.NewActionConfig(filePath)
	sectionConf, err := config.LoadConfigBySectionName(sectionName)
	if err != nil {
		logger.Error(err)
		os.Exit(2)
	}
	return sectionConf
}

func execAction(modName, modActionName string, modExecContent map[string]string) {
	// 根据modName和modActionName动态调用执行函数
	// 函数命名规则：modName与modActionName拼接，但是为驼峰命名
	funcName := upperFirstChar(strings.ToLower(modName)) + upperFirstChar(strings.ToLower(modActionName))
	// 函数名称方式调用逻辑函数
	t := T{}
	in := []reflect.Value{ reflect.ValueOf(modExecContent) }
	reflect.ValueOf(t).MethodByName(funcName).Call(in)
}

// 首字母大写
func upperFirstChar(str string) string {
	if len(str) > 0 {
		runeList := []rune(str)
		if int(runeList[0]) >= 97 && int(runeList[0]) <= 122 {
			runeList[0] = rune(int(runeList[0]) - 32)
			str = string(runeList)
		}
	}
	return str
}
