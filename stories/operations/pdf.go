package operations

import (
	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/hoenn/go-hn/pkg/hnapi"
	"hntoebook/stories"
	"log"
	"strconv"
)

func HTMLtoPDFGenerator(story *stories.Story, storyItem *hnapi.Story, commentItem *hnapi.Comment, pdfPath string, mobiPath string) {
	// Create new PDF generator
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		log.Fatal("PDF Generator, Error generating pdf", err)
		return
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtml.OrientationLandscape)
	pdfg.Grayscale.Set(true)

	if story != nil {
		pdfg.Title.Set(story.Title)
	} else if storyItem != nil {
		pdfg.Title.Set(storyItem.Title)
	} else if commentItem != nil {
		pdfg.Title.Set("Comment by " + commentItem.By)
	}

	// Create a new input page from an URL
	var page *wkhtml.Page
	if story != nil {
		page = wkhtml.NewPage(story.URL)
	} else if storyItem != nil {
		page = wkhtml.NewPage("https://news.ycombinator.com/item?id=" + strconv.Itoa(storyItem.ID))
	} else if commentItem != nil {
		page = wkhtml.NewPage("https://news.ycombinator.com/item?id=" + strconv.Itoa(commentItem.ID))
	}

	page.UserStyleSheet.Set("./pdfstyles/hn.css")
	//page.FooterHTML.Set("./src/app/pdfs/footer.html")

	// Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Fatal("PDF Generator, Error creating pdf", err)
		return
	}

	// Write buffer contents to file on disk
	if story != nil {
		err = pdfg.WriteFile(pdfPath + strconv.Itoa(story.Id) + ".pdf")
		if err != nil {
			log.Fatal("PDF Generator, Error writing pdf", err)
			return
		}
	} else if storyItem != nil {
		err = pdfg.WriteFile(pdfPath + strconv.Itoa(storyItem.ID) + ".pdf")
		if err != nil {
			log.Fatal("PDF Generator, Error writing pdf", err)
			return
		}
	} else if commentItem != nil {
		err = pdfg.WriteFile(pdfPath + strconv.Itoa(commentItem.ID) + ".pdf")
		if err != nil {
			log.Fatal("PDF Generator, Error writing pdf", err)
			return
		}
	}

	log.Println("PDF Generator, Creating pdf file: Success")

	PDFToMobiGenerator(story, storyItem, commentItem, pdfPath, mobiPath)
}
