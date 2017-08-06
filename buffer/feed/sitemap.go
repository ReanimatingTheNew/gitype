// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package feed 包提供了根据数据生成 sitemap，atom，rss 的功能。
package feed

import (
	"bytes"
	"strconv"
	"time"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/vars"
)

const (
	sitemapHeader = `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`

	sitemapFooter = `</urlset>`
)

// BuildSitemap 生成一个符合 sitemap 规范的 XML 文本 buffer。
func BuildSitemap(d *data.Data) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	w := &writer{
		buf: buf,
	}

	w.writeString(xmlHeader)

	if len(d.Config.Sitemap.XslURL) > 0 {
		w.writePI("xml-stylesheet", map[string]string{
			"type": "text/xsl",
			"href": d.Config.Sitemap.XslURL,
		})
	}

	w.writeString(sitemapHeader)

	addPostsToSitemap(w, d)

	if d.Config.Sitemap.EnableTag {
		addTagsToSitemap(w, d)
	}

	w.writeString(sitemapFooter)

	if w.err != nil {
		return nil, w.err
	}
	return buf, nil
}

func addPostsToSitemap(w *writer, d *data.Data) {
	sitemap := d.Config.Sitemap
	for _, p := range d.Posts {
		loc := d.Config.URL + p.Permalink
		addItemToSitemap(w, loc, sitemap.PostChangefreq, p.Modified, sitemap.PostPriority)
	}
}

func addTagsToSitemap(w *writer, d *data.Data) error {
	now := time.Now().Unix()
	sitemap := d.Config.Sitemap

	loc := d.Config.URL + vars.TagsURL()
	addItemToSitemap(w, loc, sitemap.TagChangefreq, now, sitemap.TagPriority)

	for _, tag := range d.Tags {
		loc = d.Config.URL + tag.Permalink
		addItemToSitemap(w, loc, sitemap.TagChangefreq, now, sitemap.TagPriority)
	}
	return nil
}

func addItemToSitemap(w *writer, loc, changefreq string, lastmod int64, priority float64) {
	w.writeString("<url>\n")

	w.writeElement("loc", loc, nil)
	t := time.Unix(lastmod, 0)
	w.writeElement("lastmod", t.Format("2006-01-02T15:04:05-07:00"), nil)
	w.writeElement("changefreq", changefreq, nil)
	w.writeElement("priority", strconv.FormatFloat(priority, 'f', 1, 32), nil)

	w.writeString("</url>\n")
}
