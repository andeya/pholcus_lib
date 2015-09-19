package spider_lib

// 基础包
import (
	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/context" //必需
	// "github.com/henrylee2cn/pholcus/logs"           //信息输出
	. "github.com/henrylee2cn/pholcus/app/spider" //必需
	// . "github.com/henrylee2cn/pholcus/app/spider/common" //选用

	// net包
	// "net/http" //设置http.Header
	"net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	// "regexp"
	"strconv"
	"strings"

	// 其他包
	"fmt"
	// "math"
	// "time"
)

func init() {
	Zaojia.AddMenu()
}

var cityNum int = 23

var zaojiaCity = map[int]string{1: "广州市",
	3:  "珠海市",
	4:  "江门市",
	5:  "佛山市",
	6:  "东莞市",
	7:  "惠州市",
	8:  "韶关市",
	9:  "中山市",
	10: "深圳市",
	12: "揭阳市",
	13: "汕尾市",
	14: "河源市",
	15: "梅州市",
	16: "阳江市",
	17: "茂名市",
	18: "汕头市",
	19: "湛江市",
	20: "潮州市",
	21: "清远市",
	22: "肇庆市",
	23: "云浮市",
}

var timeNum = []string{"2009-04",
	//"2008-01",
	//	"2008-02",
	//	"2008-03",
	//	"2008-04",
	//	"2008-05",
	"2008-06",
	//	"2008-07",
	//	"2008-08",
	//	"2008-09",
	//	"2008-10",
	//	"2008-11",
	//	"2008-12",
	//	"2009-01",
	//	"2009-02",
	"2009-03",
	"2009-05",
	"2009-06",
	"2009-07",
	"2009-08",
	"2009-09",
	"2009-10",
	"2009-11",
	"2009-12",
	"2010-01",
	"2010-02",
	"2010-03",
	"2010-04",
	"2010-05",
	"2010-06",
	"2010-07",
	"2010-08",
	"2010-09",
	"2010-10",
	"2010-11",
	"2010-12",
	"2011-01",
	"2011-02",
	"2011-03",
	"2011-04",
	"2011-05",
	"2011-06",
	"2011-07",
	"2011-08",
	"2011-09",
	"2011-10",
	"2011-11",
	"2011-12",
	"2012-01",
	"2012-02",
	"2012-03",
	"2012-04",
	"2012-05",
	"2012-06",
	"2012-07",
	"2012-08",
	"2012-09",
	"2012-10",
	"2012-11",
	"2012-12",
	"2013-01",
	"2013-02",
	"2013-03",
	"2013-04",
	"2013-05",
	"2013-06",
	"2013-07",
	"2013-08",
	"2013-09",
	"2013-10",
	"2013-11",
	"2013-12",
	"2014-01",
	"2014-02",
	"2014-03",
	"2014-04",
	"2014-05",
	"2014-06",
	"2014-07",
	"2014-08",
	"2014-09",
	"2014-10",
	"2014-11",
	"2014-12",
	"2015-01",
	"2015-02",
	"2015-03",
	"2015-04",
	"2015-05"}

