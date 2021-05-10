package main

import (
	"github.com/dongxiaoyi/toolBox/configs"
	"github.com/dongxiaoyi/toolBox/libs"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	Logger *zap.SugaredLogger
	rootCmd = &cobra.Command{Use: "toolBox"}
)

func init() {
	Logger = configs.NewLogger(false, true, true, true, false, "console")

	rootCmd.AddCommand(libs.CmdMod)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		Logger.Error(err)
	}
}