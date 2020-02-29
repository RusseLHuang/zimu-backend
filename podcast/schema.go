package podcast

import (
	"strconv"

	"github.com/graphql-go/graphql"
)

type FeedResp struct {
	Feed Feed `json:"feed"`
}

type Feed struct {
	Title     string   `json:"title"`
	ID        string   `json:"id"`
	Copyright string   `json:"copyright"`
	Country   string   `json:"country"`
	Icon      string   `json:"icon"`
	Updated   string   `json:"updated"`
	Results   []Result `json:"results"`
}

type Result struct {
	ArtistName            string  `json:"artistName"`
	ID                    string  `json:"id"`
	ReleaseDate           string  `json:"releaseDate"`
	Name                  string  `json:"name"`
	Kind                  string  `json:"kind"`
	Copyright             string  `json:"copyright"`
	ArtistID              string  `json:"artistId"`
	ContentAdvisoryRating string  `json:"contentAdvisoryRating"`
	ArtistURL             string  `json:"artistUrl"`
	ArtworkURL100         string  `json:"artworkUrl100"`
	Genres                []Genre `json:"genres"`
	URL                   string  `json:"url"`
}

type CollectionResp struct {
	ResultCount int          `json:"resultCount"`
	Results     []Collection `json:"results"`
}

type Collection struct {
	CollectionID   int           `json:"collectionId"`
	ArtistID       int           `json:"artistId"`
	ArtistName     string        `json:"artistName"`
	CollectionName string        `json:"collectionName"`
	FeedURL        string        `json:"feedUrl"`
	ContentFeed    []ContentFeed `json:"contentFeed"`
	ArtworkURL100  string        `json:"artworkUrl100"`
	ArtworkURL600  string        `json:"artworkUrl600"`
	ReleaseDate    string        `json:"releaseDate"`
	Country        string        `json:"country"`
	GenreIds       string        `json:"genreIds"`
	Genres         string        `json:"genres"`
}

type ContentFeed struct {
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	PublishedDate string `json:"publishedDate"`
	ContentURL    string `json:"contentUrl"`
}

type Genre struct {
	GenreID string `json:"genreId"`
	Name    string `json:"name"`
	URL     string `json:"url"`
}

// Podcast services
type Podcast struct {
	ArtistName            string  `json:"artistName"`
	ID                    string  `json:"id"`
	ReleaseDate           string  `json:"releaseDate"`
	Name                  string  `json:"name"`
	Kind                  string  `json:"kind"`
	Copyright             string  `json:"copyright"`
	ArtistID              string  `json:"artistId"`
	ContentAdvisoryRating string  `json:"contentAdvisoryRating"`
	ArtistURL             string  `json:"artistUrl"`
	ArtworkURL100         string  `json:"artworkUrl100"`
	Genres                []Genre `json:"genres"`
	URL                   string  `json:"url"`
}

var genreType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Genre",
		Fields: graphql.Fields{
			"genreId": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"url": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var contentFeedType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ContentFeed",
		Fields: graphql.Fields{
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"desc": &graphql.Field{
				Type: graphql.String,
			},
			"publishedDate": &graphql.Field{
				Type: graphql.String,
			},
			"contentUrl": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var CollectionType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Collection",
		Fields: graphql.Fields{
			"collectionId": &graphql.Field{
				Type: graphql.Int,
			},
			"artistId": &graphql.Field{
				Type: graphql.Int,
			},
			"artistName": &graphql.Field{
				Type: graphql.String,
			},
			"collectionName": &graphql.Field{
				Type: graphql.String,
			},
			"contentFeed": &graphql.Field{
				Type: graphql.NewList(contentFeedType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					parent := p.Source.(Collection)
					collectionID := strconv.Itoa(parent.CollectionID)
					collection := GetCollection(collectionID)

					return collection.ContentFeed, nil
				},
			},
			"artworkUrl100": &graphql.Field{
				Type: graphql.String,
			},
			"artworkUrl600": &graphql.Field{
				Type: graphql.String,
			},
			"releaseDate": &graphql.Field{
				Type: graphql.String,
			},
			"country": &graphql.Field{
				Type: graphql.String,
			},
			"genreIds": &graphql.Field{
				Type: graphql.String,
			},
			"genres": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var PodcastType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Podcast",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"artistName": &graphql.Field{
			Type: graphql.String,
		},
		"releaseDate": &graphql.Field{
			Type: graphql.String,
		},
		"kind": &graphql.Field{
			Type: graphql.String,
		},
		"copyright": &graphql.Field{
			Type: graphql.String,
		},
		"artistId": &graphql.Field{
			Type: graphql.String,
		},
		"contentAdvisoryRating": &graphql.Field{
			Type: graphql.String,
		},
		"artistUrl": &graphql.Field{
			Type: graphql.String,
		},
		"artworkUrl100": &graphql.Field{
			Type: graphql.String,
		},
		"url": &graphql.Field{
			Type: graphql.String,
		},
		"genre": &graphql.Field{
			Type: graphql.NewList(genreType),
		},
	},
})
