job "rabbitmq" {
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
}