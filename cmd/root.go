package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd表示在没有任何子命令的情况下的基本命令
var rootCmd = &cobra.Command{
	// Use : Command的用法，Use是一个行用法消息
	Use: "grpc",
	// Short : 是help命令输出中显示的简短描述
	Short: "Run the gRPC & gRPC-Gateway hello-world server",
	// Run : 运行 典型的实际工作功能 大多数命令只会实现这一点
	// 另外还有 PreRun PreRunE PostRun PostRunE 等不同时期的运行命令,比较少用
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
