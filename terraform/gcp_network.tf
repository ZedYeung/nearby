resource "google_compute_network" "default" {
  name                    = "${var.APP}-network"
  auto_create_subnetworks = "true"
}

resource "google_compute_firewall" "default" {
  name    = "${var.APP}-network"
  network = "${google_compute_network.default.name}"

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["80", "8080", "8081", "9200", "9300", "5601", "22"]
  }
}