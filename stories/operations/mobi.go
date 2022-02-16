package operations

import (
	"fmt"
	"git.mills.io/prologic/bitcask"
	"github.com/hoenn/go-hn/pkg/hnapi"
	"hntoebook/stories"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func PDFToMobiGenerator(story *stories.Story, storyItem *hnapi.Story, commentItem *hnapi.Comment, pdfPath string, mobiPath string) {

	var out []byte
	var err error
	if story != nil {
		out, err = exec.Command("ebook-convert", pdfPath+strconv.Itoa(story.Id)+".pdf", mobiPath+strconv.Itoa(story.Id)+".mobi", "--authors=HN to Kindle", "--remove-first-image", "--title="+strings.ReplaceAll(story.Title, `"`, `\"`)).Output()

		// if there is an error with our execution
		// handle it here
		if err != nil {
			log.Fatal("Mobi, Error executing command check the mobiPath ", err)
			return
		}
	} else if storyItem != nil {
		out, err = exec.Command("ebook-convert", pdfPath+strconv.Itoa(storyItem.ID)+".pdf", mobiPath+strconv.Itoa(storyItem.ID)+".mobi", "--authors=HN to Kindle", "--remove-first-image", "--title="+strings.ReplaceAll(storyItem.Title, `"`, `\"`)).Output()

		// if there is an error with our execution
		// handle it here
		if err != nil {
			log.Fatal("Mobi, Error executing command check the mobiPath ", err)
			return
		}
	} else if commentItem != nil {
		out, err = exec.Command("ebook-convert", pdfPath+strconv.Itoa(commentItem.ID)+".pdf", mobiPath+strconv.Itoa(commentItem.ID)+".mobi", "--authors=HN to Kindle", "--remove-first-image", "--title="+strings.ReplaceAll("Comment by "+commentItem.By, `"`, `\"`)).Output()

		// if there is an error with our execution
		// handle it here
		if err != nil {
			log.Fatal("Mobi, Error executing command check the mobiPath ", err)
			return
		}
	}

	// as the out variable defined above is of type []byte we need to convert
	// this to a string or else we will see garbage printed out in our console
	// this is how we convert it to a string
	fmt.Println("Command Successfully Executed")
	output := string(out[:])
	fmt.Println(output)

	db, err := bitcask.Open("db")
	if err != nil {
		log.Fatal("Error opening the db", err)
	}
	err = db.Put([]byte(strconv.Itoa(story.Id)), []byte("True"))
	db.Close()
	if err != nil {
		log.Fatal("Error storing the item id in the db")
	}
	fmt.Println("Stored item id for preventing duplicates")
}
