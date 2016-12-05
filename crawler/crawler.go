package crawler

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sonhnguyen/pcchecker/mlabConnector"
	. "github.com/sonhnguyen/pcchecker/model"
)

func GetPriceToday(price int) PriceToday {
	datetime := time.Now()
	return PriceToday{
		Price:    price,
		Datetime: datetime,
	}
}
func ScrapeTanDoanhVer2(chProduct chan PcItem, tandoanhFinished chan bool) {
	ROOT_URL := "http://tandoanh.vn"
	doc, err := goquery.NewDocument(ROOT_URL + "/trang-chu/")
	if err != nil {
		return
	}
	var linkGroup sync.WaitGroup

	doc.Find("tr td.menutitle a").Each(func(i int, s *goquery.Selection) {
		linkGroup.Add(1)
		go func() {
			hrefCategory, _ := s.Attr("href")
			categoryTitle := s.Text()
			if hrefCategory != "javascript:void(0)" {
				docCategory, errCategory := goquery.NewDocument(ROOT_URL + hrefCategory)
				if errCategory != nil {
					return
				}
				docCategory.Find("table tbody tr td table tbody tr td table tbody tr td table tbody tr td table tbody tr td table tbody tr td table tbody tr td table tbody tr td a").Each(func(iCategory int, sCategory *goquery.Selection) {
					hrefItem, _ := sCategory.Attr("href")
					docItem, errItem := goquery.NewDocument(ROOT_URL + hrefItem)
					if errItem != nil {
						return
					}
					var imageList []string
					var title string
					var status string
					var guarantee string
					var from string
					var available string
					var priceString string
					var price int
					docItem.Find("table tbody tr td table tbody tr td table tbody tr td table tbody tr td table tbody tr td table").Each(func(iItem int, sItem *goquery.Selection) {
						if iItem == 6 {
							sItem.Find("tbody tr td a img").Each(func(iImage int, sImage *goquery.Selection) {
								hrefImage, _ := sImage.Attr("src")
								if strings.Index(hrefImage, "jpg") > 0 {
									hrefImage = ROOT_URL + "/" + hrefImage[strings.Index(hrefImage, "upload/shop"):strings.Index(hrefImage, "jpg")+3]
								} else if strings.Index(hrefImage, "png") > 0 {
									hrefImage = ROOT_URL + "/" + hrefImage[strings.Index(hrefImage, "upload/shop"):strings.Index(hrefImage, "png")+3]
								}
								imageList = append(imageList, hrefImage)
							})
						} else if iItem == 7 {
							title = sItem.Find("tbody tr td p:nth-child(1) span.product_name_view").Text()
							status = sItem.Find("tbody tr td p:nth-child(2) span").Text()
							guarantee = sItem.Find("tbody tr td p:nth-child(3) span").Text()
							from = sItem.Find("tbody tr td p:nth-child(4) span").Text()
							available = sItem.Find("tbody tr td p:nth-child(5) span").Text()
							priceString = strings.Replace(strings.Replace(sItem.Find("tbody tr td b font").Text(), ",", "", -1), " ", "", -1)
							price, err = strconv.Atoi(priceString)
							if err != nil {
								price = 0
							}
						}
					})
					item := PcItem{Title: title, Price: GetPriceToday(price), Guarantee: guarantee, Image: imageList, Available: available, Vendor: "tandoanh", Category: categoryTitle, Link: ROOT_URL + hrefItem}
					chProduct <- item
				})
			}
			linkGroup.Done()
			return
		}()
	})
	linkGroup.Wait()
	tandoanhFinished <- true
	return
}

