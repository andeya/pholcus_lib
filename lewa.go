package spider_lib

// 基础包
import (
	// "github.com/PuerkitoBio/goquery" //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/context" //必需
	// "github.com/henrylee2cn/pholcus/logs"           //信息输出
	. "github.com/henrylee2cn/pholcus/app/spider"        //必需
	. "github.com/henrylee2cn/pholcus/app/spider/common" //选用

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
	Lewa.AddMenu()
}

var Lewa = &Spider{
	Name:        "乐蛙登录测试",
	Description: "乐蛙登录测试 [Auto Page] [http://accounts.lewaos.com]",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	// Keyword:   USE,
	UseCookie: true,
	RuleTree: &RuleTree{
		Root: func(self *Spider) {
			self.AddQueue(map[string]interface{}{"Url": "http://accounts.lewaos.com/", "Rule": "登录页"})
		},

		Trunk: map[string]*Rule{

			"登录页": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					// self.AddQueue(map[string]interface{}{
					// 	"Url":    "http://accounts.lewaos.com",
					// 	"Rule":   "登录后",
					// 	"Method": "POST",
					// 	"PostData": url.Values{
					// 		"username":  []string{""},
					// 		"password":  []string{""},
					// 		"login_btn": []string{"login_btn"},
					// 		"submit":    []string{"login_btn"},
					// 	},
					// })
					NewForm(
						self,
						"登录后",
						"http://accounts.lewaos.com",
						resp.GetDom().Find(".userlogin.lw-pl40"),
					).Inputs(map[string]string{
						"username": "",
						"password": "",
					}).Submit()
				},
			},
			"登录后": {
				//注意：有无字段语义和是否输出数据必须保持一致
				OutFeild: []string{
					"全部",
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					// 结果存入Response中转
					resp.AddItem(map[string]interface{}{
						self.OutFeild(resp, 0): resp.GetText(),
					})
				},
			},
		},
	},
}
