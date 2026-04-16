package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

type Article struct {
	Title       string
	Description string
	Date        string
}

func printBanner() {
	banner := `
__  __                          __             
\ \/ /____ __   __ __  __ ____ / /   ____ _ _____
 \  // __ \\ \ / // / / //_  // /   / __ \// ___/
 / // /_/ / \ V // /_/ /  / /_ / /___/ /_/ // /    
/_/ \__,_/   \_/ \__,_/  /___//_____/\__,_//_/     
 _       __  _____  ______ ____   ____   __        
| |     / / / ___/ /_  __// __ \ / __ \ / /        
| | /| / /  \__ \   / /  / / / // / / // /         
| |/ |/ /  ___/ /  / /  / /_/ // /_/ // /___       
|__/|__/  /____/  /_/   \____/ \____//_____/       
`
	fmt.Println(banner)
}

func cleanText(s string) string {
	s = strings.ReplaceAll(s, "", "")
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}


func fetchWithChromedp(siteName, url, waitSelector, containerSelector, titleSelector, descSelector, dateSelector string) []Article {
	fmt.Printf("%s ten haberler cekiliyor\n", siteName)

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true), 
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	var html string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(waitSelector, chromedp.ByQuery), 
		chromedp.OuterHTML("html", &html),
	)

	if err != nil {
		log.Fatalf("\n[%s] Hata %v", siteName, err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalf("\n[%s] Hata %v", siteName, err)
	}

	var articles []Article
	count := 0

	doc.Find(containerSelector).Each(func(i int, s *goquery.Selection) {
		if count >= 100 { 
			return
		}

		title := cleanText(s.Find(titleSelector).First().Text())
		desc := cleanText(s.Find(descSelector).First().Text())
		date := cleanText(s.Find(dateSelector).First().Text())

		if title != "" {
			if len(desc) > 250 {
				desc = desc[:247] + "..."
			}
			articles = append(articles, Article{
				Title:       title,
				Description: desc,
				Date:        date,
			})
			count++
		}
	})

	return articles
}

// ---- Site 1: The Hacker News ----
func scrapeTheHackerNews() []Article {
	return fetchWithChromedp(
		"The Hacker News",
		"https://thehackernews.com/",
		".body-post",  
		".body-post",  
		".home-title", 
		".home-desc",  
		".h-datetime", 
	)
}

// ---- Site 2: Krebs on Security ----
func scrapeKrebsOnSecurity() []Article {
	return fetchWithChromedp(
		"Krebs on Security",
		"https://krebsonsecurity.com/",
		"article",          
		"article",          
		".entry-title",     
		".entry-content p", 
		".entry-date",      
	)
}

// ---- Site 3: Cyber Security News ----
func scrapeCyberSecurityNews() []Article {
	return fetchWithChromedp(
		"Cyber Security News",
		"https://cybersecuritynews.com/",
		".td_module_wrap", 
		".td_module_wrap", 
		".entry-title",    
		".td-excerpt",     
		".td-module-date", 
	)
}

// ---- Ekrana Basma ve Dosyaya Kaydetme ----
func processAndSaveArticles(siteName string, articles []Article, hideDate, hideDesc bool) {
	if len(articles) == 0 {
		fmt.Println("Haber bulunamadi.")
		return
	}

	today := time.Now().Format("02-01-2006")
	safeSiteName := strings.ReplaceAll(siteName, " ", "_")
	fileName := fmt.Sprintf("%s_%s.txt", safeSiteName, today)

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Dosya olusturulamadi: %v", err)
	}
	defer file.Close()

	output := func(format string, a ...interface{}) {
		text := fmt.Sprintf(format, a...)
		fmt.Print(text)
		file.WriteString(text)
	}

	output("\n==================================================\n")
	output("Site: %s\n", siteName)
	output("Cekilme Tarihi: %s\n", time.Now().Format("02-01-2006 15:04:05"))
	output("==================================================\n\n")

	for i, article := range articles {
		count := i + 1
		fmt.Printf("\033[35m%d.Haber\033[0m\n", count)
		file.WriteString(fmt.Sprintf("%d.Haber\n", count))
		
		fmt.Printf("\033[35mBaslik:\033[0m ")
		fmt.Printf("%s\n", article.Title)
		file.WriteString(fmt.Sprintf("Baslik: %s\n", article.Title))

		if !hideDesc && article.Description != "" {
			fmt.Printf("\033[34mAciklama:\033[0m\n")
			fmt.Printf("%s\n", article.Description)
			file.WriteString(fmt.Sprintf("Aciklama:\n%s\n", article.Description))
		}

		if !hideDate && article.Date != "" {
			fmt.Printf("\033[31mTarih:\033[0m\n")
			fmt.Printf("%s\n", article.Date)
			file.WriteString(fmt.Sprintf("Tarih:\n%s\n", article.Date))
		}

		fmt.Printf("--------------------------------------------------\n")
		file.WriteString("--------------------------------------------------\n")
	}

	fmt.Printf("\nIslem tamamlandi!\n\n")
}

func main() {
	site1 := flag.Bool("1", false, "displays the first news site (The Hacker News)")
	site2 := flag.Bool("2", false, "displays the second news site (Krebs on Security)")
	site3 := flag.Bool("3", false, "displays the third news site (Cyber Security News)")
	exitCmd := flag.Bool("4", false, "exits the application")
	hideDate := flag.Bool("date", false, "filters the date part")
	hideDesc := flag.Bool("description", false, "filters the description part")

	flag.Usage = func() {
		printBanner()
		fmt.Println("usage of scraper:")
		fmt.Println("  -1\n\tdisplays the first news site")
		fmt.Println("  -2\n\tdisplays the second news site")
		fmt.Println("  -3\n\tdisplays the third news site")
		fmt.Println("  -4\n\texits the application")
		fmt.Println("  -date\n\tfilters the date part")
		fmt.Println("  -description\n\tfilters the description part")
	}

	flag.Parse()

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	if *exitCmd {
		printBanner()
		fmt.Println("Uygulamadan cikis yapiliyor")
		os.Exit(0)
	}

	printBanner()

	if *site1 {
		articles := scrapeTheHackerNews()
		processAndSaveArticles("The Hacker News", articles, *hideDate, *hideDesc)
	}

	if *site2 {
		articles := scrapeKrebsOnSecurity()
		processAndSaveArticles("Krebs on Security", articles, *hideDate, *hideDesc)
	}

	if *site3 {
		articles := scrapeCyberSecurityNews()
		processAndSaveArticles("Cyber Security News", articles, *hideDate, *hideDesc)
	}
}