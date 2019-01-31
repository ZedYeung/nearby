package model

type Config struct {
	SIGN_KEY        string
	DB_USER         string
	DB_PASSWORD     string
	DB_NAME         string
	INDEX           string
	TYPE            string
	DISTANCE        string
	PROJECT_ID      string
	BT_INSTANCE     string
	ES_URL          string
	REDIS_URL       string
	REDIS_PASSWORD  string
	REDIS_DB        string
	BUCKET_NAME     string
	ENABLE_MEMCACHE bool
	ENABLE_BIGTABLE bool
}
