package crawler

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	//"github.com/sonhnguyen/pcchecker/mlabConnector"
	. "github.com/sonhnguyen/pcchecker/model"
)

func ScrapeTanDoanh(res []PcItem) ([]PcItem, error) {
	ROOT_URL := "http://tandoanh.vn"
	doc, err := goquery.NewDocument(ROOT_URL + "/baogia.php")
	if err != nil {
		return nil, err
	}
	doc.Find("div.accordion table tbody tr").Each(func(i int, s *goquery.Selection) {
		category := s.Closest(".accordion").PrevFiltered("h3").Text()
		category = category[14:]
		if s.Find("td:nth-child(1)").Text() != "STT" {

			price := 0
			col2 := s.Find("td:nth-child(2)").Text()
			col3 := s.Find("td:nth-child(3)").Text()
			col3 = col3[:len(col3)-3]
			col3 = strings.Replace(col3, ".", "", -1)

			price, err = strconv.Atoi(col3)
			if err != nil {
				price = 0
			}
			col4 := s.Find("td:nth-child(4)").Text()
			col5 := s.Find("td:nth-child(5)").Text()
			item := PcItem{Title: col2, Price: price, Guarantee: col4, Available: col5, Vendor: "tandoanh", Category: category}
			res = append(res, item)
			fmt.Println("TanDoanh reading", len(res))
		}
	})
	return res, nil
}

func ScrapeTanDoanhVer2(res []PcItem) ([]PcItem, error) {
	ROOT_URL := "http://tandoanh.vn"
	doc, err := goquery.NewDocument(ROOT_URL + "/trang-chu/")
	if err != nil {
		return nil, err
	}
	counting := 0
	doc.Find("tr td.menutitle a").Each(func(i int, s *goquery.Selection) {
		hrefCategory, _ := s.Attr("href")
		categoryTitle := s.Text()
		if hrefCategory != "javascript:void(0)" {
			docCategory, errCategory := goquery.NewDocument(ROOT_URL + hrefCategory)
			if errCategory != nil {
				return
			} else {
				docCategory.Find("table tbody tr td table tbody tr td table tbody tr td table tbody tr td table tbody tr td table tbody tr td table tbody tr td table tbody tr td a").Each(func(iCategory int, sCategory *goquery.Selection) {
					hrefItem, _ := sCategory.Attr("href")
					docItem, errItem := goquery.NewDocument(ROOT_URL + hrefItem)
					if errItem != nil {
						return
					} else {
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
						item := PcItem{Title: title, Price: price, Guarantee: guarantee, Image: imageList, Available: available, Vendor: "tandoanh", Category: categoryTitle, Link: ROOT_URL + hrefItem}
						counting++
						fmt.Println("Done item ", counting, item.Title, item.Link)
						res = append(res, item)
					}
				})
			}
		}
	})
	return res, nil
}

