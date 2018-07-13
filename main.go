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
	"fmt"
	"log"
	"strconv"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	LOG_INFO("Crawling itellyou.cn...")

	craw()
	// 获取分类信息
	catas, err := getCatasUrl(url)
	if err != nil {
		log.Fatal("getCatasUrl:", err)
	}

	// 获取每一类的项
	for _, cata := range catas {
		fmt.Printf("Fetching %s items url...\n", cata.Name)
		cata.Items, err = getCataItemsUrl(cataUrl, cata.Id)
		if err != nil {
			log.Print("getCataItemsUrl, Name:%s, ID:%s\n", cata.Name, cata.Id)
			continue
		}

		// 获取每一项的语言
		for _, item := range cata.Items {
			fmt.Printf("Fetching item %s language\n", item.Name)
			item.ItemLangs, err = getCataItemLangsUrl(langUrl, item.Id)
			if err != nil {
				log.Print("getCataItemsUrl, Name:%s, ID:%s\n", cata.Name, cata.Id)
				continue
			}

			// 获取每一种语言可下载的文件列表
			for _, list := range item.ItemLangs {
				fmt.Printf("Fetching item %s list\n", list.Lang)
				list.ItemSums, err = getCataItemLangListsUrl(listUrl,
					item.Id+"&lang="+list.Id+"&filter=true")
				if err != nil {
					log.Print("getCataItemLangListsUrl, Lang:%s, ID:%s\n", list.Lang, list.Id)
					continue
				}

				for _, summ := range list.ItemSums {
					fmt.Printf("Fetching item detail %s list\n", summ.Name)
					summ.Detail, err = getCataItemLangListDetail(productUrl, summ.Id)
					if err != nil {
						log.Print("getCataItemLangListDetail, Name:%s\n", summ.Name)
						continue
					}
				}
			}
		}
	}

	LOG_INFO("Craw finished.")
	log.Print(catas)
}

func bytes2float(b []byte) float64 {
	f, err := strconv.ParseFloat(string(b), 8)
	if err != nil {
		log.Fatalln("convert bytes to float error:", err)
	}

	return f
}
