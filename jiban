package pholcus_lib

import (
	"strconv"
	"strings"

	"github.com/henrylee2cn/pholcus/app/downloader/request"
	. "github.com/henrylee2cn/pholcus/app/spider" //必需
	"github.com/henrylee2cn/pholcus/common/goquery"
)

func init() {
	Jiban.Register()
}

var Jiban = &Spider{
	Name:         "羁绊",
	Description:  "羁绊动漫二次元资讯 [http://www.005.tv/zx/]",
	EnableCookie: true,

	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{
				Url:         "http://www.005.tv/zx/list_526_1.html",
				Rule:        "请求",
				Temp:        map[string]interface{}{"p": 1},
				ConnTimeout: -1,
				Reloadable:  true,
			})

		},
		Trunk: map[string]*Rule{
			"请求": {
				ParseFunc: func(ctx *Context) {
					var curr = ctx.GetTemp("p", int(0)).(int)
					ctx.GetDom().Find(".pages .dede_pages  .pagelist  .thisclass a").Each(func(ii int, iio *goquery.Selection) {
						url2, _ := iio.Attr("href")
						if url2 != "javascript:void(0);" {
							if curr > 100 {
								return
							}
						}
					})
					ctx.AddQueue(&request.Request{
						Url:         "http://www.005.tv/zx/list_526_" + strconv.Itoa(curr+1) + ".html",
						Rule:        "请求",
						Temp:        map[string]interface{}{"p": curr + 1},
						ConnTimeout: -1,
						Reloadable:  true,
					})
					ctx.Parse("获取列表")
				},
			},

			"获取列表": {
				ParseFunc: func(ctx *Context) {
					ctx.GetDom().
						Find(".article-list ul li .xs-100 div h3 a").
						Each(func(i int, s *goquery.Selection) {
							url, _ := s.Attr("href")
							ctx.AddQueue(&request.Request{
								Url:         url,
								Rule:        "news",
								ConnTimeout: -1,
							})
						})
				},
			},

			"news": {
				ItemFields: []string{
					"title",
					"content",
					"img_url",
					"time",
					"vedio_url",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					var title, time, img_url, content, vedio_url string
					query.Find(".article-list-wrap").
						Each(func(j int, jo *goquery.Selection) {
							title = jo.Find(".articleTitle-name").Text()
							time = jo.Find("span.time").Text()
							jo.Find(".articleContent img").Each(func(x int, xo *goquery.Selection) {
								if img, ok := xo.Attr("src"); ok {
									img_url = img_url + img + ","
								}
							})
							jo.Find(".articleContent vedio").Each(func(b int, bo *goquery.Selection) {
								if vedio, ok := bo.Attr("src"); ok {
									vedio_url = vedio_url + vedio + ","
								}
							})
							jo.Find(".articleContent img").ReplaceWithHtml("#image#")
							jo.Find(".articleContent img").Remove()
							content, _ = jo.Find(".articleContent").Html()
							content = strings.Replace(content, `"`, `'`, -1) + "<p>来源：ACGdoge</p>"

						})
					ctx.Output(map[int]interface{}{
						0: title,
						1: content,
						2: img_url,
						3: time,
						4: vedio_url,
					})
				},
			},
		},
	},
}
