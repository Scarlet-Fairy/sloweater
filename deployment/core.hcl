job "core" {
  datacenters = ["dc1"]

  group "rabbitmq" {
    network {
      mode = "bridge"

      port "queues" {
        static = 5672
        to = 5672
      }

      port "manager" {
        static = 15672
        to = 15672
      }
    }

    service {
      name = "rabbitmq"
      port = 5672

      connect {
        sidecar_service {}
      }
    }

    task "rabbitmq" {
      driver = "docker"

      config {
        image = "rabbitmq:3-management-alpine"
      }
    }
  }

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