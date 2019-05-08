package util

import (
	"net/http"
	"strings"

	"google.golang.org/grpc"
)

// GrpcHandlerFunc函数用于判断请求是来源于Rpc客户端 或者 Restful Api请求
// 根据不同的请求注册不同的ServeHTTP服务
func GrpcAndGrpcgatewayAdapterHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	if otherHandler == nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			grpcServer.ServeHTTP(w, r)
		})
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// r.ProtoMajor == 2 表示请求必须基于HTTP/2
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}
