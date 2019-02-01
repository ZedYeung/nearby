# Nearby-backend
Geospatial based photo shared social network backend with Golang<br>
Frontend is implementated with ReactJS, please see [here](https://github.com/ZedYeung/nearby-frontend).

# Overview of project
- Develop Web services with **Golang** to handle user auth(register, login, logout), post and search
- Implement token based authentication with **PostgreSQL**, **JWT**, **OAuth 2.0**
- Implement posts storage, geo-location based search for nearby posts, data visualization  with **ElasticSearch** and **Kibana** in **GKE(Google Kubernetes Engine)**
- Store post image with **Google Cloud Storage(GCS)**
- Improve read performance with a little data consistency sacrifice with **Redis**
- Deploy with **Docker**, **GKE(Google Kubernetes Engine)** and **GOOGLE CLOUD LOAD BALANCING**

# API
- **/signup**
  * save to PostgreSQL.
- **/login**
  * check login credential in PostgreSQL, if correct return token
- **/search** - search nearby posts
  1. have token-based authentication first
  2. search in Redis cache, if not found then search ElasticSearch
  3. use  `"type" : "geo_point"` to map (lat, lon) to geo_point, ElasticSearch will use geo-indexing to search(**KD tree**)

  e.g.
  ```
  /search?lat=10.0&lon=20.0
  /search?lat=10.0&lon=20.0&range=10
  ```
- **/post**
  1. save post image in GCS. 
  2. save post info in ElasticSearch, bigTable(optional).

  post body
  ```
  {
    "message": "<3",
    "location": {
      "lat": 47.651977307189256,
      "lon": -122.3316657342046
    },
    "url": [
      "https://www.googleapis.com/download/storage/v1/b/nearby-posts/o/d46f6e06-669b-41ec-9bb4-25a6e08c259f0?generation=1548962551529115&alt=media"
    ]
  }
  ```

# TODO
* CRUD for user's own posts
* Follow and Like

