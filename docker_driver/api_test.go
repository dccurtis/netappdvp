package docker_driver

import (
  "os"
  "testing"
  "github.com/docker/go-plugins-helpers/volume"
  log "github.com/Sirupsen/logrus"
)

func newDriver() (*ndvpDriver) {
  d := newNdvpDriverWithPrefix(``, "")
  return d
}

func cleanup() {
  os.RemoveAll(tempRoot)
}

func createVolumeFromList(d *ndvpDriver, volumes []string) {
    for _, vol := range volumes {
      request := volume.Request {
        Name: vol,
        Options: map[string]string{},
      }
      response := d.Create(request)

      if response.Err != "" {
        panic(0)
      }
      }
}

func TestCreate(t *testing.T) {
  log.Infof("Docker Volume Interface Create(): Starting")

  create_cases := []struct {
    request volume.Request
    expected_path string
  } {
  { volume.Request {
    Name: "myvolume",
    Options: map[string]string{}},
    tempRoot + "/fake_myvolume"},

    { volume.Request {
      Name: "myvolume",
      //test_driver is being used and it currenlty doesn't do anythign for snaphots
      Options: map[string]string{"from": "myvol", "fromSnapshot": "mysnap"}},
      tempRoot + "/fake_myvolume"},
  }

  for _, c := range create_cases {
    d := newDriver()
    defer cleanup()
    response := d.Create(c.request)

    if response.Err != "" {
      log.Infof("Docker Volume Interface Create(): Failed")
      t.Errorf("response: Mountpoint %v, Err: %s", response.Mountpoint, response.Err)
    }

    _, err := os.Stat(c.expected_path)
    if os.IsNotExist(err) {
      log.Infof("Docker Volume Interface Create(): Failed")
      t.Errorf("ndvpDriver.Create() expected directory (%s) does not exist.", c.expected_path)
    }
    cleanup()
  }

  log.Infof("Docker Volume Interface Create(): Passed")
}

func TestList(t *testing.T) {
  log.Infof("Docker Volume Interface List(): Starting")
  list_cases := []struct {
    request volume.Request
    req_volumes []string
    expected_vol_qty int
  } {
  {volume.Request {
    Name: "doesnt_do_anything",
    Options: map[string]string{}},
    []string{},
    0},
  {volume.Request {
    Name: "doesnt_do_anything",
    Options: map[string]string{}},
    []string{"testvolume"},
    1},
  {volume.Request {
    Name: "doesnt_do_anything",
    Options: map[string]string{}},
    []string{"testvolume", "secondvolume"},
    2},
  }

  for _, c := range list_cases {
    d := newDriver()
    defer cleanup()

    createVolumeFromList(d, c.req_volumes)

    list_response := d.List(c.request)

    if list_response.Err != "" {
      log.Infof("Docker Volume Interface List(): Failed")
      t.Errorf("ndvpDriver.List() Failed: Mountpoint %v, Err: %s", list_response.Mountpoint, list_response.Err)
    }

    //Verify the volume list
    if (len(list_response.Volumes) != c.expected_vol_qty) {
      t.Errorf("ndvpDriver.List() Unexpected number of volumes, found: %s, expected: %s",
        len(list_response.Volumes), c.expected_vol_qty)
    }

    //Check volume names
    for idx, vol := range list_response.Volumes {
        if c.req_volumes[idx] != vol.Name {
          t.Errorf("ndvpDriver.List() Unexpected volume name.  found: %s, expected %s",
          list_response.Volumes, c.req_volumes)
        }
        _, err := os.Stat(vol.Mountpoint)
        if os.IsNotExist(err) {
          log.Infof("Docker Volume Interface List(): Failed")
          t.Errorf("ndvpDriver.List() expected directory (%s) does not exist.", vol.Mountpoint)
        }
    }
    cleanup()
  }
  log.Infof("Docker Volume Interface List(): Passed")
}

func TestGet(t *testing.T) {
  log.Infof("Docker Volume Interface Get(): Starting")
  get_cases := []struct {
    request volume.Request
    req_volumes []string
    expected_vol_qty int
  } {
  {volume.Request {
    Name: "testvolume",
    Options: map[string]string{}},
    []string{"testvolume"},
    1},
  }

  for _, c := range get_cases {
    d := newDriver()
    defer cleanup()

    createVolumeFromList(d, c.req_volumes)

    get_response := d.Get(c.request)

    if get_response.Err != "" {
      log.Infof("Docker Volume Interface Get(): Failed")
      t.Errorf("ndvpDriver.Get() Failed: Err: %s", get_response.Err)
    }

    //Verify the volume name
    if (get_response.Volume.Name != c.request.Name) {
      t.Errorf("ndvpDriver.get() Unexpected volume found: %s, expected: %s",
        get_response.Volume.Name, c.request.Name)
    }

    expected_mount_point := tempRoot + "/fake_" + c.request.Name
    if get_response.Volume.Mountpoint != expected_mount_point {
      log.Infof("Docker Volume Interface Get(): Failed")
      t.Errorf("volume_get_response Mountpoint: %s, expected: %s", get_response.Volume.Mountpoint, expected_mount_point)
    }

    cleanup()
  }

  log.Infof("Docker Volume Interface Get(): Passed")
}

