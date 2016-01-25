// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package front

import (
	"github.com/caixw/typing/app"
	"github.com/caixw/typing/models"
	"github.com/issue9/logs"
)

type Tag struct {
	ID          int64
	Name        string
	Title       string
	Description string
	//Count       int // 关联的文章数量
}

func (t *Tag) Count() int {
	if cnt, found := stats.Tags[t.ID]; found {
		return cnt
	}

	r := &models.Relationship{TagID: t.ID}
	cnt, err := db.Count(r)
	if err != nil {
		logs.Error("front.Tag.Count:", err)
	}
	stats.Tags[t.ID] = cnt
	return cnt
}

func (t *Tag) Permalink() string {
	return app.TagURL(t.Name, 1)
}
