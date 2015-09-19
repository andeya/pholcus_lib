package spider_lib

// 基础包
import (
	// "github.com/PuerkitoBio/goquery"                          //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/context" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	// . "github.com/henrylee2cn/pholcus/app/spider/common" //选用
	// "github.com/henrylee2cn/pholcus/logs"

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	//"encoding/json"

	// 字符串处理包
	//"regexp"
	// "strconv"
	//	"strings"

	// 其他包
	// "fmt"
	// "math"
	// "time"
)

func init() {
	FileTest.AddMenu()
}

var FileTest = &Spider{
	Name:        "文件下载测试",
	Description: "文件下载测试",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	// Keyword:   USE,
	UseCookie: false,
	RuleTree: &RuleTree{
		Root: func(self *Spider) {
			self.AddQueue(map[string]interface{}{"Url": "https://www.baidu.com/img/bd_logo1.png", "Rule": "百度图片"})
			self.AddQueue(map[string]interface{}{"Url": "https://github.com/henrylee2cn/pholcus", "Rule": "Pholcus页面"})
		},

		Trunk: map[string]*Rule{

			"百度图片": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					resp.AddFile("baidu")
				},
			},
			"Pholcus页面": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					resp.AddFile()
				},
			},
		},
	},
}
