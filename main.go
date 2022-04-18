package main

import (
	"bufio"
	"context"
	"fmt"
	"hntoebook/stories/operations"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/dgraph-io/badger/v3"
	"github.com/hoenn/go-hn/pkg/hnapi"
)

func startClassifierServer(ctx context.Context) {
	cmd := exec.CommandContext(ctx, "uvicorn", "main:app")
	err := cmd.Start()
	if err != nil {
		// Run could also return this error and push the program
		// termination decision to the `main` method.
		log.Fatalln("Starting classifier server, Error when starting the server. Check if all the requirements are fulfilled", err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Println("waiting on cmd:", err)
	}
}

func OperationsMode(db *badger.DB, mode string) {
	var mobiPath string

	fmt.Println("Entering operations mode")

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("mobiPath"))

		if err != nil {
			log.Fatalln("Error accessing mobiPath, Did you set the config using -c?ing db for mobiPath, Did you set the config using -c?", err)
		} else {

			err := item.Value(func(val []byte) error {
				// This func with val would only be called if item.Value encounters no error.

				// Accessing val here is valid.
				fmt.Printf("The .mobi path is: %s\n", val)

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
			case "filter":

				var wg sync.WaitGroup
				ctx, cancel := context.WithCancel(context.Background())

				// Increment the WaitGroup synchronously in the main method, to avoid
				// racing with the goroutine starting.
				wg.Add(1)
				go func() {
					startClassifierServer(ctx)
					// Signal the goroutine has completed
					wg.Done()
				}()

				fmt.Println("Enter the categories for filtering separted by a comma, e.g. Tech,Climate,Gaming:")
				var categories []string
				var categoryParam string

				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				if scanner.Err() != nil {
					log.Fatalln("Error in getting the categories")
				} else {
					categoryParam = scanner.Text()
				}

				if len(categoryParam) == 0 {
					log.Fatalln("No categories were entered, try again")
				} else {
					re := regexp.MustCompile("^(\\w+ *\\w*)+( *, *\\w* *\\w*)*$")
					if !(re.MatchString(categoryParam)) {
						log.Fatalln("Invalid category entered, Enter categories separated by a comma, e.g. Tech,Climate,Gaming:")
					}
				}

				categoriesTemp := strings.Split(categoryParam, ",")

				for _, category := range categoriesTemp {
					if len(category) > 1000 {
						fmt.Println("Was the category name in german?")
					}
					categories = append(categories, strings.TrimSpace(category))
				}

				log.Println("Categories: ", categories)

				log.Println("Creating temporary directory for storing .pdf files")
				dir, err := ioutil.TempDir("", "hn")
				if err != nil {
					log.Fatal(err)
				}

				log.Println("Temporary directory name:", dir)

				operations.UpdateStories(db, dir+"/", mobiPath, categories)
				os.Remove(dir)

				log.Println("closing via ctx")
				cancel()

				// Wait for the child goroutine to finish, which will only occur when
				// the child process has stopped and the call to cmd.Wait has returned.
				// This prevents main() exiting prematurely.
				wg.Wait()

				break

			default:
				log.Println("Creating temporary directory for storing .pdf files")
				dir, err := ioutil.TempDir("", "hn")
				if err != nil {
					log.Fatal(err)
				}

				log.Println("Temporary directory name:", dir)

				operations.UpdateStories(db, dir+"/", mobiPath, nil)
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

	log.Println("***HN to E-book***")

	db, err := badger.Open(badger.DefaultOptions("db"))
	if err != nil {
		log.Fatal(err)
	}

	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "-c":
			log.Println("Entering config mode")

			log.Println("Would like to set up the path for .mobi? Y/N")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			if scanner.Err() != nil {
				log.Fatalln("Error in getting the answer for storing path for the ebook")
			} else {
				mobiAnswer := scanner.Text()
				if mobiAnswer == "Y" {
					fmt.Println("Enter a path for storing .mobi files on the e-reader(After creating the folder) e.g. /run/media/username/Kindle/documents/Downloads/hn/ :")
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
					log.Println("Stored mobiPath for future operations")
				} else if mobiAnswer == "N" {
					log.Println("Not setting up path for .mobi this time")
				} else {
					log.Println("Invalid answer, Enter Y (or) N")
				}
			}

			fmt.Println("Would you like to setup category filter? Y/N")
			scanner.Scan()
			if scanner.Err() != nil {
				log.Fatalln("Error in getting the answer for category filter")
			} else {
				categoryFilter := scanner.Text()

				if categoryFilter == "Y" {
					log.Println("Downloading model for classification")
					log.Println("This would take a while....")

					var out []byte
					out, err = exec.Command("git", "clone", "https://huggingface.co/typeform/distilbert-base-uncased-mnli", "models/distilbert-base-uncased-mnli/").CombinedOutput()

					// if there is an error with our execution
					// handle it here
					if err != nil {
						log.Println("Downloading models, Error executing command to download models. Check if the models folder is empty", err)
						return
					}
					log.Println("Command Successfully Executed")
					output := string(out[:])
					log.Println(output)

					log.Println("Installing necessary python packages")

					out, err = exec.Command("pip", "install", "-r", "requirements.txt").CombinedOutput()

					// if there is an error with our execution
					// handle it here
					if err != nil {
						log.Println("Installing python packages, Error executing command to install python packages. Install the packages manually.", err)
						return
					}
					log.Println("Command Successfully Executed")
					output = string(out)
					log.Println(output)

				} else if categoryFilter == "N" {
					log.Println("Category filter not enabled")
				} else {
					log.Println("Invalid answer, Enter Y (or) N")
				}
			}

			log.Println("Configuration done, You can now use ./hntoebook")
			break
		case "-i":
			log.Println("Entering item mode")
			OperationsMode(db, "item")
		case "-f":
			log.Println("Entering filter mode")
			OperationsMode(db, "filter")
		default:
			log.Fatalln("Invalid argument, Use -c for config mode (or) -i for item mode")
		}
	} else {
		OperationsMode(db, "default")
	}

	defer db.Close()
}
