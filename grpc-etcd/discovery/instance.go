package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Server struct {
	Project string `json:"project"`
	Name    string `json:"name"`
	Addr    string `json:"addr"`
	Version string `json:"version"`
	Weight  int    `json:"weight"`
}

// BuildPrefix key前缀
func BuildPrefix(server Server) string {
	if server.Version != "" && server.Project != "" {
		return fmt.Sprintf("%s/%s/%s/", server.Project, server.Name, server.Version)
	}
	if server.Version == "" {
		return fmt.Sprintf("%s/%s/", server.Project, server.Name)
	}
	if server.Project == "" {
		return fmt.Sprintf("%s/%s/", server.Name, server.Version)
	}
	return fmt.Sprintf("%s/", server.Name)
}

// BuildRegisterPath 拼装key
func BuildRegisterPath(server Server) string {
	return fmt.Sprintf("/%s%s", BuildPrefix(server), server.Addr)
}

// SplitPath 切割路径（key）
func SplitPath(path string) (Server, error) {
	server := Server{}
	strArr := strings.Split(path, "/")
	if len(strArr) == 0 {
		return server, errors.New("invalid path")
	}

	// 从1开始，因为开头有一个斜杠
	server.Project = strArr[1]
	server.Name = strArr[2]
	server.Version = strArr[3]
	server.Addr = strArr[4]
	return server, nil
}

// ParseValue 解析 value 到 Server
func ParseValue(value []byte) (Server, error) {
	server := Server{}
	if err := json.Unmarshal(value, &server); err != nil {
		return server, err
	}
	return server, nil
}
