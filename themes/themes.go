// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package themes

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/caixw/typing/core"
	"github.com/issue9/logs"
)

var (
	cfg       *core.Config
	tpl       *template.Template // 当前使用的模板
	themesMap map[string]*Theme  // 所有的主题列表
)

// Theme 用于描述主题的相关信息，一般从主题目录下的theme.json获取。
type Theme struct {
	ID          string  `json:"id"`                    // 主题的唯一id，即其文件夹名称
	Name        string  `json:"name"`                  // 主题名称
	Description string  `json:"description,omitempty"` // 详细描述
	Version     string  `json:"version,omitempty"`     // 版本
	Screenshot  string  `json:"screenshot"`            // 截图地址，相对于当前目录下的public
	Author      *Author `json:"author"`                // 作者相关信息
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	URL   string `json:"url"`
}

// 从主题根目录加载所有的主题内容，并初始所有的主题下静态文件的路由。
// defaultTheme 为默认的主题。
func Init(config *core.Config, defaultTheme string) error {
	cfg = config

	fs, err := ioutil.ReadDir(cfg.ThemeDir)
	if err != nil {
		return err
	}
	if len(fs) == 0 {
		return errors.New("不存在任何主题目录")
	}
	themesMap = make(map[string]*Theme, len(fs))

	p := cfg.ThemeURLPrefix + "/"
	for _, file := range fs {
		if !file.IsDir() {
			continue
		}

		id := file.Name()
		themePath := cfg.ThemeDir + id + string(os.PathSeparator)

		theme, err := loadThemeFile(themePath + "theme.json")
		if err != nil {
			return err
		}
		theme.ID = id

		theme.Screenshot = p + id + "/" + theme.Screenshot
		themesMap[id] = theme
		cfg.Core.Static[p+id] = themePath + "public/"
	}

	return Switch(defaultTheme)
}

// 加theme.json文件
func loadThemeFile(path string) (*Theme, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	t := &Theme{}
	if err = json.Unmarshal(data, t); err != nil {
		return nil, err
	}

	return t, nil
}

// 返回所有的主题列表
func Themes() []*Theme {
	ret := make([]*Theme, 0, len(themesMap))
	for _, theme := range themesMap {
		ret = append(ret, theme)
	}
	return ret
}

// 切换主题
func Switch(id string) (err error) {
	tpl, err = template.ParseGlob(cfg.ThemeDir + id + "/*.html")
	return
}

// 输出指定模板
func Render(w http.ResponseWriter, name string, data interface{}) {
	err := tpl.ExecuteTemplate(w, name, data)
	if err != nil {
		logs.Error("themes.Render:", err)
		core.RenderJSON(w, http.StatusInternalServerError, nil, nil)
		return
	}
}
