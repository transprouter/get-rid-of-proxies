provider "docker" {
  host = "unix:///var/run/docker.sock"
}

resource "docker_network" "transprouter_private" {
  name = "transprouter_private"
}

resource "docker_container" "private_proxy" {
  name = "proxy.private"
  image = "${docker_image.proxy.latest}"
  networks = [
    "${docker_network.transprouter_private.id}",
    "${docker_network.transprouter_public.id}",
  ]
}

resource "docker_container" "workstation" {
  name  = "workstation.private"
  image = "${docker_image.workstation.latest}"
  networks = ["${docker_network.transprouter_private.id}"]
  privileged = true
}

resource "docker_container" "private_webserver" {
  name = "webserver.private"
  image = "${docker_image.nginx.latest}"
  networks = ["${docker_network.transprouter_private.id}"]
}

resource "docker_container" "private_sshserver" {
  name = "sshserver.private"
  image = "${docker_image.openssh.latest}"
  networks = ["${docker_network.transprouter_private.id}"]
}

resource "docker_network" "transprouter_public" {
  name = "transprouter_public"
}

resource "docker_container" "public_webserver" {
  name = "webserver.public"
  image = "${docker_image.nginx.latest}"
  networks = ["${docker_network.transprouter_public.id}"]
}

resource "docker_container" "public_sshserver" {
  name = "sshserver.public"
  image = "${docker_image.openssh.latest}"
  networks = ["${docker_network.transprouter_public.id}"]
}

resource "docker_image" "workstation" {
  name = "transprouter/workstation:latest"
  keep_locally = true
}

resource "docker_image" "proxy" {
  name = "transprouter/proxy:latest"
  keep_locally = true
}

resource "docker_image" "nginx" {
  name = "transprouter/nginx:latest"
  keep_locally = true
}

resource "docker_image" "openssh" {
  name = "transprouter/openssh:latest"
  keep_locally = true
}