func ScrapeHH(res []PcItem) ([]PcItem, error) {
	ROOT_URL := "http://huuhoang.com"
	doc, err := goquery.NewDocument(ROOT_URL + "/ban-phim/")
	if err != nil {
		return nil, err
	}
	productsLink := []string{}
	doc.Find("li[class^='cat-']").Each(func(i int, s *goquery.Selection) {
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
					productsLink = append(productsLink, ROOT_URL+productLink)
				}
			})
		}

		//fetching pagination:

	})
	//have productsLink contains all the products

	for i := 0; i < len(productsLink); i++ {
		doc, err := goquery.NewDocument(productsLink[i])
		if err == nil {
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
			item := PcItem{Title: title, Link: productsLink[i], Price: price, Vendor: "huuhoang", Category: category, Desc: desc, Image: images}
			res = append(res, item)
			fmt.Println("HH reading %v / %v", i, len(productsLink))
		}

	}

	return res, nil
}
func ScrapeGamebank(res []PcItem) ([]PcItem, error) {
	ROOT_URL := "https://gear.gamebank.vn"
	doc, err := goquery.NewDocument(ROOT_URL + "/")
	if err != nil {
		return nil, err
	}
	productsLink := []string{}

	doc.Find("ul.navbar-nav > li").Each(func(i int, s *goquery.Selection) {
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

		for i := 0; i < len(pagination); i++ {
			catPage, err2 := goquery.NewDocument(pagination[i])
			if err2 != nil {
				return
			}

			catPage.Find("div.product-thumb > div.image > a").Each(func(i int, s *goquery.Selection) {
				productLink, err := s.Attr("href")
				if err && productLink != "" {
					productsLink = append(productsLink, productLink)
				}
			})
		}
		fmt.Printf("%v", len(productsLink))
	})

	//have productsLink contains all the products
	for i := 0; i < len(productsLink); i++ {
		doc, err := goquery.NewDocument(productsLink[i])
		if err == nil {
			images := []string{}
			category := doc.Find("ul.breadcrumb li:nth-child(2) a").Text()
			title := doc.Find("div#content h1").Text()
			desc := doc.Find("div#tab-description").Text()
			image, err := doc.Find("img#zoomImg").Attr("src")
			if err && image != "" {
				images = append(images, image)
			}
			priceString := doc.Find("span.price-new").Text()
			priceString = strings.Replace(priceString, "Giá", "", -1)
			priceString = strings.Replace(priceString, ":", "", -1)
			priceString = strings.Replace(priceString, ".", "", -1)
			priceString = strings.Replace(priceString, "đ", "", -1)
			priceString = strings.Replace(priceString, " ", "", -1)

			price, err2 := strconv.Atoi(priceString)
			if err2 != nil {
				price = 0
			}

			desc = desc + doc.Find("div#product-content-tab").Text()
			available := doc.Find("div#content ul.list-unstyled > li:nth-child(1)").Text()
			origin := doc.Find("div#content ul.list-unstyled > li:nth-child(2) a").Text()
			guarantee := doc.Find("div#content ul.list-unstyled > li:nth-child(3)").Text()

			item := PcItem{Title: title, Link: productsLink[i], Price: price, Vendor: "gamebank", Category: category, Desc: desc, Image: images, Available: available, Origin: origin, Guarantee: guarantee}
			res = append(res, item)
			fmt.Println("gamebank reading %#v / %#v", i, len(productsLink))
		}

	}

	return res, nil
}

func ScrapeGearvn(res []PcItem) ([]PcItem, error) {
	ROOT_URL := "https://gearvn.com"

	categoryLinks := []string{"http://gearvn.com/collections/ban-phim-co-gaming/",
		"http://gearvn.com/collections/gaming-mouse/",
		"http://gearvn.com/collections/headphones/",
		"http://gearvn.com/collections/mouse-pad/",
		"http://gearvn.com/collections/ghe-choi-game/",
		"http://gearvn.com/collections/linh-kien-may-tinh/",
		"http://gearvn.com/collections/laptop-gaming-1/",
		"http://gearvn.com/collections/phu-kien/"}
	productsLink := []string{}

	//creating unsafe connection
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	for i := 0; i < len(categoryLinks); i++ {
		doc, err := goquery.NewDocument(categoryLinks[i])
		if err != nil {
			return nil, err
		}
		for {
			doc.Find("div.product-row > a").Each(func(i int, s *goquery.Selection) {
				productLink, _ := s.Attr("href")
				if productLink != "" {
					productsLink = append(productsLink, ROOT_URL+productLink)
				}
			})
			nextPage, _ := doc.Find("ul.pagination-list > li:last-child a").Attr("href")
			if nextPage != "" && (ROOT_URL+nextPage) != doc.Url.String() {
				web, err := client.Get(ROOT_URL + nextPage)
				if err != nil {
					return nil, err
				}

				doc, err = goquery.NewDocumentFromResponse(web)
				if err != nil {
					return nil, err
				}
			} else {
				break
			}
		}
	}
	fmt.Println("gearvn reading %#v", len(categoryLinks))

	//have productsLink contains all the products
	for i := 0; i < len(productsLink); i++ {
		web, err := client.Get(productsLink[i])
		if err != nil {
			return nil, err
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
			item := PcItem{Title: title, Link: productsLink[i], Price: price, Status: status, Vendor: "gearvn", Category: category, Desc: desc, Image: images, Origin: origin, Guarantee: guarantee}
			res = append(res, item)
			fmt.Println("gearvn reading %#v / %#v", i, len(productsLink))
		}
	}
	return res, nil
}

