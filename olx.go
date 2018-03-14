package main

import (
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/docopt/docopt-go"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Item struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Link string `json:"link"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadItems(filename string) []Item {
	var items []Item
	if _, err := os.Stat(filename); err == nil {
		data, readErr := ioutil.ReadFile(filename)
		check(readErr)
		json.Unmarshal(data, &items)
	}
	return items
}

func saveItems(items []Item, filename string) {
	data, err := json.Marshal(&items)
	check(err)
	writeErr := ioutil.WriteFile(filename, data, 0644)
	check(writeErr)
}

func main() {
	usage := `Olx.ph Scraper

Usage:
  olx [-hp] <pattern> <jsonfile>

Options:
  -h, --help     Show this screen.
  -p, --print    Print new items.`

	args, _ := docopt.ParseDoc(usage)
  printNew, _ := args.Bool("--print")
  pattern, _ := args.String("<pattern>")
  filename, _ := args.String("<jsonfile>")

	olxURL := "https://www.olx.ph"

	var searchURL strings.Builder
	searchURL.WriteString(olxURL)
	searchURL.WriteString("/all-results?q=")
  searchURL.WriteString(pattern)

	items := loadItems(filename)
	var foundItems []Item
	var newItems []Item
	re, _ := regexp.Compile("[0-9]+")
	resp, err := soup.Get(searchURL.String())
	check(err)

	doc := soup.HTMLParse(resp)
	results := doc.FindAll("div", "itemid", "#product")
	for _, result := range results {
    var itemURL strings.Builder
    link := result.Find("a", "itemprop", "url")
    itemURL.WriteString(olxURL)
    itemURL.WriteString(link.Attrs()["href"])

		name := result.Find("span", "itemprop", "name")

		priceLine := result.Find("div", "itemprop", "offers").Find("span", "class", "price")
		price, err := strconv.Atoi(re.FindString(strings.Replace(priceLine.Text(), ",", "", -1)))
		check(err)

		foundItems = append(foundItems, Item{name.Text(), price, itemURL.String()})
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
	if printNew && len(newItems) > 0 {
		for _, item := range newItems {
			fmt.Printf("%v -> %v -> %v\n", item.Price, item.Name, item.Link)
		}
	}
}
