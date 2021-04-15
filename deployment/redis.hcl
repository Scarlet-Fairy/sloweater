job "redis" {
  datacenters = ["dc1"]

  group "redis" {
    network {
      mode = "bridge"
    }

    service {
      name = "redis"
      port = 6379

      connect {
        sidecar_service {}
      }
    }

    task "redis" {
      driver = "docker"

      config {
        image = "redis:alpine"
      }
    }
  }
}