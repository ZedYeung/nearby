package main

import (
	"encoding/json"
	"fmt"
	"os"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
	"github.com/ZedYeung/nearby/model"
	"github.com/go-redis/redis"
	"golang.org/x/net/context"
	elastic "gopkg.in/olivere/elastic.v3"
)

// http://localhost:8080/search?lat=10.0&lon=20.0
// http://localhost:8080/search?lat=10.0&lon=20.0&range=10
func Search(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

	fmt.Println("Received one request for search")

	if req.Method != "GET" {
		fmt.Println(req.Method)
		return
	}

	lat, _ := strconv.ParseFloat(req.URL.Query().Get("lat"), 64)
	lon, _ := strconv.ParseFloat(req.URL.Query().Get("lon"), 64)
	ran := os.Getenv("DISTANCE")
	if val := req.URL.Query().Get("range"); val != "" {
		ran = val + "km"
	}
	fmt.Println("range is ", ran)

	fmt.Printf("Search received: %f %f %s\n", lat, lon, ran)

	// use cache, reference: https://github.com/go-redis/redis
	// use lat+lon+range as key
	key := req.URL.Query().Get("lat") + ":" + req.URL.Query().Get("lon") + ":" + ran
	ENABLE_MEMCACHE, err := strconv.ParseBool(os.Getenv("ENABLE_MEMCACHE"))
	// First find query with cache.
	if ENABLE_MEMCACHE {
		// build connection with Redis
		rs_client := redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_URL"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})

		val, err := rs_client.Get(key).Result()
		if err != nil {
			fmt.Printf("Cache miss! Redis cannot find the key %s as %v.\n", key, err)
		} else { // cache hit
			fmt.Printf("Cache hit! Redis find the key %s.\n", key)
			res.Header().Set("Content-Type", "application/json")
			res.Write([]byte(val))
			return
		}
	}

	ctx := context.Background()

	// for here, not enable cache or cache miss, we need to search in database(ES).
	// Create a client，which means we create a connection to ES. If there is err, return.
	client, err := elastic.NewClient(elastic.SetURL(os.Getenv("ES_URL")), elastic.SetSniff(false))
	if err != nil {
		panic(err)
		return
	}

	// Define geo distance query as specified in
	// https://www.elastic.co/guide/en/elasticsearch/reference/5.2/query-dsl-geo-distance-query.html
	// Prepare a geo based query to find posts within a geo box.
	query := elastic.NewGeoDistanceQuery("location")
	query = query.Distance(ran).Lat(lat).Lon(lon)

	// Some delay may range from seconds to minutes.
	// Get the results based on Index (similar to dataset) and query (q that we just prepared).
	// Pretty means to format the output.
	searchResult, err := client.Search().
		Index(os.Getenv("ES_INDEX")).
		Query(query).
		Pretty(true).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)
	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d post\n", searchResult.TotalHits())

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization.
	var typ model.Post
	var posts []model.Post
	// Iterate the result results and if they are type of Post (typ)
	for _, item := range searchResult.Each(reflect.TypeOf(typ)) {
		post := item.(model.Post) // Cast an item to Post, equals to p = (Post) item in java
		fmt.Printf("Post by %s: %s at lat %v and lon %v\n", post.User, post.Message, post.Location.Lat, post.Location.Lon)
		// Perform filtering based on keywords such as web spam etc.
		if !containsSensitiveWords(&post.Message) {
			posts = append(posts, post)
		}
	}
	js, err := json.Marshal(posts) // Convert the go object to a string
	if err != nil {
		panic(err)
		return
	}

	// for here, we find result from ES, we need to write result into Redis, use TTL(time to live) as
	// caching strategy to avoid result inconsistent
	if ENABLE_MEMCACHE {
		// build connection with Redis
		rs_client := redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_URL"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})

		// Set the cache expiration to be 30 seconds
		err := rs_client.Set(key, string(js), time.Second*30).Err()
		if err != nil {
			fmt.Printf("Redis cannot save the key %s as %v.\n", key, err)
		}

	}

	res.Write(js)
}

// private method, filter sensitive words
func containsSensitiveWords(post *string) bool {
	data, err := ioutil.ReadFile("../filteredWords.json")
	if err != nil {
		fmt.Println(err)
		return false
	}

	var sensitiveWords []string
	err = json.Unmarshal(data, &sensitiveWords)
	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, word := range sensitiveWords {
		if strings.Contains(*post, word) {
			return true
		}
	}
	return false
}
