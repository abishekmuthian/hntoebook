package main

import (
	"bufio"
	"fmt"
	"git.mills.io/prologic/bitcask"
	"github.com/hoenn/go-hn/pkg/hnapi"
	"hntoebook/stories/operations"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var mobiPath string

	fmt.Println("***HN to E-book***")

	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "-c":
			fmt.Println("Entering config mode")

			fmt.Println("Enter a path for storing .mobi files on the e-reader e.g. /run/media/username/Kindle/documents/Downloads/hn/ :")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			if scanner.Err() != nil {
				log.Fatalln("Error in getting the path for storing the ebook")
			} else {
				mobiPath = scanner.Text()
			}

			db, err := bitcask.Open("db")
			if err != nil {
				fmt.Println("Error opening db ", err)
			}
			db.Put([]byte("mobiPath"), []byte(mobiPath))
			db.Close()
			fmt.Println("Stored mobiPath for future operations")

			OperationsMode(mobiPath)
			break
		case "-i":
			fmt.Println("Entering item mode")

			db, err := bitcask.Open("db")
			if err != nil {
				fmt.Println("Error opening db ", err)
			}
			mobiPath, err := db.Get([]byte("mobiPath"))
			db.Close()
			if err != nil {
				log.Fatalln("Error accessing db for mobiPath, Did you set the config using -c?")
			} else {
				fmt.Println("Enter the HN story or comment item id:")

				var itemId string

				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				if scanner.Err() != nil {
					log.Fatalln("Error in getting the HN item id")
				} else {
					itemId = scanner.Text()
				}

				c := hnapi.NewHNClient()
				// Get the details of the current max item.
				item, err := c.Item(itemId)

				if err != nil {
					log.Fatalln("Item mode, Error fetching story item")
				}

				fmt.Println("Creating temporary directory for storing .pdf files")
				dir, err := ioutil.TempDir("", "hn")
				if err != nil {
					log.Fatal(err)
				}
				defer os.Remove(dir)

				fmt.Println("Temporary directory name:", dir)

				switch item.(type) {
				case *hnapi.Story:
					fmt.Println("Found HN story")
					storyItem := item.(*hnapi.Story)
					operations.HTMLtoPDFGenerator(nil, storyItem, nil, dir+"/", string(mobiPath))

				case *hnapi.Comment:
					fmt.Println("Found HN comment")
					commentItem := item.(*hnapi.Comment)
					operations.HTMLtoPDFGenerator(nil, nil, commentItem, dir+"/", string(mobiPath))

				}
			}
		default:
			log.Fatalln("Invalid argument, Use -c for config mode (or) -i for item mode")
		}
	} else {
		db, _ := bitcask.Open("db")
		if db.Has([]byte("mobiPath")) {
			mobiPath, err := db.Get([]byte("mobiPath"))
			if err != nil {
				db.Close()
				log.Fatalln("Error accessing db for mobiPath, Did you set the config using -c?")
			} else {
				OperationsMode(string(mobiPath))
			}
		} else {
			db.Close()
			log.Fatalln("Error accessing db for mobiPath, Did you set the config using -c?")
		}

	}
}

func OperationsMode(mobiPath string) {
	fmt.Println("Entering operations mode")

	fmt.Println("Creating temporary directory for storing .pdf files")
	dir, err := ioutil.TempDir("", "hn")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(dir)

	fmt.Println("Temporary directory name:", dir)

	operations.UpdateStories(dir+"/", mobiPath)
}
