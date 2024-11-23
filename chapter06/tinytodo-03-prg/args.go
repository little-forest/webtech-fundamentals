package main

import (
	"flag"
	"os"
	"strconv"
)

const defaultPort = 8080

// 環境変数またはコマンドライン引き数からポート番号を取得する。
func getPortNumber() int {
	var port int
	portStr := os.Getenv("PORT")
	if portStr == "" {
		p := flag.Int("p", defaultPort, "Port number")
		flag.Parse()
		port = *p
	} else {
		port, _ = strconv.Atoi(portStr)
	}
	return port
}
