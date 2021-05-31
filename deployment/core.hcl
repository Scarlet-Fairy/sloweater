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

    service {
      name = "rabbit-manager"
      port = 15672

      check {
        name     = "alive"
        type     = "http"
        path     = "/"
        interval = "10s"
        timeout  = "2s"
      }

    }

    task "rabbitmq" {
      driver = "docker"

      config {
        image = "rabbitmq:3-management-alpine"

        logging {
          type = "elastic/elastic-logging-plugin:7.12.1"

          config {
            hosts = "http://localhost:9200"
            index = "scarlet-fairy-core"
          }
        }
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

  group "elasticsearch" {
    network {
      mode = "bridge"

      port "api" {
        to = 9200
        static = 9200
      }

      port "transport" {
        to = 9300
        static = 9300
      }
    }

    service {
      name = "elasticsearch-api"
      port = 9200


      connect {
        sidecar_service {}
      }
    }

    task "elasticsearch" {
      driver = "docker"

      config {
        image = "docker.elastic.co/elasticsearch/elasticsearch:7.12.1"

        ulimit {
          nofile = "65535:65535"
        }
      }

      env = {
        "bootstrap.memory_lock" = "true"
        ES_JAVA_OPTS = "-Xms512m -Xmx512m"
        "discovery.type" = "single-node"
      }

      resources {
        cpu = 1000
        disk = 512
        memory = 1024
      }
    }
  }

  group "kibana" {
    network {
      mode = "bridge"

      port "ui" {
        to = 5601
        static = 5601
      }
    }

    service {
      name = "kibana-ui"
      port = 5601

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "elasticsearch-api"
              local_bind_port  = 9200
            }
          }
        }
      }
    }

    task "kibana" {
      driver = "docker"

      config {
        image = "docker.elastic.co/kibana/kibana:7.12.1"
      }

      env {
        ELASTICSEARCH_HOSTS = "http://localhost:9200"
      }
    }
  }
}