func ScrapeHH(chProduct chan PcItem, huuhoangFinished chan bool) {
	ROOT_URL := "http://huuhoang.com"
	doc, err := goquery.NewDocument(ROOT_URL + "/ban-phim/")
	if err != nil {
		return
	}
	productLinkCh := make(chan string, 10000)
	var linkGroup sync.WaitGroup

	doc.Find("li[class^='cat-']").Each(func(i int, s *goquery.Selection) {
		linkGroup.Add(1)
		go func() {

			categoryPage, err := s.Find("a").Attr("href")
			categoryLink := ""
			if !err {
				categoryLink = ""
			} else {
				categoryLink = ROOT_URL + categoryPage
			}

			catPage, err2 := goquery.NewDocument(categoryLink)
			if err2 != nil {
				return
			}
			pagination := []string{}
			//page 1
			pagination = append(pagination, categoryLink)

			catPage.Find(".pagination li:not(:last-child) a").Each(func(i int, s *goquery.Selection) {
				paginationLink, err := s.Attr("href")
				if err && paginationLink != "" {
					pagination = append(pagination, ROOT_URL+paginationLink)
				}
			})

			for i := 0; i < len(pagination); i++ {
				catPage, err2 := goquery.NewDocument(pagination[i])
				if err2 != nil {
					return
				}

				catPage.Find("div.detail-product-slider").Each(func(i int, s *goquery.Selection) {
					productLink, err := s.Find("h3 a").Attr("href")
					if err && productLink != "" {
						productLinkCh <- ROOT_URL + productLink
					}
				})
			}
			linkGroup.Done()
			return

		}()
	})
	linkGroup.Wait()
	close(productLinkCh)

	var productGroup sync.WaitGroup
	productGroup.Add(2)

	workerCount := 2
	for i := 0; i < workerCount; i++ {
		go func(productLinkCh <-chan string) {
			for productLink := range productLinkCh {
				doc, err := goquery.NewDocument(productLink)
				if err != nil {
					fmt.Println("ERROR", err)

					return
				}
				images := []string{}
				category := doc.Find("ul.breadcrums li:nth-child(2) a").Text()
				title := doc.Find("div.detail-header h1").Text()
				desc := doc.Find("div.detail-description").Text()
				doc.Find("ul#product-gallery li a").Each(func(i int, s *goquery.Selection) {
					image, err := s.Attr("data-image")
					if err && image != "" {
						images = append(images, ROOT_URL+image)
					}
				})
				if len(images) == 0 {
					image, err := doc.Find("img#product-main-image").Attr("src")
					if err && image != "" {
						images = append(images, ROOT_URL+image)
					}
				}
				priceString := doc.Find("div.price span").Text()
				priceString = strings.Replace(priceString, ".", "", -1)
				priceString = strings.Replace(priceString, "đ", "", -1)
				priceString = strings.Replace(priceString, " ", "", -1)

				price, err := strconv.Atoi(priceString)
				if err != nil {
					price = 0
				}

				desc = desc + doc.Find("div#product-content-tab").Text()
				item := PcItem{Title: title, Link: productLink, Price: GetPriceToday(price), Vendor: "huuhoang", Category: category, Desc: desc, Image: images}
				chProduct <- item
			}
			productGroup.Done()
		}(productLinkCh)
	}
	productGroup.Wait()
	huuhoangFinished <- true
	return
}

