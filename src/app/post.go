package main

import (
	"context"
	"fmt"
	"io"

	// "mime/multipart"
	"model"
	"net/http"
	"strconv"

	"cloud.google.com/go/bigtable"
	"cloud.google.com/go/storage"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pborman/uuid"
	elastic "gopkg.in/olivere/elastic.v3"
)

// how to parse multipart form in Go?
// https://golang.org/pkg/net/http/#Request.ParseMultipartForm
// https://github.com/golang-samples/http/blob/master/fileupload/main.go

func Post(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

	if req.Method != "POST" {
		fmt.Println(req.Method)
		return
	}

	// get username from context instead of letting user to input username themselves
	user := req.Context().Value("user")
	if user == nil {
		m := fmt.Sprintf("Unable to find user in context")
		fmt.Println(m)
		http.Error(res, m, http.StatusBadRequest)
		return
	}
	claims := user.(*jwt.Token).Claims
	username := claims.(jwt.MapClaims)["username"]

	// 32 << 20 is the maxMemory param for ParseMultipartForm, equals to 32MB (1MB = 1024 * 1024 bytes = 2^20 bytes)
	// After you call ParseMultipartForm, the file will be saved in the server memory with maxMemory size.
	// If the file size is larger than maxMemory, the rest of the data will be saved in a system temporary file.
	err := req.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Println(err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse from form data.
	fmt.Printf("Received one post request %s\n", req.FormValue("message"))
	lat, _ := strconv.ParseFloat(req.FormValue("lat"), 64)
	lon, _ := strconv.ParseFloat(req.FormValue("lon"), 64)
	post := &model.Post{
		User:    username.(string),
		Message: req.FormValue("message"),
		Location: model.Location{
			Lat: lat,
			Lon: lon,
		},
	}

	id := uuid.New()

	files := req.MultipartForm.File["images[]"]

	// testLocal(images)
	for i, f := range files {
		// open uploaded file
		file, err := f.Open()
		defer file.Close()
		if err != nil {
			// log.Fatal(err);
			http.Error(res, "Image is not available", http.StatusInternalServerError)
			fmt.Println(err)
			continue
		}

		// attrs is file
		_, attrs, err := saveToGCS(context.Background(), file, BUCKET_NAME, fmt.Sprintf("%s%d", id, i))
		if err != nil {
			http.Error(res, "GCS is not setup", http.StatusInternalServerError)
			fmt.Printf("GCS is not setup %v\n", err)
			return
		}

		// Update the media link after saving to GCS.
		post.URLs = append(post.URLs, attrs.MediaLink)
	}

	ctx := context.Background()

	// Save to ES.
	go saveToES(ctx, post, id)

	// Save to BigTable as well.
	if ENABLE_BIGTABLE {
		go saveToBigTable(ctx, post, id)
	}
}

// private method, Save a post to ElasticSearch
func saveToES(ctx context.Context, post *model.Post, id string) {
	post.Mu.Lock()
	defer post.Mu.Unlock()

	// Create a client
	es_client, err := elastic.NewClient(elastic.SetURL(ES_URL), elastic.SetSniff(false))
	if err != nil {
		panic(err)
		return
	}

	// Save it to index, example taken from https://github.com/olivere/elastic
	_, err = es_client.Index().
		Index(INDEX).
		Type(TYPE).
		Id(id).
		BodyJson(post).
		Refresh("true").
		Do(ctx)
	if err != nil {
		panic(err)
		return
	}

	fmt.Printf("Post is saved to Index: %s\n", post.Message)
}

// private method, save image to Google Cloud Storage
// Google example of open a client connection to GCS
// https://cloud.google.com/storage/docs/reference/libraries#client-libraries-install-go
// Google example of writing an object to GCS, see write function
// https://github.com/GoogleCloudPlatform/golang-samples/blob/master/storage/objects/main.go
// r is image file, name is id of this image
func saveToGCS(ctx context.Context, reader io.Reader, bucketName, name string) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer client.Close()

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)
	// Next check if the bucket exists
	if _, err = bucket.Attrs(ctx); err != nil {
		return nil, nil, err
	}

	obj := bucket.Object(name) // stored file
	writer := obj.NewWriter(ctx)
	// copy(write) image file to GCS's bucket
	if _, err := io.Copy(writer, reader); err != nil {
		return nil, nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, nil, err
	}

	// already finish writting, now modify access permission to all users, ACL is access control list
	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return nil, nil, err
	}

	attrs, err := obj.Attrs(ctx)
	fmt.Printf("Post is saved to GCS: %s\n", attrs.MediaLink)
	return obj, attrs, err
}

// private method, Save a post to BigTable
// how to write obj to bigtable reference: https://cloud.google.com/bigtable/docs/samples-go-hello
// type conversion: string to byte array. byteArray := []byte(myString)
func saveToBigTable(ctx context.Context, post *model.Post, id string) {
	post.Mu.Lock()
	defer post.Mu.Unlock()

	bt_client, err := bigtable.NewClient(ctx, PROJECT_ID, BT_INSTANCE)
	if err != nil {
		fmt.Println(err)
		panic(err)
		return
	}
	tbl := bt_client.Open("post")
	mut := bigtable.NewMutation()
	t := bigtable.Now()
	// params: columnFamilyName, columnName, timestamp value
	mut.Set("post", "user", t, []byte(post.User))
	mut.Set("post", "message", t, []byte(post.Message))
	mut.Set("location", "lat", t, []byte(strconv.FormatFloat(post.Location.Lat, 'f', -1, 64)))
	mut.Set("location", "lon", t, []byte(strconv.FormatFloat(post.Location.Lon, 'f', -1, 64)))

	err = tbl.Apply(ctx, id, mut)
	if err != nil {
		panic(err)
		return
	}
	fmt.Printf("Post is saved to BigTable: %s\n", post.Message)
}

// func testLocal(files []*multipart.FileHeader) {
// 	for _, f := range files {
// 		// open uploaded file
// 		file, err := f.Open()
// 		defer file.Close()
// 		if err != nil {
// 			// log.Fatal(err);
// 			fmt.Println(err)
// 		}
// 		// create local folder
// 		os.Mkdir("./upload", os.ModePerm)

// 		// save to local file
// 		cur, err := os.Create("./upload/" + f.Filename)
// 		defer cur.Close()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		io.Copy(cur, file)
// 	}
// }
