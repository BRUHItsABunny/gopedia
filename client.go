package gopedia

import (
	"errors"
	gokhttp "github.com/BRUHItsABunny/gOkHttp"
	"golang.org/x/net/html"
	"net/http"
	"regexp"
	"strings"
)

func GetClient() WikiClient {
	options := gokhttp.HttpClientOptions{
		Headers: map[string]string{
			"User-Agent":      "WikipediaApp/2.7.269-r-2018-12-11 (Android 8.0.0; Phone) Google Play",
			"Cache-Control":   "max-stale=0",
			"Accept-Language": "en",
		},
	}
	client := gokhttp.GetHTTPClient(&options)
	return WikiClient{
		Client:       &client,
		BaseAPIURL:   "https://en.wikipedia.org/w/api.php",
		BasePageURL:  "https://en.wikipedia.org/api/rest_v1/page/",
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

func (client *WikiClient) SearchAdvanced(term string) ([]*WikiAdvancedResult, error) {

	var err error
	var req *http.Request
	var resp *gokhttp.HttpResponse
	var result *WikiSearchAdvancedResponse

	params := map[string]string{
		"format":        "json",
		"formatversion": "2",
		"action":        "query",
		"prop":          "description|pageimages",
		"generator":     "prefixsearch",
		"gpsnamespace":  "0",
		"srnamespace":   "0",
		"piprop":        "thumbnail",
		"pilicense":     "any",
		"list":          "search",
		"srwhat":        "text",
		"srinfo":        "suggestion",
		"srlimit":       "1",
		"pithumbsize":   "320",
		"gpssearch":     term,
		"gpslimit":      "20",
		"srsearch":      term,
	}
	req, err = client.Client.MakeGETRequest(client.BaseAPIURL, params, map[string]string{})
	if err == nil {
		resp, err = client.Client.Do(req)
		if err == nil {
			err = resp.Object(&result)
			if err == nil {
				if len(result.Query.Pages) > 0 {
					return result.Query.Pages, nil
				} else {
					err = errors.New("no result found")
				}
			}
		}
	}
	return nil, err
}

func (client *WikiClient) SearchBasic(term string) ([]*WikiBasicResult, error) {

	var err error
	var req *http.Request
	var resp *gokhttp.HttpResponse
	var result WikiSearchBasicResponse

	params := map[string]string{
		"format":        "json",
		"formatversion": "2",
		"action":        "query",
		"prop":          "description|pageimages",
		"piprop":        "thumbnail",
		"pilicense":     "any",
		"list":          "search",
		"srwhat":        "text",
		"srinfo":        "suggestion",
		"srlimit":       "1",
		"pithumbsize":   "320",
		"srsearch":      term,
	}
	req, err = client.Client.MakeGETRequest(client.BaseAPIURL, params, map[string]string{})
	if err == nil {
		resp, err = client.Client.Do(req)
		if err == nil {
			err = resp.Object(&result)
			if err == nil {
				if len(result.Query.Search) > 0 {
					for index, res := range result.Query.Search {
						result.Query.Search[index].Snippet = client.filterHTML(res.Snippet)
					}
					return result.Query.Search, nil
				} else {
					err = errors.New("no results found")
				}
			}
		}
	}
	return nil, err
}

func (client *WikiClient) getPage(name string) (*WikiPage, error) {

	var err error
	var req *http.Request
	var resp *gokhttp.HttpResponse
	var result WikiPage

	req, err = client.Client.MakeGETRequest(client.BasePageURL+"mobile-sections-remaining/"+name, map[string]string{}, map[string]string{})
	if err == nil {
		resp, err = client.Client.Do(req)
		if err == nil {
			err = resp.Object(&result)
			if err == nil {
				result.Sections = result.Sections[:1]
				return &result, nil
			}
		}
	}
	return nil, err
}

func (client *WikiClient) getPageDetails(name string, page *WikiPage) (*WikiPage, error) {

	var err error
	var req *http.Request
	var resp *gokhttp.HttpResponse
	var result WikiDetailedResponse

	req, err = client.Client.MakeGETRequest(client.BasePageURL+"mobile-sections-remaining/"+name, map[string]string{}, map[string]string{})
	if err == nil {
		resp, err = client.Client.Do(req)
		if err == nil {
			err = resp.Object(&result)
			if err == nil {
				for _, element := range result.Sections {
					page.Sections = append(page.Sections, element)
				}
				for index, section := range page.Sections {
					section.Text = client.filterHTML(section.Text)
					page.Sections[index] = section
				}
				for index, note := range page.HatNotes {
					page.HatNotes[index] = client.filterHTML(note)
				}
				return page, nil
			}
		}
	}
	return nil, err
}

func (client *WikiClient) GetPage(name string) (*WikiPage, error) {

	var err error
	var page *WikiPage

	page, err = client.getPage(name)
	if err == nil {
		page, err = client.getPageDetails(name, page)
		if err == nil {
			return page, nil
		}
	}
	return nil, err
}