func ScrapeGamebank(chProduct chan PcItem, gamebankFinished chan bool) {
	ROOT_URL := "https://gear.gamebank.vn"
	productLinkCh := make(chan string, 10000)

	doc, err := goquery.NewDocument(ROOT_URL + "/")
	if err != nil {
		return
	}

	var linkGroup sync.WaitGroup

	doc.Find("ul.navbar-nav > li").Each(func(i int, s *goquery.Selection) {
		linkGroup.Add(1)
		go func() {
			categoryPage, err := s.Find("a").Attr("href")
			categoryLink := ""
			if !err {
				categoryLink = ""
			} else {
				categoryLink = categoryPage
			}

			catPage, err2 := goquery.NewDocument(categoryLink)
			if err2 != nil {
				return
			}
			pagination := []string{}
			//page 1
			pagination = append(pagination, categoryLink)

			maxPage, err := catPage.Find("ul.pagination > li:last-child a").Attr("href")
			if err && maxPage != "" {
				splitPage := strings.Split(maxPage, "=")
				maxPageInt, err := strconv.Atoi(splitPage[1])
				if err != nil {
					maxPageInt = 0
				}
				for i := 2; i <= maxPageInt; i++ {
					pagination = append(pagination, categoryLink+"?page="+strconv.Itoa(i))
				}
			}
			fmt.Println("gamebank", pagination)

			for i := 0; i < len(pagination); i++ {
				catPage, err2 := goquery.NewDocument(pagination[i])
				if err2 != nil {
					return
				}

				catPage.Find("div.product-thumb > div.image > a").Each(func(i int, s *goquery.Selection) {
					productLink, err := s.Attr("href")
					if err && productLink != "" {
						productLinkCh <- productLink
						fmt.Println("gamebank", productLink)
					}
				})
			}
			linkGroup.Done()
			return
		}()
	})
	linkGroup.Wait()
	close(productLinkCh)

	var productGroup sync.WaitGroup
	productGroup.Add(2)

	workerCount := 2
	for i := 0; i < workerCount; i++ {
		go func(productLinkCh <-chan string) {
			for productLink := range productLinkCh {
				doc, err := goquery.NewDocument(productLink)
				if err != nil {
					fmt.Println("ERROR", err)

					return
				}
				images := []string{}
				category := doc.Find("ul.breadcrumb li:nth-child(2) a").Text()
				title := doc.Find("div#content h1").Text()
				desc := doc.Find("div#tab-description").Text()
				image, err2 := doc.Find("img#zoomImg").Attr("src")
				if err2 && image != "" {
					images = append(images, image)
				}
				priceString := doc.Find("span.price-new").Text()
				priceString = strings.Replace(priceString, "Giá", "", -1)
				priceString = strings.Replace(priceString, ":", "", -1)
				priceString = strings.Replace(priceString, ".", "", -1)
				priceString = strings.Replace(priceString, "đ", "", -1)
				priceString = strings.Replace(priceString, " ", "", -1)

				price, err3 := strconv.Atoi(priceString)
				if err3 != nil {
					price = 0
				}

				desc = desc + doc.Find("div#product-content-tab").Text()
				available := doc.Find("div#content ul.list-unstyled > li:nth-child(1)").Text()
				origin := doc.Find("div#content ul.list-unstyled > li:nth-child(2) a").Text()
				guarantee := doc.Find("div#content ul.list-unstyled > li:nth-child(3)").Text()

				item := PcItem{Title: title, Link: productLink, Price: GetPriceToday(price), Vendor: "gamebank", Category: category, Desc: desc, Image: images, Available: available, Origin: origin, Guarantee: guarantee}
				chProduct <- item
				fmt.Println("item", item)
			}
			productGroup.Done()
		}(productLinkCh)
	}
	productGroup.Wait()
	gamebankFinished <- true

	return
}

