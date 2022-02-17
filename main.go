package main

import (
	"bufio"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/hoenn/go-hn/pkg/hnapi"
	"hntoebook/stories/operations"
	"io/ioutil"
	"log"
	"os"
)

func OperationsMode(db *badger.DB, mode string) {
	var mobiPath string

	fmt.Println("Entering operations mode")

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("mobiPath"))

		if err != nil {
			log.Fatalln("Error accessing db for mobiPath, Did you set the config using -c?", err)
		} else {

			err := item.Value(func(val []byte) error {
				// This func with val would only be called if item.Value encounters no error.

				// Accessing val here is valid.
				fmt.Printf("The answer is: %s\n", val)

				// Copying or parsing val is valid.
				mobiPath = string(append([]byte{}, val...))

				return nil
			})

			if err != nil {
				log.Fatalln("Item not found in the database", err)
			}

			switch mode {
			case "item":
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

				fmt.Println("Temporary directory name:", dir)

				switch item.(type) {
				case *hnapi.Story:
					fmt.Println("Found HN story")
					storyItem := item.(*hnapi.Story)
					operations.HTMLtoPDFGenerator(db, nil, storyItem, nil, dir+"/", string(mobiPath))
					break
				case *hnapi.Comment:
					fmt.Println("Found HN comment")
					commentItem := item.(*hnapi.Comment)
					operations.HTMLtoPDFGenerator(db, nil, nil, commentItem, dir+"/", string(mobiPath))
					break
				}
				os.Remove(dir)
				break
			default:
				fmt.Println("Creating temporary directory for storing .pdf files")
				dir, err := ioutil.TempDir("", "hn")
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("Temporary directory name:", dir)

				operations.UpdateStories(db, dir+"/", mobiPath)
				os.Remove(dir)
				break
			}

		}
		return nil
	})

	if err != nil {
		log.Fatalln("Db couldn't be viewed", err)
	}

}

func main() {
	var mobiPath string

	fmt.Println("***HN to E-book***")

	db, err := badger.Open(badger.DefaultOptions("db"))
	if err != nil {
		log.Fatal(err)
	}

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
			err = db.Update(func(txn *badger.Txn) error {
				err := txn.Set([]byte("mobiPath"), []byte(mobiPath))
				return err
			})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Stored mobiPath for future operations")

			OperationsMode(db, "default")
			break
		case "-i":
			fmt.Println("Entering item mode")
			OperationsMode(db, "item")
		default:
			log.Fatalln("Invalid argument, Use -c for config mode (or) -i for item mode")
		}
	} else {
		OperationsMode(db, "default")
	}

	db.Close()
}
