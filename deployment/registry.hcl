job "registry" {
  datacenters = ["dc1"]

  group "registry" {
    network {
      mode = "bridge"
    }

    service {
      name = "image-registry"
      port = 5000

      connect {
        sidecar_service {}
      }
    }

    task "registry" {
      driver = "docker"

      config {
        image = "registry:2"
      }
    }
  }
}