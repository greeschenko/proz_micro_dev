//setam sync system
package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	Dbuser string = "root"
	Dbpass string = "rootpass"
	Dbname string = "prozorrodev"
)

type Document struct {
	XMLName xml.Name `xml:"Document"`
	Items   Items    `xml:"items"`
}

type Items struct {
	XMLName xml.Name    `xml:"items"`
	Items   []SetamItem `xml:"item"`
}

// a simple struct which contains all our
// social links
type SetamItem struct {
	gorm.Model
	XMLName         xml.Name       `xml:"item"`
	Title           string         `xml:"title"`
	Link            string         `xml:"link"`
	Description     string         `xml:"description"`
	Category        int            `xml:"category"`
	Category_name   string         `xml:"category_name"`
	StartDate       string         `xml:"startDate"`
	EndDate         string         `xml:"endDate"`
	RequestsEndDate string         `xml:"requestsEndDate"`
	Region          int            `xml:"region"`
	Region_name     string         `xml:"region_name"`
	Seller          int            `xml:"seller"`
	Seller_name     string         `xml:"seller_name"`
	LotNumber       int            `xml:"lotNumber"`
	StartPrice      float32        `xml:"startPrice"`
	Enclosure       SetamEnclosure `xml:"enclosure"`
	EnclosureID     uint           `xml:"-"`
}

type SetamEnclosure struct {
	gorm.Model
	XMLName xml.Name `xml:"enclosure"`
	Url     string   `xml:"url,attr"`
	Length  string   `xml:"length,attr"`
	Type    string   `xml:"type,attr"`
}

func main() {
	url := "https://setam.net.ua/partners/data/pDK4PDXFnpo.xml"
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	//fmt.Println("response from GET request", res)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var doc Document

	d := xml.NewDecoder(bytes.NewReader(body))
	d.Entity = map[string]string{
		"laquo": "'",
		"raquo": "'",
	}
	decodeerr := d.Decode(&doc)
	if decodeerr != nil {
		fmt.Printf("error: %v", decodeerr)
		return
	}

	fmt.Println(len(doc.Items.Items))

	connectstr := fmt.Sprintf(
		"%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		Dbuser,
		Dbpass,
		Dbname,
	)
	db, err := gorm.Open("mysql", connectstr)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.AutoMigrate(&SetamItem{}, &SetamEnclosure{})

	var (
		T int = 20
		N int = len(doc.Items.Items)
		D int = N / T
	)

	ack := make(chan bool, N) // Acknowledgement channel

	for i := 0; i < N; i += D {
		var end int
		if i+D < N {
			end = i + D
		} else {
			end = N
		}

		fmt.Println(i, end)

		go func(i int) { // Point #1
			for j := i; j < end; j++ {
				var item = doc.Items.Items[j]
				db.Where("lot_number = ?", item.LotNumber).First(&item)
				if item.ID == 0 {
					db.Create(&item)
				} else {
					db.Model(&item).Updates(&item)
				}
				fmt.Println(j)
				ack <- true // Point #2
			}
		}(i) // Point #3
	}

	for i := 0; i < N; i++ {
		<-ack // Point #2
	}
}
