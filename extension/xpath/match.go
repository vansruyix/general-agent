// Package xpath
// @author: fengyi
// @date: 2024/7/18
// @note:
package xpath

import (
	"general-agent/framework/xcache"
	"strings"

	"github.com/gobwas/glob"
)

// UriMatch 路径匹配器，支持*通配符
func UriMatch(sourcePath, targetPath string) bool {
	if sourcePath == targetPath {
		return true
	} else if strings.Contains(sourcePath, "*") {
		// 通配符比对
		// 为每个通配符权限缓存一个匹配器
		_, err := xcache.Instance.Load("path.match."+sourcePath, func() (any, error) {
			return glob.MustCompile(sourcePath), nil
		})
		if err != nil {
			return false
		}
		// m, ok := matcher.(glob.Glob)
		// return ok && m.Match(targetPath)
	}
	return false
}
