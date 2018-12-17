package lib

import (
	"encoding/json"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type WikiClient struct {
	Client       *http.Client
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

func GetClient() WikiClient {
	return WikiClient{Client: &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSHandshakeTimeout: 5 * time.Second,
			DisableCompression:  false,
			DisableKeepAlives:   false,
		},
	},
		RegexOutside: regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`),
		RegexInside:  regexp.MustCompile(`[\s\p{Zs}]{2,}`),
	}
}

func (client WikiClient) filterHTML(response string) string {
	domDocTest := html.NewTokenizer(strings.NewReader(response))
	result := ""
	previousStartTokenTest := domDocTest.Token()
loopDomTest:
	for {
		tt := domDocTest.Next()
		switch {
		case tt == html.ErrorToken:
			break loopDomTest // EOF
		case tt == html.StartTagToken:
			previousStartTokenTest = domDocTest.Token()
		case tt == html.TextToken:
			if previousStartTokenTest.Data == "script" {
				continue
			}
			TxtContent := html.UnescapeString(string(domDocTest.Text()))
			if len(TxtContent) > 0 {
				TxtContent = client.RegexOutside.ReplaceAllString(TxtContent, "")
				TxtContent = client.RegexInside.ReplaceAllString(TxtContent, "")
				result += TxtContent
			}
		}
	}
	return result
}

func (client WikiClient) makeRequestSearchBasic(url string, term string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	query := req.URL.Query()
	query.Add("format", "json")
	query.Add("formatversion", "2")
	query.Add("action", "query")
	query.Add("prop", "description|pageimages")
	query.Add("piprop", "thumbnail")
	query.Add("pilicense", "any")
	query.Add("list", "search")
	query.Add("srwhat", "text")
	query.Add("srinfo", "suggestion")
	query.Add("srlimit", "1")
	query.Add("pithumbsize", "320")
	query.Add("srsearch", term)
	req.URL.RawQuery = query.Encode()
	req.Header.Add("User-Agent", "WikipediaApp/2.7.269-r-2018-12-11 (Android 8.0.0; Phone) Google Play")
	req.Header.Add("Cache-Control", "max-stale=0")
	req.Header.Add("Accept-Language", "en")
	return req
}

func (client WikiClient) makeRequestSearchAdvanced(url string, term string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	query := req.URL.Query()
	query.Add("format", "json")
	query.Add("formatversion", "2")
	query.Add("action", "query")
	query.Add("prop", "description|pageimages")
	query.Add("generator", "prefixsearch")
	query.Add("gpsnamespace", "0")
	query.Add("srnamespace", "0")
	query.Add("piprop", "thumbnail")
	query.Add("pilicense", "any")
	query.Add("list", "search")
	query.Add("srwhat", "text")
	query.Add("srinfo", "suggestion")
	query.Add("srlimit", "1")
	query.Add("pithumbsize", "320")
	query.Add("gpssearch", term)
	query.Add("gpslimit", "20")
	query.Add("srsearch", term)
	req.URL.RawQuery = query.Encode()
	req.Header.Add("User-Agent", "WikipediaApp/2.7.269-r-2018-12-11 (Android 8.0.0; Phone) Google Play")
	req.Header.Add("Cache-Control", "max-stale=0")
	req.Header.Add("Accept-Language", "en")
	return req
}

func (client WikiClient) makeRequestFirst(term string) *http.Request {
	req, _ := http.NewRequest("GET", "https://en.wikipedia.org/api/rest_v1/page/mobile-sections-lead/"+term, nil)
	req.Header.Add("User-Agent", "WikipediaApp/2.7.269-r-2018-12-11 (Android 8.0.0; Phone) Google Play")
	req.Header.Add("Cache-Control", "max-stale=0")
	req.Header.Add("Accept-Language", "en")
	return req
}

func (client WikiClient) makeRequestDetailed(term string) *http.Request {
	req, _ := http.NewRequest("GET", "https://en.wikipedia.org/api/rest_v1/page/mobile-sections-remaining/"+term, nil)
	req.Header.Add("User-Agent", "WikipediaApp/2.7.269-r-2018-12-11 (Android 8.0.0; Phone) Google Play")
	req.Header.Add("Cache-Control", "max-stale=0")
	req.Header.Add("Accept-Language", "en")
	return req
}

func (client WikiClient) SearchAdvanced(term string) ([]*WikiAdvancedResult, error, error) {
	resp, err := client.Client.Do(client.makeRequestSearchAdvanced("https://en.wikipedia.org/w/api.php", term))
	if err != nil {
		log.Fatalln(err)
		return nil, err, nil
	}
	bodyBytes, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatalln(err2)
		return nil, err, err2
	}
	defer resp.Body.Close()
	result := new(WikiSearchAdvancedResponse)
	_ = json.Unmarshal(bodyBytes, result)
	if len(result.Query.Pages) > 0 {
		return result.Query.Pages, err, err2
	} else {
		return nil, err, err2
	}
}

func (client WikiClient) SearchBasic(term string) ([]*WikiBasicResult, error, error) {
	resp, err := client.Client.Do(client.makeRequestSearchBasic("https://en.wikipedia.org/w/api.php", term))
	if err != nil {
		log.Fatalln(err)
		return nil, err, nil
	}
	bodyBytes, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatalln(err2)
		return nil, err, err2
	}
	defer resp.Body.Close()
	result := new(WikiSearchBasicResponse)
	_ = json.Unmarshal(bodyBytes, result)
	if len(result.Query.Search) > 0 {
		result.Query.Search[0].Snippet = client.filterHTML(result.Query.Search[0].Snippet)
		return result.Query.Search, err, err2
	} else {
		return nil, err, err2
	}
}

func (client WikiClient) GetPage(name string) (*WikiPage, error, error) {
	resp, err := client.Client.Do(client.makeRequestFirst(name))
	if err != nil {
		log.Fatalln(err)
		return nil, err, nil
	}
	bodyBytes, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatalln(err2)
		return nil, err, err2
	}
	resp.Body.Close()
	result := new(WikiPage)
	_ = json.Unmarshal(bodyBytes, result)
	result.Sections = result.Sections[:1]
	resp2, err3 := client.Client.Do(client.makeRequestDetailed(name))
	if err3 != nil {
		log.Fatalln(err3)
		return nil, err3, nil
	}
	bodyBytes2, err4 := ioutil.ReadAll(resp2.Body)
	if err4 != nil {
		log.Fatalln(err4)
		return nil, err, err4
	}
	resp2.Body.Close()
	result2 := new(WikiDetailedResponse)
	_ = json.Unmarshal(bodyBytes2, result2)
	for _, element := range result2.Sections {
		result.Sections = append(result.Sections, element)
	}
	for index, section := range result.Sections {
		section.Text = client.filterHTML(section.Text)
		result.Sections[index] = section
	}
	for index, note := range result.HatNotes {
		result.HatNotes[index] = client.filterHTML(note)
	}
	return result, err, err2
}
