job "busybox" {
  datacenters = ["dc1"]

  group "busybox" {
    network {
      mode = "bridge"
    }

    job "busybox" {
      driver = "docker"

      config {
        image = "busybox"
        command = "sleep"
        args = [
          "infinity"
        ]
      }
    }
  }
}