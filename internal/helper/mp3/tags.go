package MP3_HELPER

import (
	"fmt"
	"os"

	"github.com/bogem/id3v2/v2"
)

type Tags struct {
	FilePath string
	Cover    string

	Artist string
	Title  string
}

func SetTags(tags Tags) error {
	tag, err := id3v2.Open(tags.FilePath, id3v2.Options{Parse: true})
	if err != nil {
		return fmt.Errorf("failed open mp3 file: %s", err)
	}
	defer tag.Close()

	if len(tags.Artist) != 0 && len(tags.Title) != 0 {
		tag.SetArtist(tags.Artist)
		tag.SetTitle(tags.Title)
	}

	if art, err := os.ReadFile(tags.Cover); err == nil {
		pic := id3v2.PictureFrame{
			Encoding:    id3v2.EncodingUTF8,
			MimeType:    "image/png",
			PictureType: id3v2.PTFrontCover,
			Description: "Front cover",
			Picture:     art,
		}
		tag.AddAttachedPicture(pic)
	} else {
		return fmt.Errorf("failed add cover to file: %s", err)
	}

	if err = tag.Save(); err != nil {
		return fmt.Errorf("Error while saving a tag: %s", err)
	}

	return nil
}
