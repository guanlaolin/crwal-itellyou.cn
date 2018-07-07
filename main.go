package main

import (
	"fmt"
	"log"
	"strconv"
)

const url = `https://msdn.itellyou.cn/`
const cataUrl = `https://msdn.itellyou.cn/Category/Index`
const langUrl = `https://msdn.itellyou.cn/Category/GetLang`
const listUrl = `https://msdn.itellyou.cn/Category/GetList`
const productUrl = `https://msdn.itellyou.cn/Category/GetProduct`

func init() {
	log.SetFlags(log.Llongfile)
}

func main() {
	log.Println("Crawling itellyou.cn...")
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

	log.Println("Craw finished.")
	log.Print(catas)
}

func bytes2float(b []byte) float64 {
	f, err := strconv.ParseFloat(string(b), 8)
	if err != nil {
		log.Fatalln("convert bytes to float error:", err)
	}

	return f
}