func ScrapeGearvn(chProduct chan PcItem, gearvnFinished chan bool) {
	productLinkCh := make(chan string, 10000)
	ROOT_URL := "https://gearvn.com"

	categoryLinks := []string{"http://gearvn.com/collections/ban-phim-co-gaming/",
		"http://gearvn.com/collections/gaming-mouse/",
		"http://gearvn.com/collections/headphones/",
		"http://gearvn.com/collections/mouse-pad/",
		"http://gearvn.com/collections/ghe-choi-game/",
		"http://gearvn.com/collections/linh-kien-may-tinh/",
		"http://gearvn.com/collections/laptop-gaming-1/",
		"http://gearvn.com/collections/phu-kien/"}
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConnsPerHost: 50,
	}
	client := &http.Client{
		Transport: tr}
	var linkGroup sync.WaitGroup
	linkGroup.Add(len(categoryLinks))

	for i := 0; i < len(categoryLinks); i++ {
		categoryLink := categoryLinks[i]

		go func() {
			//creating unsafe connection
			web, err := client.Get(categoryLink)
			if err != nil {
				fmt.Println("ERROR", err)
				return
			}

			doc, err := goquery.NewDocumentFromResponse(web)
			if err != nil {
				fmt.Println("ERROR", err)
				return
			}
			for true {
				doc.Find("div.product-row > a").Each(func(i int, s *goquery.Selection) {

					productLink, _ := s.Attr("href")
					if productLink != "" {
						productLinkCh <- productLink
					}
				})

				nextPage, _ := doc.Find("ul.pagination-list > li:last-child a").Attr("href")
				if nextPage != "" && (ROOT_URL+nextPage) != doc.Url.String() {
					web, err := client.Get(ROOT_URL + nextPage)
					if err != nil {
						fmt.Println("ERROR", err)
						return
					}

					doc, err = goquery.NewDocumentFromResponse(web)
					if err != nil {
						fmt.Println("ERROR", err)
						return
					}

				} else {
					linkGroup.Done()
					return
				}
			}

		}()
	}
	linkGroup.Wait()
	close(productLinkCh)

	var productGroup sync.WaitGroup
	productGroup.Add(len(categoryLinks))

	workerCount := len(categoryLinks)
	for i := 0; i < workerCount; i++ {
		go func(productLinkCh <-chan string) {
			for productLink := range productLinkCh {
				web, err := client.Get(ROOT_URL + productLink)
				if err != nil {
					fmt.Println("ERROR", err)

					return
				}

				doc, err := goquery.NewDocumentFromResponse(web)
				if err == nil {
					images := []string{}
					category := doc.Find("#breadcrumb > div > div > span:nth-child(4) > a").Text()
					title := doc.Find("h1.product_name").Text()
					desc := doc.Find("div.tab-content").Text()
					doc.Find("div.product_thumbnail img").Each(func(i int, s *goquery.Selection) {
						src, exists := s.Attr("src")
						if exists {
							images = append(images, src)
						}
					})
					priceString := doc.Find("span.product_sale_price").First().Text()
					priceString = strings.Replace(priceString, ",", "", -1)
					priceString = strings.Replace(priceString, "₫", "", -1)
					price, err2 := strconv.Atoi(priceString)
					if err2 != nil {
						price = 0
					}
					status := doc.Find("div.product_parameters > p:nth-child(4) > span").Text()
					origin := doc.Find("div.product_parameters > p:nth-child(3) > span").Text()
					guarantee := doc.Find("div.product_parameters > p:nth-child(5) > span").Text()
					item := PcItem{Title: title, Link: ROOT_URL + productLink, Price: GetPriceToday(price), Status: status, Vendor: "gearvn", Category: category, Desc: desc, Image: images, Origin: origin, Guarantee: guarantee}
					chProduct <- item
				} else {
					fmt.Println("ERROR", err)
					return
				}
			}
			productGroup.Done()
		}(productLinkCh)
	}
	productGroup.Wait()
	gearvnFinished <- true
}

