package spider_lib

// 基础包
import (
	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/context" //必需
	// "github.com/henrylee2cn/pholcus/logs"              //信息输出
	. "github.com/henrylee2cn/pholcus/app/spider" //必需
	// . "github.com/henrylee2cn/pholcus/app/spider/common" //选用

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	// "regexp"
	// "strconv"
	// "strings"

	// 其他包
	// "fmt"
	// "math"
	// "time"
)

func init() {
	Kaola.AddMenu()
}

// 考拉海淘,海外直采,7天无理由退货,售后无忧!考拉网放心的海淘网站!
var Kaola = &Spider{
	Name:        "考拉海淘",
	Description: "考拉海淘商品数据 [Auto Page] [www.kaola.com]",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	// Keyword:   USE,
	UseCookie: false,
	RuleTree: &RuleTree{
		Root: func(self *Spider) {
			self.AddQueue(map[string]interface{}{"Url": "http://www.kaola.com", "Rule": "获取版块URL"})
		},

		Trunk: map[string]*Rule{

			"获取版块URL": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					query := resp.GetDom()
					lis := query.Find("#funcTab li a")
					lis.Each(func(i int, s *goquery.Selection) {
						if i == 0 {
							return
						}
						if url, ok := s.Attr("href"); ok {
							self.AddQueue(map[string]interface{}{"Url": url, "Rule": "商品列表", "Temp": map[string]interface{}{"goodsType": s.Text()}})
						}
					})
				},
			},

			"商品列表": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					query := resp.GetDom()
					query.Find(".proinfo").Each(func(i int, s *goquery.Selection) {
						if url, ok := s.Find("a").Attr("href"); ok {
							self.AddQueue(map[string]interface{}{
								"Url":  "http://www.kaola.com" + url,
								"Rule": "商品详情",
								"Temp": map[string]interface{}{"goodsType": resp.GetTemp("goodsType").(string)},
							})
						}
					})
				},
			},

			"商品详情": {
				//注意：有无字段语义和是否输出数据必须保持一致
				OutFeild: []string{
					"标题",
					"价格",
					"品牌",
					"采购地",
					"评论数",
					"类别",
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					query := resp.GetDom()
					// 获取标题
					title := query.Find(".product-title").Text()

					// 获取价格
					price := query.Find("#js_currentPrice span").Text()

					// 获取品牌
					brand := query.Find(".goods_parameter li").Eq(0).Text()

					// 获取采购地
					from := query.Find(".goods_parameter li").Eq(1).Text()

					// 获取评论数
					discussNum := query.Find("#commentCounts").Text()

					// 结果存入Response中转
					resp.AddItem(map[string]interface{}{
						self.OutFeild(resp, 0): title,
						self.OutFeild(resp, 1): price,
						self.OutFeild(resp, 2): brand,
						self.OutFeild(resp, 3): from,
						self.OutFeild(resp, 4): discussNum,
						self.OutFeild(resp, 5): resp.GetTemp("goodsType"),
					})
				},
			},
		},
	},
}
