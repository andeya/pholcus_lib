package spider_lib

// 基础包
import (
	// "github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/context" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	. "github.com/henrylee2cn/pholcus/app/spider/common"    //选用
	//	"github.com/henrylee2cn/pholcus/logs" //信息输出

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	//	"encoding/json"
	//	"encoding/xml"

	// 字符串处理包
	//	"regexp"
	//	"strconv"
	// "strings"

	// 其他包
	"fmt"
	//	"log"
	// "math"
	// "github.com/henrylee2cn/pholcus/common/mgo" //信息输出
	// "time"
)

func init() {
	SinaWeibo.Register()
}

var weiboList = []string{
	// "http://weibo.com/2074674805",
	// "http://weibo.com/1831775890",
	// "http://weibo.com/2083716120",
	// "http://weibo.com/1863128044",
	// "http://weibo.com/1948532334",
	// "http://weibo.com/2083385932",
	// "http://weibo.com/2085991971",
	// "http://weibo.com/2091639913",
	// "http://weibo.com/2074139685",
	// "http://weibo.com/5052253202",
	// "http://weibo.com/2089145113",
	// "http://weibo.com/2074673163",
	// "http://weibo.com/2096456297",
	// "http://weibo.com/2095476413",
	// "http://weibo.com/2591374274",
	// "http://weibo.com/1949136174",
	// "http://weibo.com/2092092551",
	// "http://weibo.com/2092984355",
	// "http://weibo.com/2073349753",
	// "http://weibo.com/2065240833",
	// "http://weibo.com/1949250614",
	// "http://weibo.com/2088327553",
	// "http://weibo.com/2089110877",
	// "http://weibo.com/2087242923",
	// "http://weibo.com/2085927973",
	// "http://weibo.com/2090404631",
	// "http://weibo.com/2088000745",
	// "http://weibo.com/2093862715",
	// "http://weibo.com/2082602061",
	// "http://weibo.com/2088311983",
	// "http://weibo.com/2086544355",
	// "http://weibo.com/2091811107",
	// "http://weibo.com/1947391624",
	// "http://weibo.com/1642351200",
	// "http://weibo.com/2074226615",
	// "http://weibo.com/1890095932",
	// "http://weibo.com/2392105140",
	// "http://weibo.com/2101569757",
	// "http://weibo.com/3089181212",
	// "http://weibo.com/1883314957",
	// "http://weibo.com/1769972891",
	// "http://weibo.com/1688018173",
	// "http://weibo.com/2824391155",
	// "http://weibo.com/2824364481",
	// "http://weibo.com/1852959903",
	// "http://weibo.com/2087063403",
	// "http://weibo.com/2087379615",
	// "http://weibo.com/1864353197",
	// "http://weibo.com/1662319474",
	"http://weibo.com/2308610287",
}

var timer_SinaWeibo = DailyFixedTimer{
	"请求列表": [3]int{18, 46, 00},
}

var SinaWeibo = &Spider{
	Name:        "IT新浪微博",
	Description: "IT新浪微博 [weibo.com]",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	// Keyword:     USE,
	EnableCookie: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {

			// ctx.Aid("入口", map[string]interface{}{})
			for _, v := range weiboList {
				ctx.AddQueue(&context.Request{
					Url:          v,
					Rule:         "请求列表",
					Duplicatable: true,
					DownloaderID: 1,
				})
			}
		},

		Trunk: map[string]*Rule{

			"请求列表": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					// defer func() {
					// 	// 循环请求
					// 	time.Sleep(60 * time.Second)
					// 	timer_SinaWeibo.Wait("请求列表")
					// 	ctx.Aid("请求列表", map[string]interface{}{})
					// }()

					return nil
				},
				ItemFields: []string{
					"微博名",
					"粉丝数",
					"微博数",
					"全文",
				},
				ParseFunc: func(ctx *Context) {
					defer func() {
						fmt.Println(recover())
					}()

					query := ctx.GetDom()
					fmt.Println(query.Find("html").Text())
					微博名 := query.Find(".username").Text()
					粉丝数 := query.Find(".W_f16").Eq(1).Text()
					微博数 := query.Find(".W_f16").Eq(2).Text()
					fmt.Println("微博名", 微博名, 粉丝数, 微博数)
					// _, err := mgo.Mgo("insert", map[string]interface{}{
					// 	"Database":   "1_3",
					// 	"Collection": "IT新浪微博",
					// 	"Docs": []map[string]interface{}{
					// 		{
					// 			"微博名":  微博名,
					// 			"粉丝数":  粉丝数,
					// 			"微博数":  微博数,
					// 			"time": time.Now().Format("2006-01-02 15:04:05"),
					// 		},
					// 	},
					// })
					// if err != nil {
					// 	fmt.Println(err)
					// }
					ctx.Output(map[int]interface{}{
						0: 微博名,
						1: 粉丝数,
						2: 微博数,
						3: query.Text(),
					})
				},
			},
		},
	},
}
