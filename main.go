// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 基于 Git 的博客系统。
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"runtime"

	"github.com/issue9/logs"
	"github.com/issue9/web"
	"github.com/issue9/web/encoding"
	"github.com/issue9/web/encoding/html"

	"github.com/caixw/gitype/app"
	"github.com/caixw/gitype/path"
	"github.com/caixw/gitype/vars"
)

const usage = `%s 是一个基于 Git 的博客系统。
源代码以 MIT 开源许可发布于：%s


常见用法：

%s -preview -appdir="./"
%s -appdir="./"


参数：

`

func main() {
	help := flag.Bool("h", false, "显示当前信息")
	version := flag.Bool("v", false, "显示程序的版本信息")
	preview := flag.Bool("preview", false, "是否启用预览模式")
	appdir := flag.String("appdir", "./", "指定运行的工作目录")
	init := flag.String("init", "", "初始化一个工作目录")
	flag.Usage = func() {
		fmt.Printf(usage, vars.Name, vars.URL, vars.Name, vars.Name)
		flag.PrintDefaults()
	}
	flag.Parse()

	switch {
	case *help:
		flag.Usage()
		return
	case *version:
		printVersion()
		return
	case len(*init) > 0:
		if err := app.Init(path.New(*init)); err != nil {
			panic(err)
		}
		fmt.Printf("操作成功，你现在可以在 %s 中修改具体的参数配置！\n", *init)
		return
	}

	path := path.New(*appdir)

	if *preview {
		fmt.Println("预览模式，监视以下数据文件：", path.DataDir)
	}

	if err := web.Init(path.ConfDir); err != nil {
		panic(err)
	}

	htmlMgr := html.New(nil)
	if err := encoding.AddMarshal("text/html", htmlMgr.Marshal); err != nil {
		panic(err)
	}

	// webhook 用到此编码
	if err := encoding.AddMarshal("application/json", json.Marshal); err != nil {
		panic(err)
	}

	logs.Critical(app.Run(path, htmlMgr, *preview))
	logs.Flush()
}

func printVersion() {
	fmt.Println(vars.Name, vars.Version(), "build with", runtime.Version())

	if len(vars.CommitHash()) > 0 {
		fmt.Println("Git commit hash:", vars.CommitHash())
	}
}
