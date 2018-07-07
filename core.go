// core.go
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Cata struct {
	Name  string // cata name
	Id    string
	Items []Item
}

type Item struct {
	Name      string
	Id        string
	ItemLangs []ItemLanguage
}

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

func fetch(method string, url string, header map[string][]string,
	body string) (*http.Response, error) {

	client := http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		log.Print("NewRequest:", err)
		return nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v[0])
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Print("Do:", err)
		return nil, err
	}

	return resp, nil
}

func anylazeAll(buf []byte, rule string) [][][]byte {
	reg := regexp.MustCompile(rule)
	res := reg.FindAllSubmatch(buf, -1)
	//log.Printf("%s\n", res)

	return res
}

/**
 * Read HTTP response body
 *
 * @param resp *http.Response
 *   http response
 * @param _close bool
 *   true : close resp.Body
 *   false : do not close resp.Body
 *
 * @return
 */
func readRespBody(resp *http.Response, _close bool) ([]byte, error) {
	if _close {
		defer resp.Body.Close()
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("ReadAll:", err)
		return nil, err
	}

	return buf, nil
}

func getCatasUrl(url string) ([]Cata, error) {
	var catas []Cata

	resp, err := fetch("GET", url, nil, "")
	if err != nil {
		log.Fatal("fetch:", err)
		return nil, err
	}

	buf, err := readRespBody(resp, true)
	if err != nil {
		log.Fatal("readRespBody:", err)
		return nil, err
	}

	rule := `<a\shref="javascript:void\(0\);".*?data-menuid="(.*?)".*?>(.*?)</a>`

	temp := anylazeAll(buf, rule)

	for _, v := range temp {
		//log.Printf("%s", v)
		catas = append(catas, Cata{Name: string(v[2]), Id: string(v[1])})
	}
	//log.Printf("%s\n", catas)

	return catas, nil
}

func getCataItemsUrl(url string, body string) ([]Item, error) {
	var items []Item

	header := make(map[string][]string)
	header["Origin"] = append(header["Origin"], "https://msdn.itellyou.cn")
	header["Referer"] = append(header["Referer"], "https://msdn.itellyou.cn/")
	header["Content-Type"] = append(header["Content-Type"],
		"application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := fetch("POST", url, header, "id="+body)
	if err != nil {
		log.Print("fetch:", err)
		return nil, err
	}
	//log.Println("status:", resp.Status)

	buf, err := readRespBody(resp, true)
	if err != nil {
		log.Print("readRespBody:", err)
		return nil, err
	}
	//log.Printf("%s\n", buf)

	if json.Unmarshal(buf, &items) != nil {
		log.Print("Unmarshal:", err)
		return nil, err
	}
	//log.Printf("%s\n", items)

	return items, nil
}

func getCataItemLangsUrl(url string, body string) ([]ItemLanguage, error) {
	var temp ItemLanguageResp

	header := make(map[string][]string)
	header["Origin"] = append(header["Origin"], "https://msdn.itellyou.cn")
	header["Referer"] = append(header["Referer"], "https://msdn.itellyou.cn/")
	header["Content-Type"] = append(header["Content-Type"],
		"application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := fetch("POST", url, header, "id="+body)
	if err != nil {
		log.Print("fetch:", err)
		return nil, err
	}
	//log.Println("status:", resp.Status)

	buf, err := readRespBody(resp, true)
	if err != nil {
		log.Print("readRespBody:", err)
		return nil, err
	}
	//log.Printf("%s\n", buf)

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

	header := make(map[string][]string)
	header["Origin"] = append(header["Origin"], "https://msdn.itellyou.cn")
	header["Referer"] = append(header["Referer"], "https://msdn.itellyou.cn/")
	header["Content-Type"] = append(header["Content-Type"],
		"application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := fetch("POST", url, header, "id="+body)
	if err != nil {
		log.Print("fetch:", err)
		return nil, err
	}
	//log.Println("status:", resp.Status)

	buf, err := readRespBody(resp, true)
	if err != nil {
		log.Print("readRespBody:", err)
		return nil, err
	}
	//log.Printf("%s\n", buf)

	if json.Unmarshal(buf, &temp) != nil {
		log.Print("Unmarshal:", err)
		return nil, err
	}
	//log.Printf("%s\n", temp)

	return temp.Result, nil
}

func getCataItemLangListDetail(url string, body string) (ItemDetail, error) {
	var temp ItemDetailResp

	header := make(map[string][]string)
	header["Origin"] = append(header["Origin"], "https://msdn.itellyou.cn")
	header["Referer"] = append(header["Referer"], "https://msdn.itellyou.cn/")
	header["Content-Type"] = append(header["Content-Type"],
		"application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := fetch("POST", url, header, "id="+body)
	if err != nil {
		log.Print("fetch:", err)
		return ItemDetail{}, err
	}
	//	log.Println("status:", resp.Status)

	buf, err := readRespBody(resp, true)
	if err != nil {
		log.Print("readRespBody:", err)
		return ItemDetail{}, err
	}
	//	log.Printf("%s\n", buf)

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
