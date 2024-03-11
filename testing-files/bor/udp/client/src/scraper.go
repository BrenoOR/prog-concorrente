package src

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Slice_CS struct {
	MutexSlc sync.Mutex
	Slc      []string
}

func (s *Slice_CS) Append(str string) {
	//fmt.Println("Appending: ", str)
	if s.Contains(str) {
		return
	}
	s.MutexSlc.Lock()
	s.Slc = append(s.Slc, str)
	s.MutexSlc.Unlock()
}

func (s *Slice_CS) Get() []string {
	s.MutexSlc.Lock()
	defer s.MutexSlc.Unlock()
	return s.Slc
}

func (s *Slice_CS) Remove(str string) {
	//fmt.Println("Removing: ", str)
	s.MutexSlc.Lock()
	defer s.MutexSlc.Unlock()
	for i, v := range s.Slc {
		if v == str {
			s.Slc = append(s.Slc[:i], s.Slc[i+1:]...)
		}
	}
}

func (s *Slice_CS) Dequeue() (string, bool) {
	s.MutexSlc.Lock()
	defer s.MutexSlc.Unlock()
	if len(s.Slc) > 0 {
		str := s.Slc[0]
		s.Slc = s.Slc[1:]
		//fmt.Println("Dequeuing", str, true)
		return str, true
	}
	//fmt.Println("Dequeuing", "", false)
	return "", false
}

func (s *Slice_CS) Pop() (string, bool) {
	s.MutexSlc.Lock()
	defer s.MutexSlc.Unlock()
	if len(s.Slc) > 0 {
		str := s.Slc[len(s.Slc)-1]
		s.Slc = s.Slc[:len(s.Slc)-1]
		//fmt.Println("Popping", str, true)
		return str, true
	}
	//fmt.Println("Popping", "", false)
	return "", false
}

func (s *Slice_CS) Contains(str string) bool {
	s.MutexSlc.Lock()
	defer s.MutexSlc.Unlock()
	for _, v := range s.Slc {
		if v == str {
			return true
		}
	}
	return false
}

type Author struct {
	Name      string
	Birthdate string
	Location  string
	Bio       string
}

func (a Author) Print() {
	fmt.Println("Name: ", a.Name)
	fmt.Println("Birthdate: ", a.Birthdate)
	fmt.Println("Location: ", a.Location)
	fmt.Println("Bio: ", a.Bio)
}

type AuthorSlc struct {
	MutexAuthors sync.Mutex
	Authors      []Author
}

func (s *AuthorSlc) Append(q Author) {
	s.MutexAuthors.Lock()
	s.Authors = append(s.Authors, q)
	s.MutexAuthors.Unlock()
}

func (s *AuthorSlc) Get() []Author {
	s.MutexAuthors.Lock()
	defer s.MutexAuthors.Unlock()
	return s.Authors
}

func (s *AuthorSlc) Remove(a Author) {
	s.MutexAuthors.Lock()
	defer s.MutexAuthors.Unlock()
	for i, v := range s.Authors {
		if v.Name == a.Name {
			s.Authors = append(s.Authors[:i], s.Authors[i+1:]...)
		}
	}
}

