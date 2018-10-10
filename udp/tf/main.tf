// Configure the Google Cloud provider
provider "google" {
  credentials = "${file("~/.gcp/kubernetes-helloworld-e1a55b4248a6.json")}"
  project     = "kubernetes-helloworld"
  region      = "us-central1"
}

// Create a new instance
resource "google_compute_instance" "default" {
  
}
