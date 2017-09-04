// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package data

import (
	"fmt"
	"strconv"
	"time"

	"github.com/caixw/typing/vars"
)

// 文章是否过时的比较方式
const (
	OutdatedTypeCreated  = "created"  // 以创建时间作为对比
	OutdatedTypeModified = "modified" // 以修改时间作为对比
)

// Outdated 描述过时文章的提示信息。
//
// 理论上把有关 Outdated 的信息，直接在模板中对文章的创建时间戳进行比较，
// 是比通过配置来比较会更加方便，也不会更任何的后期工作量。之所以把这个功能放在后端，
// 而不是模板层面，是因为觉得模板应该只负责展示页面，而不是用于处理逻辑内容。
type Outdated struct {
	Type     string        `yaml:"type"`     // 比较的类型，创建时间或是修改时间
	Duration time.Duration `yaml:"duration"` // 超时的时间，可以使用 time.Duration 的字符串值
	Content  string        `yaml:"content"`  // 提示的内容，普通文字，不能为 html
}

// Author 描述作者信息
type Author struct {
	Name   string `yaml:"name"`
	URL    string `yaml:"url,omitempty"`
	Email  string `yaml:"email,omitempty"`
	Avatar string `yaml:"avatar,omitempty"`
}

// Link 描述链接的内容
type Link struct {
	// 链接对应的图标名称，fontawesome 图标名称，不用带 fa- 前缀。
	// 也有可能是链接，模板根据情况自动选择。
	Icon  string `yaml:"icon,omitempty"`
	Title string `yaml:"title,omitempty"` // 链接的 title 属性
	Rel   string `yaml:"rel,omitempty"`   // 链接的 rel 属性
	URL   string `yaml:"url"`             // 链接地址
	Text  string `yaml:"text"`            // 链接的文本
}

// Icon 表示程序图标
type Icon struct {
	URL   string `yaml:"url"`
	Type  string `yaml:"type"` // mime type
	Sizes string `yaml:"sizes"`
}

// FieldError 表示加载文件出错时，具体的错误信息
type FieldError struct {
	File    string // 所在文件
	Message string // 错误信息
	Field   string // 所在的字段
}

func loadLinks(path *vars.Path) ([]*Link, error) {
	links := make([]*Link, 0, 20)
	if err := loadYamlFile(path.MetaLinksFile, &links); err != nil {
		return nil, err
	}

	for index, link := range links {
		if err := link.sanitize(); err != nil {
			err.File = path.MetaLinksFile
			err.Field = "[" + strconv.Itoa(index) + "]." + err.Field
			return nil, err
		}
	}

	return links, nil
}

func (err *FieldError) Error() string {
	return fmt.Sprintf("在文件 %s 中的 %s 字段发生错误：%s", err.File, err.Field, err.Message)
}

func (icon *Icon) sanitize() *FieldError {
	if len(icon.URL) == 0 {
		return &FieldError{Field: "url", Message: "不能为空"}
	}

	return nil
}

func (link *Link) sanitize() *FieldError {
	if len(link.Text) == 0 {
		return &FieldError{Field: "text", Message: "不能为空"}
	}

	if len(link.URL) == 0 {
		return &FieldError{Field: "url", Message: "不能为空"}
	}

	return nil
}

func (author *Author) sanitize() *FieldError {
	if len(author.Name) == 0 {
		return &FieldError{Field: "name", Message: "不能为空"}
	}

	return nil
}

func (o *Outdated) sanitize() *FieldError {
	if o.Type != OutdatedTypeCreated && o.Type != OutdatedTypeModified {
		return &FieldError{Message: "无效的值", Field: "outdated.type"}
	}

	if len(o.Content) == 0 {
		return &FieldError{Message: "不能为空", Field: "outdated.content"}
	}

	return nil
}
