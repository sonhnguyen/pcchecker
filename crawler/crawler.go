package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
	. "github.com/sonhnguyen/pcchecker/model"	
	"github.com/sonhnguyen/pcchecker/mlabConnector"	
)


func ScrapeTanDoanh (res []PcItem) ([]PcItem, error) {
	ROOT_URL := "http://tandoanh.vn"
	doc, err := goquery.NewDocument(ROOT_URL + "/baogia.php")
	if err != nil {
		return nil, err
	}
	doc.Find("div.accordion table tbody tr").Each(func (i int, s *goquery.Selection) {
		category := s.Closest(".accordion").PrevFiltered("h3").Text()
		category = category[14:]
		if(s.Find("td:nth-child(1)").Text()!="STT"){

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

func ScrapeHH(res []PcItem) ([]PcItem, error) {
	ROOT_URL := "http://huuhoang.com"
	doc, err := goquery.NewDocument(ROOT_URL + "/ban-phim/")
	if err != nil {
		return nil, err
	}
	productsLink := []string{}
	doc.Find("li[class^='cat-']").Each(func (i int, s *goquery.Selection) {
		categoryPage, err := s.Find("a").Attr("href")
		categoryLink := ""
		if (!err) {
			categoryLink = ""
		}	else {
			categoryLink = ROOT_URL + categoryPage
		}

		catPage, err2 := goquery.NewDocument(categoryLink)
		if err2 != nil {
			return
		}
		pagination := []string{}
		//page 1
		pagination = append(pagination, categoryLink)

		catPage.Find(".pagination li:not(:last-child) a").Each(func (i int, s *goquery.Selection) {
			paginationLink, err := s.Attr("href")
			if (err && paginationLink !="") {
				pagination = append(pagination, ROOT_URL + paginationLink)
			}
		})

		for i := 0; i < len(pagination); i++ {
			catPage, err2 := goquery.NewDocument(pagination[i])
			if err2 != nil {
				return
			}

			catPage.Find("div.detail-product-slider").Each(func (i int, s *goquery.Selection) {
				productLink, err := s.Find("h3 a").Attr("href")
				if (err && productLink !="") {
					productsLink = append(productsLink, ROOT_URL + productLink)
				}
			})
		}

		//fetching pagination:


	})
	//have productsLink contains all the products

	for i := 0; i < len(productsLink); i++ {
		doc, err := goquery.NewDocument(productsLink[i])
		if (err == nil) {
			images := []string{}
			category:= doc.Find("ul.breadcrums li:nth-child(2) a").Text()
			title:= doc.Find("div.detail-header h1").Text()
			desc:= doc.Find("div.detail-description").Text()
			doc.Find("ul#product-gallery li a").Each(func (i int, s *goquery.Selection) {
				image, err := s.Attr("data-image")
				if (err && image!="" ){
					images = append(images, ROOT_URL + image)
				}
			})
			if(len(images)==0){
				image, err := doc.Find("img#product-main-image").Attr("src")
				if (err && image!="" ){
					images = append(images, ROOT_URL + image)
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
			item := PcItem{Title: title, Link:productsLink[i],Price: price, Vendor: "huuhoang", Category: category, Desc: desc, Image: images}
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

	doc.Find("ul.navbar-nav > li").Each(func (i int, s *goquery.Selection) {
		categoryPage, err := s.Find("a").Attr("href")
		categoryLink := ""
		if (!err) {
			categoryLink = ""
		}	else {
			categoryLink = categoryPage
		}

		catPage, err2 := goquery.NewDocument(categoryLink)
		if err2 != nil {
			return
		}
		pagination := []string{}
		//page 1
		pagination = append(pagination, categoryLink)

		maxPage,err := catPage.Find("ul.pagination > li:last-child a").Attr("href")
		if (err && maxPage !="") {
			splitPage := strings.Split(maxPage, "=")
			maxPageInt, err := strconv.Atoi(splitPage[1])
			if err != nil {
				maxPageInt = 0
			}
			for i := 2; i <= maxPageInt; i++ {
				pagination = append(pagination, categoryLink +"?page="+ strconv.Itoa(i))
			}
		}

		for i := 0; i < len(pagination); i++ {
			catPage, err2 := goquery.NewDocument(pagination[i])
			if err2 != nil {
				return
			}

			catPage.Find("div.product-thumb > div.image > a").Each(func (i int, s *goquery.Selection) {
				productLink, err := s.Attr("href")
				if (err && productLink !="") {
					productsLink = append(productsLink, productLink)
				}
			})
		}
		fmt.Printf("%v", len(productsLink))
	})

	//have productsLink contains all the products
	for i := 0; i < len(productsLink); i++ {
		doc, err := goquery.NewDocument(productsLink[i])
		if (err == nil) {
			images := []string{}
			category:= doc.Find("ul.breadcrumb li:nth-child(2) a").Text()
			title:= doc.Find("div#content h1").Text()
			desc:= doc.Find("div#tab-description").Text()
			image, err := doc.Find("img#zoomImg").Attr("src")
			if (err && image!="" ){
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
			available:= doc.Find("div#content ul.list-unstyled > li:nth-child(1)").Text()
			origin:= doc.Find("div#content ul.list-unstyled > li:nth-child(2) a").Text()
			guarantee:= doc.Find("div#content ul.list-unstyled > li:nth-child(3)").Text()

			item := PcItem{Title: title, Link:productsLink[i],Price: price, Vendor: "gamebank", Category: category, Desc: desc, Image: images, Available: available, Origin: origin, Guarantee: guarantee}
			res = append(res, item)
			fmt.Println("gamebank reading %#v / %#v", i, len(productsLink))
		}

	}

	return res, nil
}

func ScrapeAZ(res []PcItem) ([]PcItem, error) {
	ROOT_URL := "http://www.azaudio.vn"
	productsLink := []string{}

	//category, may need to update in future
	categoryLinks := []string{"http://www.azaudio.vn/audio", "http://www.azaudio.vn/gaming-gear", "http://www.azaudio.vn/loa", "http://www.azaudio.vn/may-tinh"}
	for i := 0; i < len(categoryLinks); i++ {
		catPage, err2 := goquery.NewDocument(categoryLinks[i])
		if err2 != nil {
			return nil, err2
		}
		for true {
			catPage.Find(".item-prd a.center-block").Each(func (i int, s *goquery.Selection) {
				productLink, err := s.Attr("href")
				if (err && productLink !="") {
					productsLink = append(productsLink, productLink)
				}
			})

			nextPage,err := catPage.Find("a.ajaxpagerlink").Attr("href")
			if (err && nextPage !="") {
				nextPage = ROOT_URL + nextPage
				catPage, err2 = goquery.NewDocument(nextPage)
				fmt.Printf("link page", nextPage)
				if err2 != nil {
					return nil, err2
				}
			}	else {
				break;
			}
		}
	}

	//have productsLink contains all the products
	for i := 0; i < len(productsLink); i++ {
		doc, err := goquery.NewDocument(productsLink[i])
		if (err == nil) {
			images := []string{}
			category:= doc.Find("a.itemcrumb.active > span").Text()
			shortDesc := doc.Find(".briefContent p").Text()
			title:= doc.Find("div.prd-content h1").Text()
			desc:= doc.Find(".contentFull").Text()
			origin:= doc.Find(".prd-content .brands a").Text()
			guarantee:= doc.Find(".prd-content div.guarantee").Text()
			doc.Find(".prd-detail img").Each(func(i int, s *goquery.Selection) {
		        src, exists := s.Attr("src")
		        if exists {
		        	images = append(images, ROOT_URL + src)
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
			item := PcItem{Title: title, ShortDesc: shortDesc, Link:productsLink[i],Price: price, Vendor: "azaudio", Category: category, Desc: desc, Image: images, Origin: origin, Guarantee: guarantee}
			res = append(res, item)
			fmt.Println("azaudio reading %#v / %#v", i, len(productsLink))
		}
	}
	return res, nil
}

func Run() {
	var pcItems = []PcItem{}
	pcItems, err := ScrapeTanDoanh(pcItems)
	if(err!=nil){
		fmt.Println(err)
	}
	fmt.Println(len(pcItems))
	// pcItems, err = ScrapeAZ(pcItems)
	// if(err!=nil){
	// 	fmt.Println(err)
	// }
	// pcItems, err = ScrapeGamebank(pcItems)
	// if(err!=nil){
	// 	fmt.Println(err)
	// }
	// pcItems, err = ScrapeHH(pcItems)
	// if(err!=nil){
	// 	fmt.Println(err)
	// }
	mlabConnector.InsertMlab(pcItems)
}