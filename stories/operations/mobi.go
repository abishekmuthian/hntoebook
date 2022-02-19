package operations

import (
	"context"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/hoenn/go-hn/pkg/hnapi"
	"hntoebook/stories"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func PDFToMobiGenerator(db *badger.DB, story *stories.Story, storyItem *hnapi.Story, commentItem *hnapi.Comment, pdfPath string, mobiPath string) {

	var out []byte
	var err error

	// Create a new context and add a timeout to it
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	if story != nil {
		out, err = exec.CommandContext(ctx, "ebook-convert", pdfPath+strconv.Itoa(story.Id)+".pdf", mobiPath+strconv.Itoa(story.Id)+".mobi", "--authors=HN to Kindle", "--remove-first-image", "--title="+strings.ReplaceAll(story.Title, `"`, `\"`)).CombinedOutput()

		// if there is an error with our execution
		// handle it here
		if err != nil {
			log.Println("Mobi, Error executing command check the mobiPath ", err)
			return
		}
	} else if storyItem != nil {
		out, err = exec.CommandContext(ctx, "ebook-convert", pdfPath+strconv.Itoa(storyItem.ID)+".pdf", mobiPath+strconv.Itoa(storyItem.ID)+".mobi", "--authors=HN to Kindle", "--remove-first-image", "--title="+strings.ReplaceAll(storyItem.Title, `"`, `\"`)).CombinedOutput()

		// if there is an error with our execution
		// handle it here
		if err != nil {
			log.Println("Mobi, Error executing command check the mobiPath ", err)
			return
		}
	} else if commentItem != nil {
		out, err = exec.CommandContext(ctx, "ebook-convert", pdfPath+strconv.Itoa(commentItem.ID)+".pdf", mobiPath+strconv.Itoa(commentItem.ID)+".mobi", "--authors=HN to Kindle", "--remove-first-image", "--title="+strings.ReplaceAll("Comment by "+commentItem.By, `"`, `\"`)).CombinedOutput()

		// if there is an error with our execution
		// handle it here
		if err != nil {
			log.Println("Mobi, Error executing command check the mobiPath ", err)
			return
		}
	}

	// We want to check the context error to see if the timeout was executed.
	// The error returned by cmd.CombinedOutput() will be OS specific based on what
	// happens when a process is killed.
	if ctx.Err() == context.DeadlineExceeded {
		log.Println("Command timed out")
		return
	}

	// as the out variable defined above is of type []byte we need to convert
	// this to a string or else we will see garbage printed out in our console
	// this is how we convert it to a string
	fmt.Println("Command Successfully Executed")
	output := string(out[:])
	fmt.Println(output)

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(strconv.Itoa(story.Id)), []byte("true"))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Stored item id for preventing duplicates")
}