func ScrapePCX(res []PcItem) ([]PcItem, error) {
	ROOT_URL := "https://phongcachxanh.vn"
	doc, err := goquery.NewDocument(ROOT_URL + "/shop/page/1")
	if err != nil {
		return nil, err
	}
	productsLink := []string{}
	for {
		doc.Find("div.oe_product_image > a").Each(func(i int, s *goquery.Selection) {
			productLink, _ := s.Attr("href")
			if productLink != "" {
				productsLink = append(productsLink, ROOT_URL+productLink)
			}
		})
		nextPage, _ := doc.Find("ul.pagination > li:last-child a").Attr("href")
		if nextPage != "" {
			doc, err = goquery.NewDocument(ROOT_URL + nextPage)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	//have productsLink contains all the products
	for i := 0; i < len(productsLink); i++ {
		doc, err := goquery.NewDocument(productsLink[i])
		if err == nil {
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
			item := PcItem{Title: title, Link: productsLink[i], Price: price, ShortDesc: shortDesc, Status: status, Vendor: "phongcachxanh", Category: category, Desc: desc, Image: images, Origin: origin, Guarantee: guarantee}
			res = append(res, item)
			fmt.Println("PhongCachXanh reading %#v / %#v, %#v", i, len(productsLink))
		}
	}
	return res, nil
}

func ScrapeAZ(chProduct chan PcItem, chFinished chan bool) {
	productLinkCh := make(chan string, 1000)
	ROOT_URL := "http://www.azaudio.vn"
	fmt.Print("productlink")

	//category, may need to update in future
	categoryLinks := []string{"http://www.azaudio.vn/audio", "http://www.azaudio.vn/gaming-gear", "http://www.azaudio.vn/loa", "http://www.azaudio.vn/may-tinh"}
	defer func() {
		// Notify that we're done after this function
		chFinished <- true
	}()

	for i := range categoryLinks {
		categoryLink := categoryLinks[i]
		go func() {
			catPage, err2 := goquery.NewDocument(categoryLink)
			if err2 != nil {
				return
			}
			for true {

				catPage.Find(".item-prd a.center-block").Each(func(i int, s *goquery.Selection) {
					productLink, err := s.Attr("href")
					if err && productLink != "" {
						fmt.Print("something", productLink)
						//productsLink = append(productsLink, productLink)
						productLinkCh <- productLink
					}
				})

				nextPage, err := catPage.Find("a.ajaxpagerlink").Attr("href")
				if err && nextPage != "" {
					nextPage = ROOT_URL + nextPage
					catPage, err2 = goquery.NewDocument(nextPage)
					fmt.Printf("link page", nextPage)
					if err2 != nil {
						return
					}
				} else {
					return
				}
			}

		}()
	}

	const workerCount = 100
	for i := 0; i < workerCount; i++ {
		go func() {
			for {
				select {
				case productLink := <-productLinkCh:
					doc, err := goquery.NewDocument(productLink)
					fmt.Print("productlink received", productLink)
					if err == nil {
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
						item := PcItem{Title: title, ShortDesc: shortDesc, Link: productLink, Price: price, Vendor: "azaudio", Category: category, Desc: desc, Image: images, Origin: origin, Guarantee: guarantee}
						//res = append(res, item)
						fmt.Println("azaudio reading %#v / %#v", len(productLink))
						chProduct <- item
					}
				}
			}
		}()
	}

}

func Run() {
	// var pcItems = []PcItem{}
	// pcItems, err := ScrapeTanDoanhVer2(pcItems)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("done tandoanh", len(pcItems))
	// pcItems, err = ScrapeGearvn(pcItems)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("done ScrapeGearvn", len(pcItems))
	// pcItems, err = ScrapeAZ(pcItems)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("done ScrapeAZ", len(pcItems))
	// pcItems, err = ScrapePCX(pcItems)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("done ScrapePCX", len(pcItems))
	// pcItems, err = ScrapeGamebank(pcItems)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("done ScrapeGamebank", len(pcItems))
	// pcItems, err = ScrapeHH(pcItems)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("done ScrapeHH", len(pcItems))

	// Channels
	chProduct := make(chan PcItem)
	chFinished := make(chan bool)

	go ScrapeAZ(chProduct, chFinished)
	// case vendor == "azaudio":
	// 	go ScrapeAZ(chProduct, chFinished)
	pcItems := []PcItem{}
	go func() {
		for {
			select {
			case pcItem := <-chProduct:
				pcItems = append(pcItems, pcItem)
				fmt.Print("new item", len(pcItems))

			case <-chFinished:
				fmt.Printf("len pc items %v", len(pcItems))
				//mlabConnector.InsertMlab(pcItems)
			}
		}
	}()
}
