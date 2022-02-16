package operations

import (
	"git.mills.io/prologic/bitcask"
	"github.com/hoenn/go-hn/pkg/hnapi"
	Story "hntoebook/stories"
	"log"
	"strconv"
	"time"
)

func UpdateStories(pdfPath string, mobiPath string) {
	c := hnapi.NewHNClient()

	//Get top stories
	topStoryIDs, err := c.TopStoryIDs(hnapi.Best)

	if err != nil {
		log.Fatal("Update: Error fetching top stories", err)
	}

	var stories []Story.Story

	for _, topStoryID := range topStoryIDs {

		storyItem, err := c.Item(strconv.Itoa(topStoryID))

		if err != nil {
			log.Fatal("Update: Error fetching story item", err)
		}

		topStory := storyItem.(*hnapi.Story)

		db, err := bitcask.Open("db")
		if err != nil {
			log.Fatal("Error opening the db", err)
		}

		if !db.Has([]byte(strconv.Itoa(topStory.ID))) {
			db.Close()
			log.Println("Story was previously processed")
			if topStory.Descendants > 0 {
				topCommentID := topStory.Kids[0]

				commentItem, err := c.Item(strconv.Itoa(topCommentID))

				if err != nil {
					log.Fatal("Update: Error fetching comment item", err)
				}

				topComment := commentItem.(*hnapi.Comment)

				log.Println("Update, Top story ID", topStory.ID)
				log.Println("Update, Top comment ID", topComment.ID)

				if time.Now().Sub(time.Unix(topStory.Time, 0)).Hours() > 9 && time.Now().Sub(time.Unix(topStory.Time, 0)).Hours() < 24 && topStory.Descendants > 20 && time.Now().Sub(time.Unix(topComment.Time, 0)).Hours() > 2 {
					log.Println("Time difference of the story: ", time.Now().Sub(time.Unix(topStory.Time, 0)).Hours())

					story := Story.Story{
						Id:    topStory.ID,
						Time:  time.Unix(topStory.Time, 0).UTC(),
						Title: topStory.Title,
						URL:   "https://news.ycombinator.com/item?id=" + strconv.Itoa(topStory.ID),
					}

					stories = append(stories, story)

					log.Println("Created story, Story ID:", story.Id)
					HTMLtoPDFGenerator(&story, nil, nil, pdfPath, mobiPath)

				} else {
					log.Println("Story Update, Top story comment threshold not met or Story not older than 9 hours or Story older than 24 hours")
				}
			} else {
				log.Println("Story Update, No descendants to the top story")
			}
		} else {
			db.Close()
			log.Println("Story already placed in db")
		}
	}
}