func TestRemove(t *testing.T) {
  log.Infof("Docker Volume Interface Remove(): Starting")
  remove_cases := []struct {
    request volume.Request
    req_volumes []string
    expected_vol_qty int
  } {
  {volume.Request {
    Name: "testvolume"},
    []string{"testvolume"},
    1},
  }

  for _, c := range remove_cases {
    d := newDriver()
    defer cleanup()

    createVolumeFromList(d, c.req_volumes)

    remove_response := d.Remove(c.request)

    if remove_response.Err != "" {
      log.Infof("Docker Volume Interface Remove(): Failed")
      t.Errorf("vol_remove_response: Err: %s", remove_response.Err)
    }
    cleanup()
  }

  log.Infof("Docker Volume Interface Remove(): Passed")
}

func TestPath(t *testing.T) {
  log.Infof("Docker Volume Interface Path(): Starting")
  path_cases := []struct {
    request volume.Request
    req_volumes []string
    expected_vol_qty int
  } {
  {volume.Request {
    Name: "testvolume",
    Options: map[string]string{}},
    []string{"testvolume"},
    1},
  }

  for _, c := range path_cases {
    d := newDriver()
    defer cleanup()

    createVolumeFromList(d, c.req_volumes)

    path_response := d.Path(c.request)

    if path_response.Err != "" {
      log.Infof("Docker Volume Interface Path(): Failed")
      t.Errorf("Unexpected err: %s", path_response.Err)
    }

    expected_path := tempRoot + "/fake_" + c.request.Name
    if path_response.Mountpoint != expected_path {
      log.Infof("Docker Volume Interface Path(): Failed")
      t.Errorf("Unexpected volume path: got %s, expected: %s", path_response.Mountpoint, expected_path)
    }

    cleanup()
  }
  log.Infof("Docker Volume Interface Path(): Passed")
}

func TestMount(t *testing.T) {
  log.Infof("Docker Volume Interface Mount(): Starting")
  mount_cases := []struct {
    request volume.MountRequest
    req_volumes []string
    expected_vol_qty int
  } {
  {volume.MountRequest {
    Name: "testvolume"},
    []string{"testvolume"},
    1},
  }

  for _, c := range mount_cases {
    d := newDriver()
    defer cleanup()

    createVolumeFromList(d, c.req_volumes)

    mount_response := d.Mount(c.request)

    if mount_response.Err != "" {
      log.Infof("Docker Volume Interface Mount(): Failed")
      t.Errorf("Unexpected err: %s", mount_response.Err)
    }

    expected_mount := tempRoot + "/fake_" + c.request.Name
    if mount_response.Mountpoint != expected_mount {
      log.Infof("Docker Volume Interface Mount(): Failed")
      t.Errorf("Unexpected volume mount: got %s, expected: %s", mount_response.Mountpoint, expected_mount)
    }

    cleanup()
  }
  log.Infof("Docker Volume Interface Mount(): Passed")
}

func TestUnmount(t *testing.T) {
  log.Infof("Docker Volume Interface Unmount(): Starting")
  unmount_cases := []struct {
    request volume.UnmountRequest
    req_volumes []string
    expected_vol_qty int
  } {
  {volume.UnmountRequest {
    Name: "testvolume"},
    []string{"testvolume"},
    1},
  }

  for _, c := range unmount_cases {
    d := newDriver()
    defer cleanup()

    createVolumeFromList(d, c.req_volumes)

    unmount_response := d.Unmount(c.request)

    if unmount_response.Err != "" {
      log.Infof("Docker Volume Interface Unmount(): Failed")
      t.Errorf("Unexpected err: %s", unmount_response.Err)
    }
    cleanup()
  }
  log.Infof("Docker Volume Interface Unmount(): Passed")
}

func TestCapabilities(t *testing.T) {
  log.Infof("Docker Volume Interface Capabilities(): Starting")

  d := newDriver()
  defer cleanup()

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
