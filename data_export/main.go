// Package main provides ...
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	Domen string = "https://polonex.com.ua"
)

var (
	DB     *gorm.DB
	Dbuser string = "root"
	Dbpass string = "rootpass"
	Dbname string = "prozorrodev"
)

type Result struct {
	Id            string
	Proid         string
	Title         string
	Status        string
	AuctionID     string `gorm:"column:auctionID"`
	DgfID         string `gorm:"column:dgfID"`
	LotIdentifier string `gorm:"column:lotIdentifier"`
	PMT           string `gorm:"column:procurementMethodType"`
	Amount        string
	StartDate     string `gorm:"column:startDate"`
}

func getItemsAddresses(proid string) string {
	res := ""
	type AddressRes struct {
		Id            string
		Proid         string
		PostalCode    string `gorm:"column:postalCode"`
		CountryName   string `gorm:"column:countryName"`
		Region        string `gorm:"column:region"`
		Locality      string `gorm:"column:locality"`
		StreetAddress string `gorm:"column:streetAddress"`
	}
	var addrres []AddressRes

	DB.Raw(`SELECT
				proauction2_items.id,
				proauction2_items.proid,
				proauction2_addresses.postalCode,
				proauction2_addresses.countryName,
				proauction2_addresses.region,
				proauction2_addresses.locality,
				proauction2_addresses.streetAddress
			FROM proauction2_items
			LEFT JOIN proauction2_addresses
				ON proauction2_addresses.proid = proauction2_items.address
			WHERE proauction2_items.owner_id = ?`, proid).Scan(&addrres)

	for k, v := range addrres {
		res += fmt.Sprintf(
			"%d. %s %s %s %s %s",
			k,
			v.PostalCode,
			v.CountryName,
			v.Region,
			v.Locality,
			v.StreetAddress,
		)
		res += "\n\r"
	}

	return res
}

func getItemsClasif(proid string) string {
	res := ""
	type AddressRes struct {
		Proid       string
		Scheme      string
		Id          string
		Description string
	}
	var addrres []AddressRes

	DB.Raw(`SELECT
				proauction2_items.proid,
				proauction2_class.scheme,
				proauction2_class.id,
				proauction2_class.description
			FROM proauction2_items
			LEFT JOIN proauction2_class
				ON proauction2_class.proid = proauction2_items.classification
			WHERE proauction2_items.owner_id = ?`, proid).Scan(&addrres)

	for k, v := range addrres {
		res += fmt.Sprintf(
			"%d. %s %s - %s",
			k,
			v.Scheme,
			v.Id,
			v.Description,
		)
		res += "\n\r"
	}

	return res
}

func getItemsUnits(proid string) string {
	res := ""
	type AddressRes struct {
		Proid     string
		Quantity  string
		Unit_name string `gorm:"column:unit_name"`
	}
	var addrres []AddressRes

	DB.Raw(`SELECT
				proauction2_items.proid,
				proauction2_items.quantity,
				proauction2_items.unit_name
			FROM proauction2_items
			WHERE proauction2_items.owner_id = ?`, proid).Scan(&addrres)

	for k, v := range addrres {
		res += fmt.Sprintf(
			"%d. %s %s",
			k,
			v.Quantity,
			v.Unit_name,
		)
		res += "\n\r"
	}

	return res
}

func getContractPrice(proid string, status string) string {
	res := ""
	if status == "complete" {
		type AddressRes struct {
			Proid  string
			Amount string
		}
		var addrres AddressRes

		DB.Raw(`SELECT
					proauction2_contracts.proid,
					proauction2_values.amount
				FROM proauction2_contracts
				LEFT JOIN proauction2_values
					ON proauction2_values.proid = proauction2_contracts.value
				WHERE proauction2_contracts.owner_id = ? AND proauction2_contracts.status = "active"`, proid).Scan(&addrres)

		fmt.Println(addrres.Amount)
		res = addrres.Amount
	}

	return res
}