type Quote struct {
	Quote  string
	Author string
	Tags   []string
	About  string
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

type QuoteSlc struct {
	MutexQuotes sync.Mutex
	Quotes      []Quote
}

func (s *QuoteSlc) Append(q Quote) {
	s.MutexQuotes.Lock()
	s.Quotes = append(s.Quotes, q)
	s.MutexQuotes.Unlock()
}

func (s *QuoteSlc) Get() []Quote {
	s.MutexQuotes.Lock()
	defer s.MutexQuotes.Unlock()
	return s.Quotes
}

func (s *QuoteSlc) Remove(q Quote) {
	s.MutexQuotes.Lock()
	defer s.MutexQuotes.Unlock()
	for i, v := range s.Quotes {
		if v.Quote == q.Quote {
			s.Quotes = append(s.Quotes[:i], s.Quotes[i+1:]...)
		}
	}
}

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

type PokemonSlc struct {
	MutexPokemon sync.Mutex
	Pokemons     []Pokemon
}

func (s *PokemonSlc) Append(p Pokemon) {
	if p.Info.Name == "" {
		return
	}
	//fmt.Println("Appending: ", p.Info.Name)
	if s.Contains(p) {
		return
	}
	s.MutexPokemon.Lock()
	s.Pokemons = append(s.Pokemons, p)
	s.MutexPokemon.Unlock()
}

func (s *PokemonSlc) Get() []Pokemon {
	s.MutexPokemon.Lock()
	defer s.MutexPokemon.Unlock()
	return s.Pokemons
}

func (s *PokemonSlc) Remove(p Pokemon) {
	s.MutexPokemon.Lock()
	defer s.MutexPokemon.Unlock()
	for i, v := range s.Pokemons {
		if v.Info.Name == p.Info.Name {
			s.Pokemons = append(s.Pokemons[:i], s.Pokemons[i+1:]...)
		}
	}
}

func (s *PokemonSlc) Contains(p Pokemon) bool {
	s.MutexPokemon.Lock()
	defer s.MutexPokemon.Unlock()
	for _, v := range s.Pokemons {
		if v.Info.Name == p.Info.Name {
			return true
		}
	}
	return false
}

func Scrape(url string, connType string, URLsToVisit *Slice_CS, URLsVisited *Slice_CS, quotes *QuoteSlc, authors *AuthorSlc, pokemons *PokemonSlc, ch *chan int, wg *sync.WaitGroup, rttMutex *sync.Mutex, rttMean *int64, clientUDP *net.UDPConn, clientTCP *net.TCPConn, clientRPC *rpc.Client, consumer *amqp.Connection, clientMutex *sync.Mutex) {
	defer wg.Done()
	defer func() {
		<-*ch
	}()

	res := make([]byte, 500*1024)
	switch connType {
	case "tcp":
		getPageTCP(url, &res, rttMutex, rttMean)
	case "udp":
		getPageUDP(url, &res, rttMutex, rttMean)
	case "rpc":
		getPageGoRPC(url, &res, rttMutex, rttMean, clientRPC, clientMutex)
	case "rabbitmq":
		getPageRabbitMQ(url, &res, rttMutex, rttMean, consumer)
	default:
		log.Fatal("Invalid connection type")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".quote").Each(func(i int, s *goquery.Selection) {
		q := Quote{}
		var exists bool

		q.Quote = s.Find(".text").Text()
		q.Author = s.Find(".author").Text()
		q.About, exists = s.Find(".author").Next().Attr("href")
		q.Tags = []string{}
		s.Find(".tag").Each(func(i int, s *goquery.Selection) {
			q.Tags = append(q.Tags, s.Text())
		})

		quotes.Append(q)
		if exists {
			if !URLsVisited.Contains("http://quotes.toscrape.com"+q.About) && !URLsToVisit.Contains("http://quotes.toscrape.com"+q.About) {
				URLsToVisit.Append("http://quotes.toscrape.com" + q.About)
			}
		}
	})

	doc.Find(".author-details").Each(func(i int, s *goquery.Selection) {
		a := Author{}

		a.Name = s.Find("author-title").Text()
		a.Birthdate = s.Find(".author-born-date").Text()
		a.Location = s.Find(".author-born-location").Text()
		a.Bio = s.Find(".author-description").Text()

		authors.Append(a)
	})

	doc.Find(".storefront-sorting").Each(func(i int, s *goquery.Selection) {
		for _, page := range s.Find(".page-numbers").Nodes {
			for _, attr := range page.Attr {
				if attr.Key == "href" {
					if !URLsVisited.Contains(attr.Val) && !URLsToVisit.Contains(attr.Val) {
						URLsToVisit.Append(attr.Val)
					}
				}
			}
		}
	})

	doc.Find(".columns-4").Each(func(i int, s *goquery.Selection) {
		for _, pokemon := range s.Find(".woocommerce-LoopProduct-link").Nodes {
			for _, attr := range pokemon.Attr {
				if attr.Key == "href" {
					if !URLsVisited.Contains(attr.Val) && !URLsToVisit.Contains(attr.Val) {
						URLsToVisit.Append(attr.Val)
					}
				}
			}
		}
	})

	doc.Find(".product").Each(func(i int, s *goquery.Selection) {
		p := Pokemon{}
		pInfo := PokemonInfo{}

		for _, summary := range s.Find(".summary").Children().Nodes {
			for _, attr := range summary.FirstChild.Attr {
				if attr.Val == "woocommerce-Price-amount amount" {
					p.Price = summary.FirstChild.FirstChild.FirstChild.Data + summary.FirstChild.LastChild.Data
				}
			}
		}
		p.Quantity = s.Find(".in-stock").Text()
		p.SKU = s.Find(".sku").Text()
		for _, tag := range s.Find(".tag").Nodes {
			p.Tags = append(p.Tags, tag.FirstChild.Data)
		}

		pInfo.Name = s.Find(".product_title").Text()
		pInfo.Description = s.Find(".woocommerce-product-details__short-description").Text()
		pInfo.Weight = s.Find(".woocommerce-product-attributes-item--weight").Text()
		pInfo.Dimensions = s.Find(".woocommerce-product-attributes-item--dimensions").Text()
		for _, category := range s.Find(".posted_in").Children().Nodes {
			if category.FirstChild.Data != "Pokemon" {
				pInfo.Categories = append(pInfo.Categories, category.FirstChild.Data)
			}
		}
		p.Info = pInfo

		if pInfo.Name != "" && !pokemons.Contains(p) {
			if strings.Contains(url, pInfo.Name) {
				fmt.Println("Appending valid:", pInfo.Name)
			}
			pokemons.Append(p)
		}
	})

	URLsVisited.Append(url)
}

func ScrapeNC(url string, connType string, URLsToVisit *Slice_CS, URLsVisited *Slice_CS, quotes *QuoteSlc, authors *AuthorSlc, pokemons *PokemonSlc, rttMutex *sync.Mutex, rttMean *int64, clientUDP *net.UDPConn, clientTCP *net.TCPConn, clientRPC *rpc.Client, consumer *amqp.Connection, clientMutex *sync.Mutex) {
	res := make([]byte, 500*1024) // buffer size 500KB
	//fmt.Println("Scraping:", url, "with", connType)
	switch connType {
	case "tcp":
		getPageTCP(url, &res, rttMutex, rttMean)
	case "udp":
		getPageUDP(url, &res, rttMutex, rttMean)
	case "rpc":
		getPageGoRPC(url, &res, rttMutex, rttMean, clientRPC, clientMutex)
	case "rabbitmq":
		getPageRabbitMQ(url, &res, rttMutex, rttMean, consumer)
	default:
		log.Fatal("Invalid connection type")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".quote").Each(func(i int, s *goquery.Selection) {
		q := Quote{}
		var exists bool

		q.Quote = s.Find(".text").Text()
		q.Author = s.Find(".author").Text()
		q.About, exists = s.Find(".author").Next().Attr("href")
		q.Tags = []string{}
		s.Find(".tag").Each(func(i int, s *goquery.Selection) {
			q.Tags = append(q.Tags, s.Text())
		})

		quotes.Append(q)
		if exists {
			if !URLsVisited.Contains("http://quotes.toscrape.com"+q.About) && !URLsToVisit.Contains("http://quotes.toscrape.com"+q.About) {
				URLsToVisit.Append("http://quotes.toscrape.com" + q.About)
			}
		}
	})

	doc.Find(".author-details").Each(func(i int, s *goquery.Selection) {
		a := Author{}

		a.Name = s.Find("author-title").Text()
		a.Birthdate = s.Find(".author-born-date").Text()
		a.Location = s.Find(".author-born-location").Text()
		a.Bio = s.Find(".author-description").Text()

		authors.Append(a)
	})

	doc.Find(".storefront-sorting").Each(func(i int, s *goquery.Selection) {
		for _, page := range s.Find(".page-numbers").Nodes {
			for _, attr := range page.Attr {
				if attr.Key == "href" {
					if !(URLsVisited.Contains(attr.Val) || URLsToVisit.Contains(attr.Val)) {
						//fmt.Println("Appending: ", attr.Val)
						URLsToVisit.Append(attr.Val)
					}
				}
			}
		}
	})

	doc.Find(".columns-4").Each(func(i int, s *goquery.Selection) {
		for _, pokemon := range s.Find(".woocommerce-LoopProduct-link").Nodes {
			for _, attr := range pokemon.Attr {
				if attr.Key == "href" {
					if !(URLsVisited.Contains(attr.Val) || URLsToVisit.Contains(attr.Val)) {
						//fmt.Println("Appending: ", attr.Val)
						URLsToVisit.Append(attr.Val)
					}
				}
			}
		}
	})

	doc.Find(".product").Each(func(i int, s *goquery.Selection) {
		p := Pokemon{}
		pInfo := PokemonInfo{}

		for _, summary := range s.Find(".summary").Children().Nodes {
			for _, attr := range summary.FirstChild.Attr {
				if attr.Val == "woocommerce-Price-amount amount" {
					p.Price = summary.FirstChild.FirstChild.FirstChild.Data + summary.FirstChild.LastChild.Data
				}
			}
		}
		p.Quantity = s.Find(".in-stock").Text()
		p.SKU = s.Find(".sku").Text()
		for _, tag := range s.Find(".tag").Nodes {
			p.Tags = append(p.Tags, tag.FirstChild.Data)
		}

		pInfo.Name = s.Find(".product_title").Text()
		pInfo.Description = s.Find(".woocommerce-product-details__short-description").Text()
		pInfo.Weight = s.Find(".woocommerce-product-attributes-item--weight").Text()
		pInfo.Dimensions = s.Find(".woocommerce-product-attributes-item--dimensions").Text()
		for _, category := range s.Find(".posted_in").Children().Nodes {
			if category.FirstChild.Data != "Pokemon" {
				pInfo.Categories = append(pInfo.Categories, category.FirstChild.Data)
			}
		}
		p.Info = pInfo

		if pInfo.Name != "" && !pokemons.Contains(p) {
			pokemons.Append(p)
		}
	})

	URLsVisited.Append(url)
}
