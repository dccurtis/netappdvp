package test_driver

import (
	"github.com/netapp/netappdvp/storage_drivers"
	log "github.com/Sirupsen/logrus"
)
type FakeStorageDriverConfig struct {
	storage_drivers.CommonStorageDriverConfig        // embedded types replicate all fields
	ManagementLIF             string `json:"managementLIF"`
	DataLIF                   string `json:"dataLIF"`
	IgroupName                string `json:"igroupName"`
	SVM                       string `json:"svm"`
	Username                  string `json:"username"`
	Password                  string `json:"password"`
	Aggregate                 string `json:"aggregate"`
}

type FakeStorageDriver struct {
  Initialized bool
  Config FakeStorageDriverConfig
}

const FakeStorageDriverName = "fake"

func (d *FakeStorageDriver) Name() string {
	return FakeStorageDriverName
}

func (d *FakeStorageDriver) Initialize(configJSON string) error {

  //TODO: Figure out how to fill in config from configJSON
  config := &FakeStorageDriverConfig{}
  d.Config = *config
  d.Initialized = true

  //TODO:Call Validate d.Validate()
  return nil
}

func (d *FakeStorageDriver) Validate() error {
  //TODO: Validate test data?
	log.Debugf("FakeStorageDriver.Validate()")
  return nil
}

func (d *FakeStorageDriver) Create(name string, opts map[string]string) error {
  //TODO: Add logic once theres a need
	log.Debugf("FakeStorageDriver.Create()- name: %v, opts: %v", name, opts)

  return nil
}

func (d *FakeStorageDriver) CreateClone(name, source, snapshot, newSnapshotPrefix string) error {
	log.Debugf("FakeStorageDriver.CreateClone()- \n\tname: %v, \n\tsource: %v, \n\tsnapshot: %v, \n\tnewSnapshotPrefix: %v",
		name, source, snapshot, newSnapshotPrefix)

  //TODO: Add logic once theres a need
  return nil
}

func (d *FakeStorageDriver) Destroy(name string) error {
	log.Debugf("FakeStorageDriver.Destroy()- \n\tname: %v", name)
  //TODO: Add logic once theres a need
  return nil
}

func (d *FakeStorageDriver) Attach(name, mountpoint string, opts map[string]string) error {
	log.Debugf("FakeStorageDriver.Attach()- name: %v, mountpoint: %v, opts: %v", name, mountpoint, opts)
  //TODO: Add logic once theres a need
  return nil
}

func (d *FakeStorageDriver) Detach(name, mountpoint string) error {
	log.Debugf("FakeStorageDriver.Detach()- name: %v, mountpoint: %v", name, mountpoint)
  //TODO: Add logic once theres a need
  return nil
}

func (d *FakeStorageDriver) DefaultStoragePrefix() string {
	log.Debugf("FakeStorageDriver.DefaultStoragePrefix()")
	return "fake_"
}

func (d *FakeStorageDriver) DefaultSnapshotPrefix() string {
	log.Debugf("FakeStorageDriver.DefaultSnapshotPrefix()")
	return "fake_"
}

func (d *FakeStorageDriver) SnapshotList(name string) ([]storage_drivers.CommonSnapshot, error) {
	log.Debugf("FakeStorageDriver.SnapshotList()- name: %v", name)
  var snapshots []storage_drivers.CommonSnapshot
  //TODO: Add necessary stuff here
  return snapshots, nil
}
