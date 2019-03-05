provider "google" {
  credentials = "${file("${var.GOOGLE_APPLICATION_CREDENTIALS}")}"
  project     = "${var.GCP_PROJECT_ID}"
  region      = "${var.GCP_REGION}"
}