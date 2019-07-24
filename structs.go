package gopedia

import (
	gokhttp "github.com/BRUHItsABunny/gOkHttp"
	"regexp"
)

type WikiClient struct {
	Client       *gokhttp.HttpClient
	BaseAPIURL   string
	BasePageURL  string
	RegexOutside *regexp.Regexp
	RegexInside  *regexp.Regexp
}

type WikiAdvancedResult struct {
	PageId      int    `json:"pageid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ThumbNail   struct {
		Source string `json:"source"`
	} `json:"thumbnail"`
}

type WikiBasicResult struct {
	PageId    int    `json:"pageid"`
	Title     string `json:"title"`
	Snippet   string `json:"snippet"`
	TimeStamp string `json:"timestamp"`
	WordCount int    `json:"wordcount"`
}

type WikiSection struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
	Line string `json:"line"`
}

type WikiPage struct {
	Id           int    `json:"id"`
	Revision     int    `json:"revision"`
	LastModified string `json:"lastmodified"`
	LastModifier struct {
		User   string `json:"user"`
		Gender string `json:"gender"`
	} `json:"lastmodifier"`
	DisplayTitle    string        `json:"displaytitle"`
	NormalizedTitle string        `json:"normalizedtitle"`
	Description     string        `json:"descriptions"`
	HatNotes        []string      `json:"hatnotes"`
	Sections        []WikiSection `json:"sections"`
}

type WikiDetailedResponse struct {
	Sections []WikiSection
}

type WikiSearchBasicResponse struct {
	Query struct {
		Search []*WikiBasicResult `json:"search"`
	} `json:"query"`
}

type WikiSearchAdvancedResponse struct {
	Query struct {
		Pages []*WikiAdvancedResult `json:"pages"`
	} `json:"query"`
}
