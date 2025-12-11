// Copyright 2025 vistone. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package moduleinit

import "fmt"

// getDisplayValue 获取显示值，如果为空则返回默认值
func getDisplayValue(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// getBrowserList 获取浏览器列表显示字符串
func getBrowserList(browsers []string) string {
	if len(browsers) == 0 {
		return "全部浏览器"
	}
	return fmt.Sprintf("%v", browsers)
}

