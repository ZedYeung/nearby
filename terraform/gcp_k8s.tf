resource "google_container_cluster" "primary" {
  name = "${var.APP}-cluster"

  # region = "${var.GCP_REGION}"
  zone               = "${var.GCP_ZONE}"
  network            = "${google_compute_network.default.name}"
  min_master_version = "${var.GKE_VERSION}"

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true

  initial_node_count = 1

  # additional_zones = [
  #     "us-west1-a",
  #     "us-west1-c",
  # ]

  addons_config {
    kubernetes_dashboard {
      disabled = false
    }

    http_load_balancing {
      disabled = false
    }
  }
  # Setting an empty username and password explicitly disables basic auth
  master_auth {
    username = ""
    password = ""
  }
  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]

    labels = {
      app = "${var.APP}"
    }

    tags = ["${var.APP}"]
  }

  # depends_on =[
  #   "google_compute_network.default"
  # ]
}

resource "google_container_node_pool" "primary_preemptible_nodes" {
  name = "my-node-pool"

  # region     = "${var.GCP_REGION}"
  zone       = "${var.GCP_ZONE}"
  cluster    = "${google_container_cluster.primary.name}"
  node_count = 1

  node_config {
    preemptible = true

    # machine_type = "n1-standard-1"
    machine_type = "${var.GCP_INSTANCE_TYPE}"

    oauth_scopes = [
      "https://www.googleapis.com/auth/compute",
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
  }

  # provisioner "local-exec" {
  #   command     = "export ES_URL=${google_container_cluster.primary.endpoint}"
  #   interpreter = ["bash", "-c"]
  # }
}

# The following outputs allow authentication and connectivity to the GKE Cluster
# by using certificate-based authentication.
output "client_certificate" {
  value = "${google_container_cluster.primary.master_auth.0.client_certificate}"
}

output "client_key" {
  value = "${google_container_cluster.primary.master_auth.0.client_key}"
}

output "cluster_ca_certificate" {
  value = "${google_container_cluster.primary.master_auth.0.cluster_ca_certificate}"
}

output "endpoint" {
  value = "${google_container_cluster.primary.endpoint}"
}

output "cluster_name" {
  value = "${google_container_cluster.primary.name}"
}

output "cluster_region" {
  value = "${var.GCP_REGION}"
}

output "cluster_zone" {
  value = "${google_container_cluster.primary.zone}"
}
