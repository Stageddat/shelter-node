package entry

import "time"

type Entry struct {
	ID               string    `json:"id"`
	UserID           string    `json:"userId"`
	EncryptedTitle   []byte    `json:"encryptedTitle"`
	TitleIV          []byte    `json:"titleIv"`
	EncryptedContent []byte    `json:"encryptedContent"`
	ContentIV        []byte    `json:"contentIv"`
	Date             string    `json:"date"`
	Time             string    `json:"time"`
	WordCount        int       `json:"wordCount"`
	CharCount        int       `json:"charCount"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type CreateEntryRequest struct {
	UserID           string `json:"userId"`
	EncryptedTitle   []byte `json:"encryptedTitle"`
	TitleIV          []byte `json:"titleIv"`
	EncryptedContent []byte `json:"encryptedContent"`
	ContentIV        []byte `json:"contentIv"`
	Date             string `json:"date"`
	Time             string `json:"time"`
	WordCount        int    `json:"wordCount"`
	CharCount        int    `json:"charCount"`
}

type UpdateEntryRequest struct {
	EncryptedTitle   []byte `json:"encryptedTitle,omitempty"`
	TitleIV          []byte `json:"titleIv,omitempty"`
	EncryptedContent []byte `json:"encryptedContent,omitempty"`
	ContentIV        []byte `json:"contentIv,omitempty"`
}
