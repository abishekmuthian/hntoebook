package operations

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hntoebook/stories"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/hoenn/go-hn/pkg/hnapi"
)

func categoryFilter(story *stories.Story, categories []string) bool {
	if len(categories) < 1 {
		log.Fatalln("Update, User has no categories")
	}

	type bertRequest struct {
		Text   string   `json:"text"`
		Labels []string `json:"labels"`
	}

	req := &bertRequest{
		Text:   story.Title,
		Labels: categories,
	}

	postBody, err := json.Marshal(req)
	if err != nil {
		log.Fatalln("Update stories. Error creating POST body for bert")
	}

	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post("http://127.0.0.1:8000/classification", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalln("Bert, An Error Occurred", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Bert An Error Occurred", err)
	}

	type bertResponse struct {
		Labels []string  `json:"labels"`
		Scores []float64 `json:"scores"`
	}

	res := bertResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Fatalln("Bert, Error unmarshalling json", err)
	}

	log.Printf("Bert, Title: %s\n", story.Title)

	for i, score := range res.Scores {
		fmt.Printf("Bert, Label: %s, Score: %f\n", res.Labels[i], res.Scores[i])
		if score > 0.75 {
			fmt.Printf("Bert, Score matches threshold, Label: %s, Score: %f", res.Labels[i], res.Scores[i])
			return true
		}
	}

	return false
}

func UpdateStories(db *badger.DB, pdfPath string, mobiPath string, categories []string) {
	c := hnapi.NewHNClient()

	//Get top stories
	topStoryIDs, err := c.TopStoryIDs(hnapi.Best)

	if err != nil {
		log.Fatal("Update: Error fetching top stories", err)
	}

	for _, topStoryID := range topStoryIDs {

		storyItem, err := c.Item(strconv.Itoa(topStoryID))

		if err != nil {
			log.Fatal("Update: Error fetching story item", err)
		}

		topStory := storyItem.(*hnapi.Story)

		err = db.View(func(txn *badger.Txn) error {
			_, err := txn.Get([]byte(strconv.Itoa(topStory.ID)))
			if err == nil {
				log.Println("Story already found in db")
				return nil
			}

			if errors.Is(err, badger.ErrKeyNotFound) {
				log.Println("Story was not previously processed")
				if topStory.Descendants < 1 {
					log.Println("Story Update, No descendants to the top story")
					return nil
				}
				topCommentID := topStory.Kids[0]

				commentItem, err := c.Item(strconv.Itoa(topCommentID))

				if err != nil {
					log.Fatal("Update: Error fetching comment item", err)
				}

				topComment := commentItem.(*hnapi.Comment)

				log.Println("Update, Top story ID", topStory.ID)
				log.Println("Update, Top comment ID", topComment.ID)

				if !includeAsTopStory(topStory, topComment) {
					log.Println("Story Update, Top story comment threshold not met or Story not older than 9 hours or Story older than 24 hours")
					return nil
				}
				log.Println("Time difference of the story: ", time.Since(time.Unix(topStory.Time, 0)).Hours())

				story := stories.Story{
					Id:    topStory.ID,
					Time:  time.Unix(topStory.Time, 0).UTC(),
					Title: topStory.Title,
					URL:   "https://news.ycombinator.com/item?id=" + strconv.Itoa(topStory.ID),
				}

				log.Println("Created story, Story ID:", story.Id)

				if categories != nil {
					if categoryFilter(&story, categories) {
						HTMLtoPDFGenerator(db, &story, nil, nil, pdfPath, mobiPath)
					}
				} else {
					HTMLtoPDFGenerator(db, &story, nil, nil, pdfPath, mobiPath)
				}
				return nil
			}
			return err
		})
		if err != nil {
			log.Println("Update: Error while updating database", err)
		}
	}
}

// includeAsTopStory returns true if a story meets the following criteria to be categorized as "top":
// - Between 9 and 24 hours old.
// - Has at least 20 comments.
// - Top comment on the story is no more than 2 hours old.
func includeAsTopStory(topStory *hnapi.Story, topComment *hnapi.Comment) bool {
	return time.Since(time.Unix(topStory.Time, 0)).Hours() > 9 &&
		time.Since(time.Unix(topStory.Time, 0)).Hours() < 24 &&
		topStory.Descendants > 20 &&
		time.Since(time.Unix(topComment.Time, 0)).Hours() > 2
}
