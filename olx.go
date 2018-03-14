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
	Link  string `json:"link"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func withinRange(amount int, min int, max int) bool {
	// Check to see if amount is at least the minimum and at most the maximum
	if amount >= min {
		if max <= 0 {
			return true
		}
		if amount <= max {
			return true
		}
	}
	return false
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
  olx [-hp] [-a <min price>] [-b <max price>] <pattern> <jsonfile>

Options:
  -a <price>      Price must be above this amount.
  -b <price>      Price must be below this amount.
  -h, --help      Show this screen.
  -p, --print     Print new items.`

	args, _ := docopt.ParseDoc(usage)

	// Price range arguments to narrow search
	above, _ := args.Int("-a")
	below, _ := args.Int("-b")
	if above != 0 && below != 0 {
		// Above (min) price can not be more than below (max) price
		if above > below {
			panic("Nonsensical use of above and below price range.")
		}
		if above < 0 || below < 0 {
			panic("Use of negative numbers for above or below is disallowed")
		}
	}
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
		if withinRange(price, above, below) {
			foundItems = append(foundItems, Item{name.Text(), price, itemURL.String()})
		}
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
