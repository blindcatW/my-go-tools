package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Request struct {
	Url        string
	ParserFunc func([]byte) ParseResult
}

type ParseResult struct {
	Requests []Request
	Items    []interface{}
}

// 起始URL
const URL = "http://www.zhenai.com/zhenghun"

func main() {
	run(Request{
		Url:        URL,
		ParserFunc: parseCityList,
	})
}

func run(seeds ...Request) {
	var requests []Request
	for _, seed := range seeds {
		requests = append(requests, seed)
	}

	for len(requests) > 0 {
		req := requests[0]
		requests = requests[1:]
		log.Printf("fetching %s", req.Url)
		res, err := fetch(req.Url)
		if err != nil {
			log.Printf("Fetch url %s error %v \n", req.Url, err)
			continue
		}
		parseResult := req.ParserFunc(res)
		requests = append(requests, parseResult.Requests...)
		for _, item := range parseResult.Items {
			log.Printf("item: %v", item)
		}
	}
}

func fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("StatusCode: " + string(resp.StatusCode))
	}

	info, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// 使用正则解析需要的城市列表信息
const cityReg = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`

func parseCityList(contents []byte) (res ParseResult) {
	reg := regexp.MustCompile(cityReg)
	// matches := reg.FindAll(contents, -1)
	matches := reg.FindAllSubmatch(contents, -1)
	result := ParseResult{}
	for _, m := range matches {
		// fmt.Printf("City: %s, URL: %s\n", m[2], m[1])
		result.Items = append(result.Items, string(m[2]))
		result.Requests = append(result.Requests, Request{
			Url:        string(m[1]),
			ParserFunc: parseInfo,
			// ParserFunc: parserProfile,
		})
	}
	return result
}

// 把页面上的名字、性别、照片找出来
// <a href="http://album.zhenai.com/u/1283740719" target="_blank">纯善至美</a>
// <td width="180"><span class="grayL">性别：</span>男士</td>
// <div class="photo"><a href="http://album.zhenai.com/u/1283740719" target="_blank"><img src="https://photo.zastatic.com/images/photo/320936/1283740719/201357047986646980.jpg?scrop=1&amp;crop=1&amp;w=140&amp;h=140&amp;cpos=north" alt="纯善至美"></a></div>
const infoReg = `<div class="photo"><a href="(http://album.zhenai.com/u/[0-9]+)" target="_blank"><img src="(https://photo.zastatic.com/images/photo/[0-9/]+.jpg)?[^"]*" alt="([^"]*)"></a></div>`

type profile struct {
	Name string
	Img  string
}

func parseInfo(contents []byte) (res ParseResult) {
	reg := regexp.MustCompile(infoReg)
	matches := reg.FindAllSubmatch(contents, -1)
	result := ParseResult{}
	for _, m := range matches {
		fmt.Printf("info items : %s\n", m)
		result.Items = append(result.Items, string(m[2]))
		result.Requests = append(result.Requests, Request{
			Url:        string(m[1]),
			ParserFunc: nilParser,
		})
	}
	return result
}

// 使用正则解析需要的个人信息
// 不登录该问不了，这部分先不弄了
const ageReg = `<td width="180"><span class="grayL">年龄：</span>([\d])+</td>`

func parserProfile(contets []byte) (res ParseResult) {
	//reg := regexp.MustCompile(ageReg)
	//reg.FindSubmatch(contets)
	// fmt.Println(string(contets))
	return ParseResult{}
}
func nilParser([]byte) ParseResult {
	return ParseResult{}
}

func getProfile(contents []byte) {

}
