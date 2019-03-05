resource "google_storage_bucket" "image-store" {
  name     = "${var.BUCKET_NAME}"
  location = "US"
}

output "gcs_url" {
  # The base URL of the bucket, in the format gs://<bucket-name>.
  value = "${google_storage_bucket.image-store.url}"
}

output "gcs_self_link" {
  # The URI of the created resource.
  value = "${google_storage_bucket.image-store.self_link}"
}