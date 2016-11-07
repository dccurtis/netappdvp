package test_driver

import (
  "testing"
)

func TestInitialize(t *testing.T) {
  d := FakeStorageDriver{}
  config := "config-string"
  d.Initialize(config)

  if d.Initialized == false {
    t.Errorf("Fake Driver Initialized failed")
  }
}

func TestValidate(t *testing.T) {
  d := FakeStorageDriver{}
  err := d.Validate()

  if err != nil {
    t.Errorf("Fake Driver Validation failed")
  }
}

func TestCreate(t *testing.T) {
  d := FakeStorageDriver{}
  name := "myvolume"
  opts := make(map[string]string)

  err := d.Create(name, opts)

  if err != nil {
    t.Errorf("Fake Driver Create failed")
  }
}

func TestCreateClone(t *testing.T) {
  d := FakeStorageDriver{}
  name := "myname"
  source := "source"
  snapshot := "snapshot"
  snapshot_prefix := "prefix"
  err := d.CreateClone(name, source, snapshot, snapshot_prefix)

  if err != nil {
    t.Errorf("Fake Driver Create Clone failed")
  }
}

func TestDestroy(t *testing.T) {
  d := FakeStorageDriver{}
  name := "name"
  err := d.Destroy(name)

  if err != nil {
    t.Errorf("Fake Driver Destroy failed")
  }
}

func TestAttach(t *testing.T) {
  d := FakeStorageDriver{}
  name := "name"
  mountpoint := "mount"
  opts := make(map[string]string)

  err := d.Attach(name, mountpoint, opts)

  if err != nil {
    t.Errorf("Fake Driver Attach failed")
  }
}

func TestDetach(t *testing.T) {
  d := FakeStorageDriver{}
  name := "name"
  mountpoint := "mountpoint"
  err := d.Detach(name, mountpoint)

  if err != nil {
    t.Errorf("Fake Driver Detach failed")
  }
}

func TestDefaultStoragePrefix(t *testing.T) {
  d := FakeStorageDriver{}
  prefix := d.DefaultStoragePrefix()

  if prefix != "fake_" {
    t.Errorf("Fake DefaultStoragePrefix failed")
  }
}

func TestDefaultSnapshotPrefix(t *testing.T) {
  d := FakeStorageDriver{}
  prefix := d.DefaultSnapshotPrefix()

  if prefix != "fake_" {
    t.Errorf("Fake DefaultSnapshotPrefix failed")
  }
}

func TestSnapshotList(t *testing.T) {
  d := FakeStorageDriver{}
  name := "name"
  snapshots, _ := d.SnapshotList(name)

  if len(snapshots) != 0 {
    t.Errorf("Fake Snapshot List failed")
  }
}