func ScrapePCX(chProduct chan PcItem, pcxFinished chan bool) {
	ROOT_URL := "https://phongcachxanh.vn"
	productLinkCh := make(chan string, 1000)
	doc, err := goquery.NewDocument(ROOT_URL + "/shop/page/1")
	if err != nil {
		return
	}

	var linkGroup sync.WaitGroup
	linkGroup.Add(1)

	go func() {
		for {
			doc.Find("div.oe_product_image > a").Each(func(i int, s *goquery.Selection) {
				productLink, _ := s.Attr("href")
				if productLink != "" {
					productLinkCh <- ROOT_URL + productLink
				}
			})

			nextPage, _ := doc.Find("ul.pagination > li:last-child a").Attr("href")
			if nextPage != "" {
				doc, err = goquery.NewDocument(ROOT_URL + nextPage)
				if err != nil {
					return
				}
			} else {
				linkGroup.Done()
				return
			}
		}
	}()
	linkGroup.Wait()
	close(productLinkCh)

	var productGroup sync.WaitGroup
	productGroup.Add(100)

	workerCount := 100
	for i := 0; i < workerCount; i++ {
		go func(productLinkCh <-chan string) {
			for productLink := range productLinkCh {
				doc, err := goquery.NewDocument(productLink)
				if err != nil {
					fmt.Println("ERROR", err)
					return
				}

				images := []string{}
				category := doc.Find("ol.breadcrumb li:nth-child(2) a").Text()
				title := doc.Find("div#product_details > h1").Text()
				shortDesc := doc.Find("div#product_details > div:nth-child(6) > p:nth-child(3)").Text()
				doc.Find("img.img.img-responsive.optima_thumbnail").Each(func(i int, s *goquery.Selection) {
					src, exists := s.Attr("src")
					if exists {
						images = append(images, ROOT_URL+src)
					}
				})
				priceString := doc.Find("span.oe_currency_value").First().Text()
				priceString = strings.Replace(priceString, ".", "", -1)
				price, err2 := strconv.Atoi(priceString)
				if err2 != nil {
					price = 0
				}
				status := "Mới 100%, Chính hãng"
				origin := doc.Find("div#product_details > div:nth-child(5) > b > span").Text()
				desc := shortDesc + doc.Find("div#product_full_description").Text()
				guarantee := doc.Find("div#product_details > b > span").Text()
				item := PcItem{Title: title, Link: productLink, Price: GetPriceToday(price), ShortDesc: shortDesc, Status: status, Vendor: "phongcachxanh", Category: category, Desc: desc, Image: images, Origin: origin, Guarantee: guarantee}
				chProduct <- item
			}
			productGroup.Done()

		}(productLinkCh)
	}
	productGroup.Wait()
	pcxFinished <- true
}

func ScrapeAZ(chProduct chan PcItem, azFinished chan bool) {
	productLinkCh := make(chan string, 10000)
	ROOT_URL := "http://www.azaudio.vn"
	//category, may need to update in future
	categoryLinks := []string{"http://www.azaudio.vn/audio", "http://www.azaudio.vn/gaming-gear", "http://www.azaudio.vn/loa", "http://www.azaudio.vn/may-tinh"}

	var linkGroup sync.WaitGroup
	linkGroup.Add(len(categoryLinks))

	for i := 0; i < len(categoryLinks); i++ {
		categoryLink := categoryLinks[i]
		go func() {
			catPage, err2 := goquery.NewDocument(categoryLink)
			if err2 != nil {
				fmt.Println("ERROR", err2)
				return
			}
			for true {
				catPage.Find(".item-prd a.center-block").Each(func(i int, s *goquery.Selection) {
					productLink, err := s.Attr("href")
					if err && productLink != "" {
						productLinkCh <- productLink
					}
				})

				nextPage, err := catPage.Find("a.ajaxpagerlink").Attr("href")
				if err && nextPage != "" {
					nextPage = ROOT_URL + nextPage
					catPage, err2 = goquery.NewDocument(nextPage)
					if err2 != nil {
						return
					}
				} else {
					linkGroup.Done()
					return
				}
			}
		}()
	}
	linkGroup.Wait()
	close(productLinkCh)

	var productGroup sync.WaitGroup
	productGroup.Add(len(categoryLinks))

	workerCount := len(categoryLinks)
	for i := 0; i < workerCount; i++ {
		go func(productLinkCh <-chan string) {
			for productLink := range productLinkCh {

				doc, err := goquery.NewDocument(productLink)
				if err != nil {
					fmt.Println("ERROR", err)

					return
				}

				images := []string{}
				category := doc.Find("a.itemcrumb.active > span").Text()
				shortDesc := doc.Find(".briefContent p").Text()
				title := doc.Find("div.prd-content h1").Text()
				desc := doc.Find(".contentFull").Text()
				origin := doc.Find(".prd-content .brands a").Text()
				guarantee := doc.Find(".prd-content div.guarantee").Text()
				doc.Find(".prd-detail img").Each(func(i int, s *goquery.Selection) {
					src, exists := s.Attr("src")
					if exists {
						images = append(images, ROOT_URL+src)
					}
				})
				priceString := doc.Find("span.new-price").Text()
				priceString = strings.Replace(priceString, ".", "", -1)
				priceString = strings.Replace(priceString, "₫", "", -1)
				priceString = strings.Replace(priceString, " ", "", -1)
				price, err2 := strconv.Atoi(priceString)
				if err2 != nil {
					price = 0
				}
				item := PcItem{Title: title, ShortDesc: shortDesc, Link: productLink, Price: GetPriceToday(price), Vendor: "azaudio", Category: category, Desc: desc, Image: images, Origin: origin, Guarantee: guarantee}
				chProduct <- item
			}
			productGroup.Done()
		}(productLinkCh)
	}
	productGroup.Wait()
	azFinished <- true
}

