package server

import (
	"crypto/tls"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	// 实现了grpc库支持的各种凭证 该凭证封装了客户机需要的所有状态
	// 以便与服务器进行身份验证并进行各种断言
	// 如 关于客户机的身份 角色 是否授权 进行特定呼叫
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"net/http"

	"gRPC-Gateway/pkgs/util"
	pb "gRPC-Gateway/proto"
	// "gRPC-Gateway/pkgs/ui/data/swagger"
)

var (
	ServerPort     string
	CertServerName string
	CertPemPath    string
	CertKeyPath    string
	EndPoint       string
	tlsConfig      *tls.Config

	SwaggerDir string
)

func Serve() (err error) {
	log.Println(ServerPort)
	log.Println(CertServerName)
	log.Println(CertPemPath)
	log.Println(CertKeyPath)
	return nil
}

func Run() (err error) {
	EndPoint = ":" + ServerPort

	// net.Listen("tcp", EndPoint)用于监听本地的网络地址通知
	// 它的函数原型 func Listen(network, address string) (Listener, error)
	// -  参数 network必须传入tcp、tcp4、tcp6、unix、unixpacket
	// -  参数 address为空或为0则会自动选择一个端口号
	// 返回值为Listener 接口原型
	//  type Listener interface {
	//      Accept() (Conn, error) // 接受等待并将下一个连接返回给Listener
	//      Close() error          // 关闭Listener
	//      Addr() Addr            // 返回Listener的网络地址
	//  }
	// 最后net.Listen会返回一个监听器的接口,返回给接下来的动作

	conn, err := net.Listen("tcp", EndPoint)
	if err != nil {
		log.Printf("TCP Listen err:%v\n", err)
	}

	// 通过util.GetTLSConfig解析得到tls.Config
	// 传达给http.Server服务的TLSConfig配置项使用

	tlsConfig = util.GetTLSConfig(CertPemPath, CertKeyPath)
	srv := newServer(conn)
	log.Printf("gRPC and https listen on: %s\n", ServerPort)

	// 服务开始接受请求
	// srv.Serve(tls.NewListener(conn, tlsConfig))
	// 它是http.Server的方法 需要一个Listener作为参数
	if err = srv.Serve(util.NewTLSListener(conn, tlsConfig)); err != nil {
		log.Printf("ListenAndServe: %v\n", err)
	}

	return err
}

// 整个服务端的核心流转部分
func newServer(conn net.Listener) *http.Server {
	grpcServer := newGrpc()
	gwmux, err := newGrpcgateway()
	if err != nil {
		panic(err)
	}

	// http服务 http.NewServeMux 分配并返回一个新的ServeMux
	mux := http.NewServeMux()
	// mux.Handle 为给定模式注册处理程序
	mux.Handle("/", gwmux)

	// mux.HandleFunc("/swagger/", serveSwaggerFile)
	// serveSwaggerUI(mux)

	return &http.Server{
		Addr:      EndPoint,
		Handler:   util.GrpcAndGrpcgatewayAdapterHandlerFunc(grpcServer, mux),
		TLSConfig: tlsConfig,
	}
}

func newGrpc() *grpc.Server {
	// 程序采用的是HTT/2 需要支持TLS
	// 在启动grpc.NewServer前,要将认证的中间件注册进去
	// server.Run() 获取的tlsConfig仅能给HTTP/1.1使用
	//
	// 第一步 要创建grpc的TLS认证凭证
	// NewServerTLSFromFile 从输入证书文件和服务器的密钥文件构造TLS证书凭证

	creds, err := credentials.NewServerTLSFromFile(CertPemPath, CertKeyPath)
	if err != nil {
		log.Printf("Failed to create server TLS credentials %v", err)
		panic(err)
	}

	// 设置 grpc ServerOption
	opts := []grpc.ServerOption{
		// 原型 func Creds(c credentials.TransportCredentials) ServerOption
		// 返回ServerOption 为服务器连接设置凭据
		grpc.Creds(creds),
	}

	// 创建grpc服务端
	// 原型 func NewServer(opt ...ServerOption) *Server
	// 创建了一个没有注册服务的grpc服务端,还没有开始接受请求
	grpcServer := grpc.NewServer(opts...)

	// register grpc pb 注册grpc服务
	pb.RegisterHelloWorldServer(grpcServer, NewHelloService())

	return grpcServer
}

// 创建grpc-gateway关联组件
func newGrpcgateway() (http.Handler, error) {
	// 作为传入请求的顶级上下文,没有被注销 没有值 没有过期时间
	// 通常由主函数 初始化 测试使用

	ctx := context.Background()

	// 程序采用的是HTTPS 需要支持TLS
	// 在启动grpc.NewServer前,要将认证的中间件注册进去
	// 从客户机的输入证书文件构造TLS凭证

	dcreds, err := credentials.NewClientTLSFromFile(CertPemPath, CertServerName)
	if err != nil {
		log.Printf("Failed to create client TLS credentials %v", err)
		return nil, err
	}

	// grpc.WithTransportCredentials 配置一个连接级别的安全凭据
	// 如 TLS SSL 返回值为 type DialOption
	// grpc.DialOption 配置如何设置连接 内部由多个DialOption组成决定其设置连接的内容

	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}

	// 创建HTTP NewServeMux
	// runtime.NewServeMux 返回一个新的ServeMux 它的内部映射是空的
	// ServeMux 是 grpc-gateway 的一个请求多路复用器
	// 它将http请求与模式匹配 并调用相应的处理程序
	gwmux := runtime.NewServeMux()

	// register grpc-gateway pb 注册grpc-gateway逻辑 注册了HelloWorld这一个服务
	// 注册HelloWorld服务的HTTP Handler到grpc端点
	// -  上下文
	// -  grpc-gateway 请求多路复用器
	// -  服务网络地址
	// -  配置好的安全凭据
	if err := pb.RegisterHelloWorldHandlerFromEndpoint(ctx, gwmux, EndPoint, dopts); err != nil {
		return nil, err
	}
	return gwmux, nil
}

//func serveSwaggerFile(w http.ResponseWriter, r *http.Request) {
//      if ! strings.HasSuffix(r.URL.Path, "swagger.json") {
//        log.Printf("Not Found: %s", r.URL.Path)
//        http.NotFound(w, r)
//        return
//    }
//
//    p := strings.TrimPrefix(r.URL.Path, "/swagger/")
//    p = path.Join(SwaggerDir, p)
//
//    log.Printf("Serving swagger-file: %s", p)
//
//    http.ServeFile(w, r, p)
//}

//func serveSwaggerUI(mux *http.ServeMux) {
//    fileServer := http.FileServer(&assetfs.AssetFS{
//        Asset:    swagger.Asset,
//        AssetDir: swagger.AssetDir,
//        Prefix:   "third_party/swagger-ui",
//    })
//    prefix := "/swagger-ui/"
//    mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
//}
