/**
 *   Copyright 2018 Guanlaolin
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package main

import (
	"encoding/json"
	"log"
)

// urls
const url = `https://msdn.itellyou.cn/`
const cataUrl = `https://msdn.itellyou.cn/Category/Index`
const langUrl = `https://msdn.itellyou.cn/Category/GetLang`
const listUrl = `https://msdn.itellyou.cn/Category/GetList`
const productUrl = `https://msdn.itellyou.cn/Category/GetProduct`

// regex
const cataRegex = `<a\shref="javascript:void\(0\);".*?data-menuid="(.*?)".*?>(.*?)</a>`

// must be header
var header map[string][]string = map[string][]string{
	"Origin":       {"https://msdn.itellyou.cn"},
	"Referer":      {"https://msdn.itellyou.cn/"},
	"Content-Type": {"application/x-www-form-urlencoded; charset=UTF-8"},
}

// 类别
type Cata struct {
	Name  string
	Id    string
	Items []Item
}

// 类别下的具体项，如Win10、Win7...
type Item struct {
	Name      string
	Id        string
	ItemLangs []ItemLanguage
}

// 每个具体项有哪些语言类别可供选择
type ItemLanguage struct {
	Lang     string
	Id       string
	ItemSums []ItemSummary
}

type ItemLanguageResp struct {
	Status bool
	Result []ItemLanguage
}

type ItemSummary struct {
	Name   string
	Id     string
	Post   string
	URL    string
	Detail ItemDetail
}

type ItemSummaryResp struct {
	Status bool
	Result []ItemSummary
}

type ItemDetail struct {
	FileName string
	SHA1     string
	Size     string
	PubTime  string `json:"PostDateString"`
	URL      string `json:"DownLoad"`
}

type ItemDetailResp struct {
	Status bool
	Result ItemDetail
}

func getCatasUrl(url string) ([]Cata, error) {
	var catas []Cata

	buf, err := fetchBody("GET", url, nil, "")
	if err != nil {
		log.Fatal("fetchBody:", err)
		return nil, err
	}
	LOG_DEBUGF("%s response:%s\n", url, buf)

	temp := anylazeAll(buf, cataRegex)
	LOG_DEBUGF("%s\n", temp)

	for _, v := range temp {
		catas = append(catas, Cata{Name: string(v[2]), Id: string(v[1])})
	}
	LOG_DEBUGF("%s\n", catas)

	return catas, nil
}

func getCataItemsUrl(url string, body string) ([]Item, error) {
	var items []Item

	buf, err := fetchBody("POST", url, header, "id="+body)
	if err != nil {
		LOG_DEBUG("fetchBody:", err)
		return nil, err
	}
	LOG_DEBUGF("%s\n", buf)

	if json.Unmarshal(buf, &items) != nil {
		log.Print("Unmarshal:", err)
		return nil, err
	}
	LOG_DEBUG("%s\n", items)

	return items, nil
}

func getCataItemLangsUrl(url string, body string) ([]ItemLanguage, error) {
	var temp ItemLanguageResp

	buf, err := fetchBody("POST", url, header, "id="+body)
	if err != nil {
		log.Print("fetchBody:", err)
		return nil, err
	}
	LOG_DEBUGF("%s\n", buf)

	err = json.Unmarshal(buf, &temp)
	if err != nil {
		log.Print("Unmarshal:", err)
		return nil, err
	}
	//log.Println(temp)

	return temp.Result, nil
}

func getCataItemLangListsUrl(url string, body string) ([]ItemSummary, error) {
	var temp ItemSummaryResp

	buf, err := fetchBody("POST", url, header, "id="+body)
	if err != nil {
		log.Print("fetchBody:", err)
		return nil, err
	}
	LOG_DEBUGF("%s\n", buf)

	if json.Unmarshal(buf, &temp) != nil {
		log.Print("Unmarshal:", err)
		return nil, err
	}
	//log.Printf("%s\n", temp)

	return temp.Result, nil
}

func getCataItemLangListDetail(url string, body string) (ItemDetail, error) {
	var temp ItemDetailResp

	buf, err := fetchBody("POST", url, header, "id="+body)
	if err != nil {
		log.Print("fetchBody:", err)
		return ItemDetail{}, err
	}
	LOG_DEBUGF("%s\n", buf)

	err = json.Unmarshal(buf, &temp)
	if err != nil {
		log.Print("Unmarshal:", err)
		return ItemDetail{}, err
	}
	//log.Printf("%s\n", temp)

	return temp.Result, nil
}

func joinBytes(o ...[]byte) []byte {
	var temp []byte

	for _, v1 := range o {
		for _, v2 := range v1 {
			temp = append(temp, v2)
		}
	}

	return temp
}
