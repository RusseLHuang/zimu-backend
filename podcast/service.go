package podcast

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/RusseLHuang/zimu-backend/utils"

	"github.com/mmcdole/gofeed"
)

func GetCollection(itunesID string) Collection {
	const itunesLookupURI string = "https://itunes.apple.com/lookup?entity=podcast&limit=1&id="
	itunesLookupAPI := itunesLookupURI + url.QueryEscape(itunesID)

	val := utils.RedisGet(itunesID)
	if val != "" {
		var collection Collection
		searchByte := []byte(val)
		json.Unmarshal(searchByte, &collection)
		return collection
	}

	var collectionResp CollectionResp

	val = fetch(itunesLookupAPI)
	searchByte := []byte(val)
	json.Unmarshal(searchByte, &collectionResp)
	results := collectionResp.Results

	if len(results) == 0 {
		var result Collection
		utils.RedisSet(itunesID, "")
		return result
	}

	collection := results[0]
	podcastList := parseRSSFeed(collection.FeedURL)
	collection.ContentFeed = podcastList

	collectionStr, err := json.Marshal(collection)
	if err != nil {
		panic(err)
	}
	utils.RedisSet(itunesID, collectionStr)

	return collection
}

func GetAll() []Result {
	const itunesFeedAPI string = "https://rss.itunes.apple.com/api/v1/us/podcasts/top-podcasts/all/100/explicit.json"
	const keyword string = "feeds"

	var results []Result
	var feedResp FeedResp

	val := utils.RedisGet(keyword)
	if val == "" {
		val = fetch(itunesFeedAPI)
		utils.RedisSet(keyword, val)
	}
	feedByte := []byte(val)
	json.Unmarshal(feedByte, &feedResp)
	results = feedResp.Feed.Results

	return results
}

func Search(keyword string) []Collection {
	const itunesSearchURI string = "https://itunes.apple.com/search?entity=podcast&limit=25&term="
	itunesSearchAPI := itunesSearchURI + url.QueryEscape(keyword)

	var results []Collection
	var collectionResp CollectionResp

	val := utils.RedisGet(keyword)
	if val == "" {
		val = fetch(itunesSearchAPI)
		utils.RedisSet(keyword, val)
	}

	searchByte := []byte(val)
	json.Unmarshal(searchByte, &collectionResp)
	results = collectionResp.Results

	return results
}

func fetch(uri string) string {
	var netClient = &http.Client{
		Timeout: time.Second * 60,
	}

	resp, err := netClient.Get(uri)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	val := string(body)
	return val
}

func parseRSSFeed(uri string) []ContentFeed {
	contentFeeds := []ContentFeed{}
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(uri)
	if err != nil {
		panic(err)
	}

	items := feed.Items
	for i := 0; i < len(items); i++ {
		for j := 0; j < len(items[i].Enclosures); j++ {
			contentFeed := ContentFeed{
				Title:         items[i].Title,
				Desc:          items[i].ITunesExt.Summary,
				PublishedDate: items[i].Published,
				ContentURL:    items[i].Enclosures[j].URL,
			}

			contentFeeds = append(contentFeeds, contentFeed)
		}
	}

	return contentFeeds
}
