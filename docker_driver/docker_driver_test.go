package docker_driver

import (
  "os"
  "github.com/netapp/netappdvp/storage_drivers"
  "github.com/netapp/netappdvp/storage_drivers/test_driver"
  "testing"
)

func newNdvpDriverWithPrefix(storage_prefix, snapshot_prefix string) (*ndvpDriver) {
  volumeDir := "/tmp/volume"

  commonConfig := &storage_drivers.CommonStorageDriverConfig {}
  commonConfig.Version = 1
  commonConfig.StorageDriverName = "fake"
  commonConfig.Debug = false
  commonConfig.DisableDelete = false
  commonConfig.StoragePrefixRaw = []byte(storage_prefix)
  commonConfig.SnapshotPrefixRaw = []byte(snapshot_prefix)

  fakeDriver := &test_driver.FakeStorageDriver{}
  d, err := NewNetAppDockerVolumePlugin(volumeDir, *commonConfig, fakeDriver)
  if err != nil {
    panic(0)
  }

  return d
}

func TestStoragePrefix(t *testing.T) {
  prefix_cases := []struct {
    in, expected_prefix string
  } {
    {``, "fake_"},
    {`"myprefix_"`, "myprefix_"},
    {`""`, ""},
    {`"a"`, "a"},
    {`"ab"`, "ab"},
    {`"abc"`, "abc"},
  }

  for _, c := range prefix_cases {
    driver := newNdvpDriverWithPrefix(c.in, "")
    got := driver.volumePrefix()
    if got != c.expected_prefix {
      t.Errorf("ndvpDriver.volumePrefix() == %q, expected %q", got, c.expected_prefix)
    }
  }
}

func TestVolumeNames(t *testing.T) {
  volume_name_cases := []struct {
    prefix, volume, expected_volume_name string
  } {
    {``, "vol1", "fake_vol1"},
    {`"myprefix_"`, "vol2", "myprefix_vol2"},
    {`""`, "vol3", "vol3"},
  }

  for _, c := range volume_name_cases {
    driver := newNdvpDriverWithPrefix(c.prefix, "")
    got := driver.volumeName(c.volume)
    if got != c.expected_volume_name {
      t.Errorf("ndvpDriver.volumeName(%q) == %q, expected %q", c.volume, got, c.expected_volume_name)
    }
  }
}

func TestSnapshotPrefix(t *testing.T) {
  //TODO: Clean this up.  Probably need to update somethin to use the specified prefix
  snapshot_prefix_cases := []struct {
    prefix, expected_snapshot_prefix string
  } {
    {``, "fake_"},
    {`""`, ""},
    {`"myprefix"`, "myprefix"},
  }

  for _, c := range snapshot_prefix_cases {
    driver := newNdvpDriverWithPrefix("", c.prefix)
    got := driver.snapshotPrefix()
    if got != c.expected_snapshot_prefix {
      t.Errorf("ndvpDriver.snapshotPrefix() == %q, expected %q", got, c.expected_snapshot_prefix)
    }
  }
}

func TestMountPoint(t *testing.T) {
  driver := newNdvpDriverWithPrefix("", "")
  name := "abcd"
  got := driver.mountpoint(name)
  if got != ("/tmp/volume/" + name) {
    t.Errorf("d.mountpoint(%v) == %v, expected %v", name, got, got)
  }
}

func TestGetMountPoint(t *testing.T) {
  volumeDir := "/tmp/volume"
  os.MkdirAll(volumeDir, 0777)
  defer os.RemoveAll(volumeDir)

  mount_point_cases := []struct {
    storage_prefix, mount_name, expected_prefix, expected_mount_point string
    expected_err error
  } {
    {``, "mount1", "fake_", volumeDir + "/fake_mount1", nil},
    {`"myprefix_"`, "mount1", "myprefix_", volumeDir + "/myprefix_mount1", nil},
    {`""`, "mount1", "", volumeDir + "/mount1", nil},
    {`"a"`, "mount1", "a", volumeDir + "/amount1", nil},
    {`"ab"`, "mount1", "ab", volumeDir + "/abmount1", nil},
    {`"abc"`, "mount1", "abc", volumeDir + "/abcmount1", nil},
  }

  for _, c := range mount_point_cases {
    fpath := volumeDir + "/" + c.expected_prefix + c.mount_name
    f, err := os.Create(fpath)
    if err != nil {
      f.Close()
      t.Errorf("Unable to create %v, err %v", fpath, err)
    }

    driver := newNdvpDriverWithPrefix(c.storage_prefix, "")
    got_volume_prefix := driver.volumePrefix()
    got_mount_point, err := driver.getMountPoint(c.mount_name)
    if got_mount_point != c.expected_mount_point {
      t.Errorf("ndvpDriver.getMountPoint(%v) == %q, expected %q (volumePrefix() = %q)", c.mount_name, got_mount_point, c.expected_mount_point, got_volume_prefix)
    }

    if err != c.expected_err {
      t.Errorf("Unexpected err: %v", err)
    }
    f.Close()
  }
}

func TestNewNetAppDockerVolumePlugin(t *testing.T) {
    root := "/tmp/volume"
    defer os.RemoveAll(root)
    commonConfig := &storage_drivers.CommonStorageDriverConfig {}
    commonConfig.Version = 1
    commonConfig.StorageDriverName = "ontap-nas"
    commonConfig.Debug = false
    commonConfig.DisableDelete = false
    commonConfig.StoragePrefixRaw = []byte(`""`)
    commonConfig.SnapshotPrefixRaw = []byte(`""`)

    fakeDriver := &test_driver.FakeStorageDriver{}

    d, err := NewNetAppDockerVolumePlugin(root, *commonConfig, fakeDriver)
    if err != nil {
      t.Errorf("Problem with %v", d)
    }
}
