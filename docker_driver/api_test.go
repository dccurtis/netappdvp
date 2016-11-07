package docker_driver

import (
  "os"
  "testing"
  "github.com/docker/go-plugins-helpers/volume"
  log "github.com/Sirupsen/logrus"
)

func TestCreate(t *testing.T) {
  log.Infof("Docker Volume Interface Create(): Starting")

  request := volume.Request{
    Name: "myvolume",
    Options: map[string]string{"from": "myvol", "fromSnapshot": "mysnap"},
  }

  d := newNdvpDriverWithPrefix(``, "")
  defer os.RemoveAll("/tmp/volume")
  response := d.Create(request)

  if response.Err != "" {
    log.Infof("Docker Volume Interface Create(): Failed")
    t.Errorf("response: Mountpoint %v, Err: %s", response.Mountpoint, response.Err)
  }
  log.Infof("Docker Volume Interface Create(): Passed")
}

func TestList(t *testing.T) {
  log.Infof("Docker Volume Interface List(): Starting")

  d := newNdvpDriverWithPrefix(``, "")
  defer os.RemoveAll("/tmp/volume")

  // Create a volume
  volume_request := volume.Request {
    Name: "testvolume",
    Options: map[string]string{"from:": "myvol"},
  }

  vol_create_response := d.Create(volume_request)
  if vol_create_response.Err != "" {
    log.Infof("Docker Volume Interface List(): Failed")
    t.Errorf("Unexpected volume create error, Err: %s", vol_create_response.Err)
  }

  list_request := volume.Request {
    Name: "doesnt_do_anything",
    Options: map[string]string{"from": "myvol", "fromSnapshot": "mysnap"},
  }

  list_response := d.List(list_request)
  if list_response.Err != "" {
    log.Infof("Docker Volume Interface List(): Failed")
    t.Errorf("response: Mountpoint %v, Err: %s", list_response.Mountpoint, list_response.Err)
  }

  if list_response.Volumes[0].Name != "testvolume" {
    log.Infof("Docker Volume Interface List(): Failed")
    t.Errorf("List(%v) = %v, expected: %v", list_request, list_response.Volumes[0], "requested_name")
  }
  log.Infof("Docker Volume Interface List(): Passed")
}

func TestGet(t *testing.T) {
  log.Infof("Docker Volume Interface Get(): Starting")

  d := newNdvpDriverWithPrefix(``, "")
  defer os.RemoveAll("/tmp/volume")

  // Create a volume
  volume_request := volume.Request {
    Name: "testvolume",
    Options: map[string]string{"from:": "myvol", "fromSnapshot": "snapz"},
  }

  vol_create_response := d.Create(volume_request)
  if vol_create_response.Err != "" {
    log.Infof("Docker Volume Interface Get(): Failed")
    t.Errorf("Unexpected volume create error, Err: %s", vol_create_response.Err)
  }

  vol_get_response := d.Get(volume_request)

  if vol_get_response.Volume.Name != "testvolume" {
    log.Infof("Docker Volume Interface Get(): Failed")
    t.Errorf("volume_get_response Name: %s, expected: %s", vol_get_response.Volume.Name, "testvolume")
  }

  if vol_get_response.Volume.Mountpoint != "/tmp/volume/fake_testvolume" {
    log.Infof("Docker Volume Interface Get(): Failed")
    t.Errorf("volume_get_response Mountpoint: %s, expected: %s", vol_get_response.Volume.Mountpoint, "/tmp/volume/fake_testvolume")
  }
  log.Infof("Docker Volume Interface Get(): Passed")

}

func TestRemove(t *testing.T) {
  log.Infof("Docker Volume Interface Remove(): Starting")

  d := newNdvpDriverWithPrefix(``, "")
  defer os.RemoveAll("/tmp/volume")

  // Create a volume
  volume_request := volume.Request {
    Name: "testvolume",
    Options: map[string]string{"from:": "myvol"},
  }

  vol_create_response := d.Create(volume_request)
  if vol_create_response.Err != "" {
    t.Errorf("Unexpected volume create error, Err: %s", vol_create_response.Err)
  }

  volume_remove_request := volume.Request {
    Name: "testvolume",
  }

  vol_remove_response := d.Remove(volume_remove_request)
  if vol_remove_response.Err != "" {
    log.Infof("Docker Volume Interface Remove(): Failed")
    t.Errorf("vol_remove_response: Err: %s", vol_remove_response.Err)
  }
  log.Infof("Docker Volume Interface Remove(): Passed")
}

