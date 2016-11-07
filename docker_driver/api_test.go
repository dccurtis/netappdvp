package docker_driver

import (
  "testing"
  "github.com/docker/go-plugins-helpers/volume"
)

func TestCreate(t *testing.T) {
  request := volume.Request{
    Name: "myvolume",
    Options: map[string]string{"from": "myvol", "fromSnapshot": "mysnap"},
  }

  d := newNdvpDriverWithPrefix(``, "")
  response := d.Create(request)

  if response.Err != "" {
    t.Errorf("response: Mountpoint %v, Err: %s", response.Mountpoint, response.Err)
  }
}

func TestList(t *testing.T) {
  request := volume.Request{
    Name: "doesnt_do_anything",
    Options: map[string]string{"from": "myvol", "fromSnapshot": "mysnap"},
  }

  d := newNdvpDriverWithPrefix(``, "")
  response := d.List(request)
  if response.Err != "" {
    t.Errorf("response: Mountpoint %v, Err: %s", response.Mountpoint, response.Err)
  }

  if response.Volumes[0].Name != "myvolume" {
    t.Errorf("List(%v) = %v, expected: %v", request, response.Volumes[0], "requested_name")
  }
  //t.Errorf("response.Volumes: %v", response.Volumes[0].Name)
}
