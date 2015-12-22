// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package front

import (
	"encoding/json"
	"runtime"
	"strconv"

	"github.com/caixw/typing/app"
	"github.com/caixw/typing/models"
	"github.com/issue9/orm/fetch"
)

type Anchor struct {
	Link  string `json:"link"`
	Title string `json:"title"`
}

// 页面的基本信息
type Info struct {
	Title       string     // 网页的title值
	SiteURL     string     // 网站地址
	SiteName    string     // 网站名称
	SecondTitle string     // 副标题
	Canonical   string     // 当前页的唯一链接
	Keywords    string     // meta.keywords的值
	Description string     // meta.description的值
	AppVersion  string     // 当前程序的版本号
	GoVersion   string     // 编译的go版本号
	PostSize    int        // 总文章数量
	CommentSize int        // 总评论数量
	Uptime      int64      // 上线时间
	RSS         *Anchor    // RSS，NOTICE:指针方便模板判断其值是否为空
	Atom        *Anchor    // Atom
	PrevPage    *Anchor    // 前一页
	NextPage    *Anchor    // 下一页
	Tags        []*Tag     // 标签列表
	Tops        []*Comment // 最新评论的10条内容
	Hots        []*Post    // 评论最多的10条内容
	Menus       []Anchor   // 菜单
}

func getInfo() (*Info, error) {
	info := &Info{
		SiteURL:     opt.SiteURL,
		SiteName:    opt.SiteName,
		SecondTitle: opt.SecondTitle,
		Keywords:    opt.Keywords,
		Description: opt.Description,
		AppVersion:  app.Version,
		GoVersion:   runtime.Version(),
		Uptime:      opt.Uptime,
	}

	menus := make([]Anchor, 0, 10)
	if err := json.Unmarshal([]byte(opt.Menus), &menus); err != nil {
		return nil, err
	}
	info.Menus = menus

	var err error
	sql := "SELECT COUNT(*) as cnt FROM #posts WHERE {state}=?" // TODO 预编译成stmt
	if info.PostSize, err = getSize(sql, models.PostStatePublished); err != nil {
		return nil, err
	}

	sql = "SELECT COUNT(*) as cnt FROM #comments WHERE {state}=?"
	if info.CommentSize, err = getSize(sql, models.CommentStateApproved); err != nil {
		return nil, err
	}

	if info.Tags, err = getTags(); err != nil {
		return nil, err
	}

	if info.Tops, err = getTops(); err != nil {
		return nil, err
	}

	return info, nil
}

func getSize(sql string, args ...interface{}) (int, error) {
	rows, err := db.Query(true, sql, args...)
	if err != nil {
		return 0, err
	}
	cnts, err := fetch.ColumnString(true, "cnt", rows)
	rows.Close()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(cnts[0])
}

func getTops() ([]*Comment, error) {
	sql := `SELECT c.{content} AS {Content}, c.{id} AS {ID}, p.{id} AS {PostID}, p.{name} AS {PostName}
	FROM #comments AS c
	LEFT JOIN #posts AS p ON c.{postID}=p.{id}
	WHERE c.{state}=?
	ORDER BY c.{id} DESC
	LIMIT ?`
	rows, err := db.Query(true, sql, models.CommentStateApproved, opt.SidebarSize)
	if err != nil {
		return nil, err
	}

	tops := make([]*Comment, 0, opt.SidebarSize)
	_, err = fetch.Obj(&tops, rows)
	rows.Close()
	if err != nil {
		return nil, err
	}
	return tops, nil
}

func getTags() ([]*Tag, error) {
	sql := "SELECT  {title} AS {Title}, {name} AS {Name} FROM #tags"
	rows, err := db.Query(true, sql)
	if err != nil {
		return nil, err
	}

	tags := make([]*Tag, 0, opt.SidebarSize)
	_, err = fetch.Obj(&tags, rows)
	rows.Close()
	if err != nil {
		return nil, err
	}

	return tags, nil
}