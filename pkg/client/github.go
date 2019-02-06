package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type GithubClient struct {
	httpClient *http.Client
	baseURL    string
	baseAPIURL string
}

func New() *GithubClient {
	return NewWithParams(baseURL, baseAPIURL)
}

func NewWithParams(baseURL, baseAPIURL string) *GithubClient {
	return &GithubClient{
		httpClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: 30 * time.Second,
			},
		},
		baseURL:    baseURL,
		baseAPIURL: baseAPIURL,
	}
}

func (gc *GithubClient) Search(query string) (*SearchResult, error) {
	query = generateQuery(query)

	resp, err := gc.httpClient.Get(fmt.Sprintf("%s/search/repositories?%s", gc.baseAPIURL, query))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	sr := &SearchResult{}
	err = json.Unmarshal(body, sr)
	if err != nil {
		return nil, err
	}

	return sr, nil
}

func generateQuery(query string) string {
	searchQuery := &SearchQuery{
		Query:   query,
		Sort:    "stars",
		Order:   "desc",
		PerPage: 10,
	}
	return searchQuery.Encode()
}

func (gc *GithubClient) GetRepo(fullname string) (*Repo, error) {
	resp, err := http.Get(fmt.Sprintf("%s/repos/%s", gc.baseAPIURL, fullname))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	repo := &Repo{}
	err = json.Unmarshal(body, repo)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (gc *GithubClient) Trending(since string) (*SearchResult, error) {
	resp, err := gc.httpClient.Get(fmt.Sprintf("%s/trending?since=%s", gc.baseURL, since))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseTrendingPage(doc)
}

func parseTrendingPage(doc *goquery.Document) (*SearchResult, error) {
	repos := []Repo{}
	total := 0
	doc.Find(".repo-list li").Each(func(i int, s *goquery.Selection) {
		repo := Repo{}

		link, exists := s.Find("div.d-inline-block.col-9.mb-1 > h3 > a").Attr("href")
		if !exists {
			fmt.Println("Repo link is not found")
			return
		}
		repo.Name = strings.Split(link, "/")[2]
		repo.URL = baseURL + link

		repo.Description = parseDescription(s.Find("div.py-1 > p").Text())

		starsStr := s.Find("a[href$=\"stargazers\"]").Text()
		stars, err := parseStars(starsStr)
		if err != nil {
			fmt.Printf("Can't parse stars number: %s\n", starsStr)
			return
		}

		repo.Stars = stars

		repos = append(repos, repo)
		total++
	})

	sr := &SearchResult{
		TotalCount: total,
		Items:      sortAndLimirSR(repos, 15),
	}

	return sr, nil
}

func sortAndLimirSR(repos []Repo, limit int) []Repo {
	sort.Sort(ByStars(repos))

	if len(repos) > limit {
		repos = repos[:limit]
	}

	return repos
}

type ByStars []Repo

func (bs ByStars) Len() int           { return len(bs) }
func (bs ByStars) Swap(i, j int)      { bs[i], bs[j] = bs[j], bs[i] }
func (bs ByStars) Less(i, j int) bool { return bs[i].Stars > bs[j].Stars }

func parseDescription(desc string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9 ]+")
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(reg.ReplaceAllString(desc, ""))
}

func parseStars(str string) (int, error) {
	trimmedString := strings.TrimSpace(str)
	fixDelimiterString := strings.Replace(trimmedString, ",", "", -1)
	val, err := strconv.ParseInt(fixDelimiterString, 10, 32)
	return int(val), err
}
