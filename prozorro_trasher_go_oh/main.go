// Package main provides ...
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	Dbuser  string = "root"
	Dbpass  string = "rootpass"
	Dbname  string = "prozorrodev"
	Webhome string = "/home/www/polonex.com.ua/web"
)

type Attachment struct {
	Id     int
	Group  string
	FileID int
}

type File struct {
	Id        int
	Name      string
	Path      string
	Ext       string
	UpdatedAt int
}

type Proelements struct {
	Id     int
	ElId   string
	Parent string
}

type Proauctions struct {
	Id     string
	Status string
}

func removeOne(path string) {
	err := os.Remove(path)

	if err != nil {
		fmt.Println(err)
		//return
	} else {
		fmt.Printf("%s...DELETED\n", path)
	}
}

func ClearIllustrations() {
	atts := []Attachment{}
	controlltime := time.Now().Add(-24 * time.Hour * 356).Unix()

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

	// LIKE
	db.Where("`group` LIKE ?", "%illustration%").Find(&atts)

	for i := 0; i < len(atts); i++ {
		file := File{}
		db.Where("`id` = ?", atts[i].FileID).Find(&file)
		if file.UpdatedAt < int(controlltime) {
			name := file.Name
			//name = strings.Replace(file.Name, " ", "\\ ", -1)
			//name = strings.Replace(name, ",", "\\,", -1)
			//name = strings.Replace(name, "(", "\\(", -1)
			//name = strings.Replace(name, ")", "\\)", -1)
			path := fmt.Sprintf("%s%s%s.%s", Webhome, file.Path, name, file.Ext)
			path_b := fmt.Sprintf("%s%s%s_big_.%s", Webhome, file.Path, name, file.Ext)
			path_m := fmt.Sprintf("%s%s%s_mid_.%s", Webhome, file.Path, name, file.Ext)
			path_t := fmt.Sprintf("%s%s%s_tumb_.%s", Webhome, file.Path, name, file.Ext)
			removeOne(path)
			removeOne(path_b)
			removeOne(path_m)
			removeOne(path_t)
		}
	}

	fmt.Println("Illustrations clear DONE!!!")
	fmt.Println(int(controlltime))
}

func ClearOldBidFiles() {
	proelems := []Proelements{}

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

	// LIKE
	db.Where("`el_type` = ?", "3").Find(&proelems)

	for i := 0; i < len(proelems); i++ {
		auct := Proauctions{}

		db.Where("`id` = ?", proelems[i].Parent).Find(&auct)

		if auct.Status == "unsuccessful" || auct.Status == "complete" {
			fmt.Println(proelems[i])
			fmt.Println(auct)

			atts := []Attachment{}

			db.Where("`group` LIKE ?", "%_bid_files_"+auct.Id+"%").Find(&atts)

			fmt.Println(len(atts))

			for i := 0; i < len(atts); i++ {
				file := File{}
				db.Where("`id` = ?", atts[i].FileID).Find(&file)
				name := file.Name
				//name = strings.Replace(file.Name, " ", "\\ ", -1)
				//name = strings.Replace(name, ",", "\\,", -1)
				//name = strings.Replace(name, "(", "\\(", -1)
				//name = strings.Replace(name, ")", "\\)", -1)
				path := fmt.Sprintf("%s%s%s.%s", Webhome, file.Path, name, file.Ext)
				path_b := fmt.Sprintf("%s%s%s_big_.%s", Webhome, file.Path, name, file.Ext)
				path_m := fmt.Sprintf("%s%s%s_mid_.%s", Webhome, file.Path, name, file.Ext)
				path_t := fmt.Sprintf("%s%s%s_tumb_.%s", Webhome, file.Path, name, file.Ext)
				removeOne(path)
				removeOne(path_b)
				removeOne(path_m)
				removeOne(path_t)
			}
		}
	}

	fmt.Println("Bids documents clear DONE!!!")
}

func main() {
	//ClearIllustrations()
	ClearOldBidFiles()
}
