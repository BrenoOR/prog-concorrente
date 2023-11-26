package scraper

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly"
)

// Common types
type URL struct {
	Url string
}

// Types for http://quotes.toscrape.com page scraping
type Quote struct {
	Quote  string
	Author string
	Tags   []string
	About  string
}

type Author struct {
	Name      string
	Birthdate string
	Location  string
	Bio       string
}

// Types for https://scrapeme.live/shop/ page scraping
type PokemonInfo struct {
	Name        string
	Categories  []string
	Description string
	Weight      string
	Dimensions  string
}

type Pokemon struct {
	Info     PokemonInfo
	Price    string
	Quantity string
	SKU      string
	Tags     []string
}

func (q Quote) Print() {
	fmt.Println("Quote: ", q.Quote)
	fmt.Println("Author: ", q.Author)
	fmt.Print("Tags: ")
	for _, tag := range q.Tags {
		fmt.Print(tag)
		fmt.Print(" ")
	}
	fmt.Println("")
	fmt.Println("About: ", q.About)
}

func (a Author) Print() {
	fmt.Println("Name: ", a.Name)
	fmt.Println("Birthdate: ", a.Birthdate)
	fmt.Println("Location: ", a.Location)
	fmt.Println("Bio: ", a.Bio)
}

func (p Pokemon) Print() {
	fmt.Println("Name: ", p.Info.Name)
	fmt.Println("Description: ", p.Info.Description)
	fmt.Println("Weight: ", p.Info.Weight)
	fmt.Println("Dimensions: ", p.Info.Dimensions)
	fmt.Print("Categories: ")
	for _, category := range p.Info.Categories {
		fmt.Print(category)
		fmt.Print(" ")
	}
	fmt.Println("")
	fmt.Println("Price: ", p.Price)
	fmt.Println("Quantity: ", p.Quantity)
	fmt.Println("SKU: ", p.SKU)
	for _, tag := range p.Tags {
		fmt.Print(tag)
		fmt.Print(" ")
	}
	fmt.Println("")
}

func (p Pokemon) CsvEntry() []string {
	categories := ""
	tags := ""
	for _, category := range p.Info.Categories {
		categories += category + " "
	}
	for _, tag := range p.Tags {
		tags += tag + " "
	}

	return []string{p.Info.Name, p.Info.Description, p.Info.Weight, p.Info.Dimensions, categories, p.Price, p.Quantity, p.SKU, tags}
}

func Scraper(url string, quotes *[]Quote, authors *[]Author, pokemons *[]Pokemon, visitedPages *[]URL, wg *sync.WaitGroup) {
	defer wg.Done()

	c := colly.NewCollector()

	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 20})

	c.OnRequest(func(req *colly.Request) {
		req.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
		fmt.Println("Visiting: ", req.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Response code: ", r.StatusCode)
		if r.StatusCode == 200 {
			*visitedPages = append(*visitedPages, URL{Url: r.Request.URL.String()})
		}

	})

	c.OnHTML(".quote", func(h *colly.HTMLElement) {
		q := Quote{}
		q.Quote = h.DOM.Find(".text").Text()
		q.Author = h.DOM.Find(".author").Text()
		q.About = h.DOM.Find(".author").Next().AttrOr("href", "Not found")
		for _, tag := range h.DOM.Find(".tag").Nodes {
			q.Tags = append(q.Tags, tag.FirstChild.Data)
		}
		if q.About != "Not found" {
			wg.Add(1)
			go Scraper("http://quotes.toscrape.com"+q.About, quotes, authors, pokemons, visitedPages, wg)
		}

		*quotes = append(*quotes, q)
	})

	// http://quotes.toscrape.com page scraping
	c.OnHTML(".author-details", func(h *colly.HTMLElement) {
		a := Author{}
		a.Name = h.DOM.Find(".author-title").Text()
		a.Birthdate = h.DOM.Find(".author-born-date").Text()
		a.Location = h.DOM.Find(".author-born-location").Text()
		a.Bio = h.DOM.Find(".author-description").Text()

		listed := false

		for i := range *authors {
			if (*authors)[i].Name == a.Name {
				listed = true
				return
			}
		}

		if !listed {
			*authors = append(*authors, a)
		}
	})

	// https://scrapeme.live/shop/ page scraping
	c.OnHTML(".storefront-sorting", func(h *colly.HTMLElement) {
		for _, page := range h.DOM.Find(".page-numbers").Nodes {
			for _, attr := range page.Attr {
				//fmt.Println(attr.Key, attr.Val)
				if attr.Key == "href" {
					wg.Add(1)
					go Scraper(attr.Val, quotes, authors, pokemons, visitedPages, wg)
				}
			}
		}
	})

	c.OnHTML(".columns-4", func(h *colly.HTMLElement) {
		for _, pokemon := range h.DOM.Find(".woocommerce-LoopProduct-link").Nodes {
			for _, attr := range pokemon.Attr {
				if attr.Key == "href" {
					wg.Add(1)
					go Scraper(attr.Val, quotes, authors, pokemons, visitedPages, wg)
				}
			}
		}
	})

	c.OnHTML(".product", func(h *colly.HTMLElement) {
		p := Pokemon{}
		pInfo := PokemonInfo{}

		for _, summary := range h.DOM.Find(".summary").Children().Nodes {
			for _, attr := range summary.FirstChild.Attr {
				if attr.Val == "woocommerce-Price-amount amount" {
					p.Price = summary.FirstChild.FirstChild.FirstChild.Data + summary.FirstChild.LastChild.Data
				}
			}
		}
		p.Quantity = h.DOM.Find(".in-stock").Text()
		p.SKU = h.DOM.Find(".sku").Text()
		for _, tag := range h.DOM.Find(".tag").Nodes {
			p.Tags = append(p.Tags, tag.FirstChild.Data)
		}

		pInfo.Name = h.DOM.Find(".product_title").Text()
		pInfo.Description = h.DOM.Find(".woocommerce-product-details__short-description").Text()
		pInfo.Weight = h.DOM.Find(".product_weight").Text()
		pInfo.Dimensions = h.DOM.Find(".product_dimensions").Text()
		for _, category := range h.DOM.Find(".posted_in").Children().Nodes {
			if category.FirstChild.Data != "Pokemon" {
				pInfo.Categories = append(pInfo.Categories, category.FirstChild.Data)
			}
		}

		p.Info = pInfo
		listed := false

		for i := range *pokemons {
			if (*pokemons)[i].Info.Name == p.Info.Name {
				listed = true
				return
			}
		}
		if pInfo.Name != "" {
			listed = true
		}

		if !listed {
			*pokemons = append(*pokemons, p)
		}
	})

	// Check if the page was already visited
	visited := false
	for _, page := range *visitedPages {
		if page.Url == url {
			//fmt.Println("Already visited: ", url)
			visited = true
			return
		}
	}

	if !visited {
		c.Visit(url)
	}
}
