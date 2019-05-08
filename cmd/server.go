package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"gRPC-Gateway/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the gRPC & gRPC-Gateway hello-world server",
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recover error : %v", err)
			}
		}()

		server.Run()
	},
}

func init() {
	// 一般,需要在init()函数中定义flags和处理配置
	// serverCmd.Flags().StringVarP() 定义了一个flag 值存储在&server.ServerPort中
	// 长命令为--port，短命令为-p 默认值为50052 命令的描述为server port
	// 这一种调用方式成为 Local Flags
	serverCmd.Flags().StringVarP(&server.ServerPort,
		"port", "p", "50052",
		"server port")
	serverCmd.Flags().StringVarP(&server.CertPemPath,
		"cert-pem", "", "./conf/certs/server.pem",
		"cert-pem path")
	serverCmd.Flags().StringVarP(&server.CertKeyPath,
		"cert-key", "", "./conf/certs/server.key",
		"cert-key path")
	serverCmd.Flags().StringVarP(&server.CertServerName,
		"cert-server-name", "", "grpc server name",
		"server's hostname")

	//serverCmd.Flags().StringVarP(&server.SwaggerDir, "swagger-dir", "", "proto", "path to the directory which contains swagger definitions")

	// AddCommand向这父命令(rootCmd)添加一个或多个命令
	rootCmd.AddCommand(serverCmd)
}