func CheckIsDoneAll(crawlerStatus map[string]bool, startTime time.Time) bool {
	for _, k := range crawlerStatus {
		if k == false {
			return false
		}
	}
	return true
}

func Run() {
	start := time.Now()

	// Channels
	chProduct := make(chan PcItem, 10000)
	azFinished := make(chan bool)
	pcxFinished := make(chan bool)
	gearvnFinished := make(chan bool)
	gamebankFinished := make(chan bool)
	huuhoangFinished := make(chan bool)
	tandoanhFinished := make(chan bool)

	// crawlerStatus := map[string]bool{"gamebank": false, "azaudio": false, "pcx": false, "gearvn": false}
	crawlerStatus := map[string]bool{"tandoanh": false}

	// go ScrapeGearvn(chProduct, gearvnFinished)
	// go ScrapeAZ(chProduct, azFinished)
	// go ScrapePCX(chProduct, pcxFinished)
	// go ScrapeGamebank(chProduct, gamebankFinished)
	// go ScrapeHH(chProduct, huuhoangFinished)
	go ScrapeTanDoanhVer2(chProduct, tandoanhFinished)
	pcItems := []PcItem{}

	const workerCount = 100
	for i := 0; i < workerCount; i++ {
		go func() {
			for {
				select {
				case pcItem := <-chProduct:
					pcItems = append(pcItems, pcItem)
					fmt.Println("new item", len(pcItems), pcItem.Link)
					break
				case <-azFinished:
					crawlerStatus["azaudio"] = true
					if CheckIsDoneAll(crawlerStatus, start) {
						elapsed := time.Since(start)
						fmt.Println("done everything %s", elapsed)
						mlabConnector.InsertMlab(pcItems)
						break
					}
					break
				case <-pcxFinished:
					crawlerStatus["pcx"] = true
					if CheckIsDoneAll(crawlerStatus, start) {
						elapsed := time.Since(start)
						fmt.Println("done everything %s", elapsed)
						mlabConnector.InsertMlab(pcItems)
						break
					}
					break
				case <-gearvnFinished:
					crawlerStatus["gearvn"] = true
					if CheckIsDoneAll(crawlerStatus, start) {
						elapsed := time.Since(start)
						fmt.Println("done everything %s", elapsed)
						mlabConnector.InsertMlab(pcItems)
						break
					}
					break
				case <-gamebankFinished:
					crawlerStatus["gamebank"] = true
					if CheckIsDoneAll(crawlerStatus, start) {
						elapsed := time.Since(start)
						fmt.Println("done everything %s", elapsed)
						mlabConnector.InsertMlab(pcItems)
						break
					}
					break
				case <-huuhoangFinished:
					crawlerStatus["huuhoang"] = true
					if CheckIsDoneAll(crawlerStatus, start) {
						elapsed := time.Since(start)
						fmt.Println("done everything %s", elapsed)
						mlabConnector.InsertMlab(pcItems)
						break
					}
					break
				case <-tandoanhFinished:
					crawlerStatus["tandoanh"] = true
					if CheckIsDoneAll(crawlerStatus, start) {
						elapsed := time.Since(start)
						fmt.Println("done everything %s", elapsed)
						mlabConnector.InsertMlab(pcItems)
						break
					}
					break
				}
			}
		}()
	}
}