func getStatusString(status string) string {
	var list = map[string]string{
		"active.enquiry":           "Період уточнень",
		"pending.activation":       "Прийняття заяв на участь",
		"active.tendering":         "Прийняття заяв на участь",
		"active.auction":           "Аукціон",
		"auction":                  "Аукціон",
		"active.qualification":     "Очікується опублікування протоколу",
		"active.awarded":           "Очікується підписання договору",
		"active.rectification":     "Період редагування",
		"pending.verification":     "Очікування перевірки",
		"unsuccessful":             "Аукціон не відбувся",
		"complete":                 "Аукціон завершено. Договір підписано.",
		"cancelled":                "Аукціон відмінено",
		"draft":                    "Чернетка",
		"active.auction.dutch":     "Період аукціону (Етап автоматичного покрокового зниження початкової ціни лоту)",
		"active.auction.sealedbid": "Період аукціону (Етап подання закритих цінових пропозицій)",
		"active.auction.bestbid":   "Період аукціону (Цінова пропозиція)",
	}

	return list[status]
}

func getAuctionNumber(dgfID string, lotIdentifier string) string {
	res := ""
	if dgfID != "" {
		res = dgfID
	}
	if lotIdentifier != "" {
		res = lotIdentifier
	}
	return res
}

func getUrl(proid string) string {
	return Domen + "/prozorrosale2/auctions/" + proid
}

func init() {
	connectstr := fmt.Sprintf(
		"%s:%s@tcp(localhost:3336)/%s?charset=utf8&parseTime=True&loc=Local",
		Dbuser,
		Dbpass,
		Dbname,
	)
	db, err := gorm.Open("mysql", connectstr)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("connected...")
	}
	DB = db
}

func main() {
	var result []Result

	file, err := os.Create("result.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = DB.Raw(`SELECT
					proauction2.id,
					proauction2.proid,
					proauction2.title,
					proauction2.status,
					proauction2.auctionID,
					proauction2.dgfID,
					proauction2.lotIdentifier,
					proauction2.procurementMethodType,
					proauction2_values.amount,
					proauction2_periods.startDate
				FROM proauction2
				LEFT JOIN proauction2_values
					ON proauction2_values.proid = proauction2.value
				LEFT JOIN proauction2_periods
					ON proauction2_periods.proid = proauction2.auctionPeriod
				LEFT JOIN proauction2_orgs
					ON proauction2.procuringEntity = proauction2_orgs.proid
				LEFT JOIN proauction2_identifiers
					ON proauction2_orgs.identifier = proauction2_identifiers.proid
				WHERE proauction2_identifiers.id != ? limit 50`, "21560045").Scan(&result).Error

	if err != nil {
		panic(err)
	}

	firstline := []string{
		"#",
		"ID аукціону",
		"Номер лоту",
		"Загальна назва аукціону",
		"Тип активу",
		"Адреса",
		"Призначення",
		"Виміри",
		"Стартова Ціна",
		"Ціна Договору",
		"Дата проведення",
		"Статус",
	}

	err = writer.Write(firstline)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(firstline)
	}

	for k, v := range result {
		var line []string
		line = append(line, fmt.Sprintf("%d", k))
		line = append(line, v.AuctionID)
		line = append(line, getAuctionNumber(v.DgfID, v.LotIdentifier))
		line = append(line, v.Title)
		line = append(line, v.PMT)
		line = append(line, getItemsAddresses(v.Proid))
		line = append(line, getItemsClasif(v.Proid))
		line = append(line, getItemsUnits(v.Proid))
		line = append(line, v.Amount)
		line = append(line, getContractPrice(v.Proid, v.Status))
		line = append(line, v.StartDate)
		line = append(line, getStatusString(v.Status))
		line = append(line, getUrl(v.Proid))

		err = writer.Write(line)
		if err != nil {
			panic(err)
		} else {
			fmt.Println(line)
		}
	}

	defer DB.Close()
}
