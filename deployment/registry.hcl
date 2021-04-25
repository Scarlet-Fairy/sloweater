job "registry" {
  datacenters = ["dc1"]

  group "registry" {
    network {
      mode = "bridge"

      port "registry" {
        to = 5000
        static = 5000
      }
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