var Zaojia = &Spider{
	Name:        "造价网",
	Description: "材价信息",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	// Keyword:   CAN_ADD,
	UseCookie: true,
	RuleTree: &RuleTree{
		Root: func(self *Spider) {
			self.AddQueue(map[string]interface{}{
				"Url":    "http://www.cjcost.com/login.php",
				"Rule":   "第1页",
				"Method": "POST",
				"PostData": url.Values{
					"account":  []string{"cenbaozong"},
					"password": []string{"123456"},
					"time":     []string{"1438566451528"},
				},
				"Temp": map[string]interface{}{"area": 3},
			})
		},

		Trunk: map[string]*Rule{
			"第1页": {
				ParseFunc: func(self *Spider, resp *context.Response) {
					for _, timeNumOne := range timeNum {
						self.AddQueue(
							map[string]interface{}{
								"Rule":     "生成请求",
								"Url":      "http://www.cjcost.com/official/iframe.php?cityid=3&areaid=0&time=" + timeNumOne + "&datatype=1&page=1",
								"Temp":     map[string]interface{}{"area": resp.GetTemp("area"), "date": timeNumOne},
								"Priority": resp.GetPriority(),
							})
					}
				},
			},

			"生成请求": {
				AidFunc: func(self *Spider, aid map[string]interface{}) interface{} {
					req := map[string]interface{}{
						"Rule":     aid["Rule"],
						"Temp":     map[string]interface{}{"area": aid["area"], "total": aid["loop"].([2]int)[1], "date": aid["data"]},
						"Priority": aid["area"].(int) * (-1),
					}
					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						req["Url"] = aid["Url"].(string) + strconv.Itoa(loop[0]+1)
						self.AddQueue(req)
					}
					return nil
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					query := resp.GetDom()
					total := query.Find("option").Last().Text()
					total2, _ := strconv.Atoi(total)
					url := strings.Split(resp.GetUrl(), "page=")[0] + "page="

					rule := fmt.Sprint(zaojiaCity[resp.GetTemp("area").(int)], "-", resp.GetTemp("date"))
					resp.SetTemp("total", total2)
					self.Trunk[rule] = self.Trunk["3-2008-01"]

					// 调用指定规则下辅助函数
					self.Aid("生成请求", map[string]interface{}{"loop": [2]int{1, total2}, "Url": url, "area": resp.GetTemp("area"), "data": resp.GetTemp("data"), "Rule": rule})

					// 用指定规则解析响应流
					self.Parse(rule, resp)
				},
			},
			"3-2008-01": {
				//注意：有无字段语义和是否输出数据必须保持一致
				OutFeild: []string{
					"序号",
					"编码",
					"名称",
					"型号规格",
					"单位",
					"品牌",
					"价格",
					"地区",
					"发布时期",
					"备注",
				},
				ParseFunc: func(self *Spider, resp *context.Response) {
					query := resp.GetDom()
					query.Find(".listtb .list").Each(func(i int, s *goquery.Selection) {
						//序列号
						number := s.Find("td").Eq(1).Text()
						//编码
						encode := s.Find("td").Eq(2).Text()
						//名称
						name := s.Find("td").Eq(3).Text()
						//型号、规格
						model := s.Find("td").Eq(4).Text()
						//单位
						unit := s.Find("td").Eq(5).Text()
						//品牌
						brand := s.Find("td").Eq(6).Text()
						//价格
						price := s.Find("td").Eq(7).Text()
						//地区
						area := s.Find("td").Eq(8).Text()
						//发布时间
						release := s.Find("td").Eq(9).Text()
						// 结果存入Response中转
						resp.AddItem(map[string]interface{}{
							self.OutFeild(resp, 0): strings.Trim(number, " \t\n"),
							self.OutFeild(resp, 1): strings.Trim(encode, " \t\n"),
							self.OutFeild(resp, 2): strings.Trim(name, " \t\n"),
							self.OutFeild(resp, 3): strings.Trim(model, " \t\n"),
							self.OutFeild(resp, 4): strings.Trim(unit, " \t\n"),
							self.OutFeild(resp, 5): strings.Trim(brand, " \t\n"),
							self.OutFeild(resp, 6): strings.Trim(price, " \t\n"),
							self.OutFeild(resp, 7): strings.Trim(area, " \t\n"),
							self.OutFeild(resp, 8): strings.Trim(release, " \t\n"),
						})
						//						page, _ := strconv.Atoi(strings.Split(resp.GetUrl(), "page=")[1])
						//						a := strconv.Itoa(resp.GetTemp("area").(int) + 1)
						//						if page == resp.GetTemp("total").(int) {
						//							resp.SetTemp("area", a)
						//							resp.SetPriority(resp.GetTemp("area").(int))
						//							self.CallRule("第1页", resp)
						//						}
					})
				},
			},
		},
	},
}
