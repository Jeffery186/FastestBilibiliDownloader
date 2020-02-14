package parser

import (
	"fmt"
	"github.com/tidwall/gjson"
	"simple-golang-crawler/engine"
	"simple-golang-crawler/fetcher"
	"simple-golang-crawler/model"
)

var getAidUrlTemp = "https://api.bilibili.com/x/space/arc/search?mid=%d&ps=30&tid=0&pn=%d&keyword=&order=pubdate&jsonp=jsonp"
var getCidUrlTemp = "https://api.bilibili.com/x/player/pagelist?aid=%d"

func UpSpaceParseFun(contents []byte, url string) engine.ParseResult {
	var retParseResult engine.ParseResult
	value := gjson.GetManyBytes(contents, "data.list.vlist", "data.page")

	var upid int64
	retParseResult.Requests, retParseResult.Items, upid = getAidDetailReqList(value[0])
	retParseResult.Requests = append(retParseResult.Requests, getNewBilibiliUpSpaceReqList(value[1], upid)...)

	return retParseResult

}

func getAidDetailReqList(pageInfo gjson.Result) ([]*engine.Request, []*engine.Item, int64) {
	var retRequests []*engine.Request
	var retItems []*engine.Item
	var upid int64
	for _, i := range pageInfo.Array() {
		aid := i.Get("aid").Int()
		upid = i.Get("mid").Int()
		title := i.Get("title").String()
		reqUrl := fmt.Sprintf(getCidUrlTemp, aid)
		videoAid := model.NewVideoAidInfo(aid, title)
		reqParseFunction := GenGetAidChildrenParseFun(videoAid)
		req := engine.NewRequest(reqUrl, reqParseFunction, fetcher.DefaultFetcher)
		retRequests = append(retRequests, req)

		item := engine.NewItem(videoAid)
		retItems = append(retItems, item)
	}
	return retRequests, retItems, upid
}

func getNewBilibiliUpSpaceReqList(pageInfo gjson.Result, upid int64) []*engine.Request {
	var retRequests []*engine.Request

	count := pageInfo.Get("count").Int()
	pn := pageInfo.Get("pn").Int()
	ps := pageInfo.Get("ps").Int()
	var extraPage int64
	if count%ps > 0 {
		extraPage = 1
	}
	totalPage := count/ps + extraPage
	for i := int64(1); i < (totalPage - totalPage + 1); i++ {
		if i == pn {
			continue
		}
		reqUrl := fmt.Sprintf(getAidUrlTemp, upid, i)
		req := engine.NewRequest(reqUrl, UpSpaceParseFun, fetcher.DefaultFetcher)
		retRequests = append(retRequests, req)
	}
	return retRequests
}

func GetRequestByUpId(upid int64) *engine.Request {
	reqUrl := fmt.Sprintf(getAidUrlTemp, upid, 1)
	return engine.NewRequest(reqUrl, UpSpaceParseFun, fetcher.DefaultFetcher)
}
