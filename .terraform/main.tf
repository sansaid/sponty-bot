resource "google_cloud_run_service" "default" {
  name     = local.app_name
  location = local.default_location

  template {
    spec {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/hello"
      }
    }
  }
}