func TestPath(t *testing.T) {
  log.Infof("Docker Volume Interface Path(): Starting")

  d := newNdvpDriverWithPrefix(``, "")
  defer os.RemoveAll("/tmp/volume")

  // Create a volume
  volume_request := volume.Request {
    Name: "testvolume",
    Options: map[string]string{"from:": "myvol"},
  }

  vol_create_response := d.Create(volume_request)
  if vol_create_response.Err != "" {
    log.Infof("Docker Volume Interface Path(): Failed")
    t.Errorf("Unexpected volume create error, Err: %s", vol_create_response.Err)
  }

  path_request := volume.Request {
    Name: "testvolume",
  }

  path_response := d.Path(path_request)

  if path_response.Err != "" {
    log.Infof("Docker Volume Interface Path(): Failed")
    t.Errorf("Unexpected err: %s", path_response.Err)
  }
  if path_response.Mountpoint != "/tmp/volume/fake_testvolume" {
    log.Infof("Docker Volume Interface Path(): Failed")
    t.Errorf("Unexpected volume path: got %s, expected: %s", path_response, "/tmp/volume/fake_testvolume")
  }
  log.Infof("Docker Volume Interface Path(): Passed")

}

func TestMount(t *testing.T) {
  log.Infof("Docker Volume Interface Mount(): Starting")

  d := newNdvpDriverWithPrefix(``, "")
  defer os.RemoveAll("/tmp/volume")

  // Create a volume
  volume_request := volume.Request {
    Name: "testvolume",
    Options: map[string]string{"from:": "myvol"},
  }

  vol_create_response := d.Create(volume_request)
  if vol_create_response.Err != "" {
    log.Infof("Docker Volume Interface Mount(): Failed")
    t.Errorf("Unexpected volume create error, Err: %s", vol_create_response.Err)
  }

  mount_request := volume.MountRequest {
    Name: "testvolume",
  }

  mount_response := d.Mount(mount_request)

  if mount_response.Err != "" {
    log.Infof("Docker Volume Interface Mount(): Failed")
    t.Errorf("Unexpected err: %s", mount_response.Err)
  }

  if mount_response.Mountpoint != "/tmp/volume/fake_testvolume" {
    log.Infof("Docker Volume Interface Mount(): Failed")
    t.Errorf("Unexpected mount point: got %s, expected: %s", mount_response, "/tmp/volume/fake_testvolume")
  }
  log.Infof("Docker Volume Interface Mount(): Passed")
}

func TestUnmount(t *testing.T) {
  log.Infof("Docker Volume Interface Unmount(): Starting")

  d := newNdvpDriverWithPrefix(``, "")
  defer os.RemoveAll("/tmp/volume")

  // Create a volume
  volume_request := volume.Request {
    Name: "testvolume",
    Options: map[string]string{"from:": "myvol"},
  }

  vol_create_response := d.Create(volume_request)
  if vol_create_response.Err != "" {
    log.Infof("Docker Volume Interface Unmount(): Failed")
    t.Errorf("Unexpected volume create error, Err: %s", vol_create_response.Err)
  }

  unmount_request := volume.UnmountRequest {
    Name: "testvolume",
  }

  unmount_response := d.Unmount(unmount_request)

  if unmount_response.Err != "" {
    log.Infof("Docker Volume Interface Unmount(): Failed")
    t.Errorf("Unexpected err: %s", unmount_response.Err)
  }
  log.Infof("Docker Volume Interface Unmount(): Passed")

}

func TestCapabilities(t *testing.T) {
  log.Infof("Docker Volume Interface Capabilities(): Starting")

  d := newNdvpDriverWithPrefix(``, "")
  defer os.RemoveAll("/tmp/volume")

  capability_request := volume.Request {
    Name: "doesntmatter",
  }
  capability_response := d.Capabilities(capability_request)

  if capability_response.Capabilities.Scope != "global" {
    log.Infof("Docker Volume Interface Capabilities(): Failed")
    t.Errorf("Capability got %s, expected: %s", capability_response.Capabilities.Scope, "global")
  }
  log.Infof("Docker Volume Interface Capabilities(): Passed")

}
