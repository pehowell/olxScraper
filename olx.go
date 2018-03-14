package main

import (
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/soup"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Item struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadItems(filename string) []Item {
	var items []Item
	if _, err := os.Stat(filename); err == nil {
		data, error := ioutil.ReadFile(filename)
		check(error)
		json.Unmarshal(data, &items)
	}
	return items
}

func saveItems(items []Item, filename string) {
	data, err := json.Marshal(&items)
	check(err)
	error := ioutil.WriteFile(filename, data, 0644)
	check(error)
}

func main() {
	if len(os.Args) != 2 {
		panic("Must provide json file")
	}
	filename := os.Args[1]
	items := loadItems(filename)
	var foundItems []Item
	var newItems []Item
	re, _ := regexp.Compile("[0-9]+")
	resp, err := soup.Get("https://www.olx.ph/all-results?q=ags-101")
	check(err)

	doc := soup.HTMLParse(resp)
	results := doc.FindAll("div", "itemid", "#product")
	for _, result := range results {
		name := result.Find("span", "itemprop", "name")
		priceLine := result.Find("div", "itemprop", "offers").Find("span", "class", "price")
		price, err := strconv.Atoi(re.FindString(strings.Replace(priceLine.Text(), ",", "", -1)))
		check(err)
		foundItems = append(foundItems, Item{name.Text(), price})
	}
	for _, foundItem := range foundItems {
		known := false
		for _, item := range items {
			if foundItem == item {
				known = true
			}
		}
		if !known {
			newItems = append(newItems, foundItem)
		}
	}
	saveItems(foundItems, filename)
	if len(newItems) > 0 {
		for _, item := range newItems {
			fmt.Printf("%v -> %v\n", item.Name, item.Price)
		}
	}
}
