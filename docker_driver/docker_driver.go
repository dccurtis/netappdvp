// Copyright 2016 NetApp, Inc. All Rights Reserved.

package docker_driver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/netapp/netappdvp/storage_drivers"

	log "github.com/Sirupsen/logrus"
)

type ndvpDriver struct {
	m      *sync.Mutex
	root   string
	config storage_drivers.CommonStorageDriverConfig
	sd     storage_drivers.StorageDriver
}

func (d *ndvpDriver) volumePrefix() string {
	defaultPrefix := d.sd.DefaultStoragePrefix()
	prefixToUse := defaultPrefix
	storagePrefixRaw := d.config.StoragePrefixRaw // this is a raw version of the json value, we will get quotes in it
	if len(storagePrefixRaw) >= 2 {
		s := string(storagePrefixRaw)
		if s == "\"\"" || s == "" {
			prefixToUse = ""
			//log.Debugf("storagePrefix is specified as \"\", using no prefix")
		} else {
			// trim quotes from start and end of string
			prefixToUse = s[1 : len(s)-1]
			//log.Debugf("storagePrefix is specified, using prefix: %v", prefixToUse)
		}
	} else {
		prefixToUse = defaultPrefix
		//log.Debugf("storagePrefix is unspecified, using default prefix: %v", prefixToUse)
	}

	return prefixToUse
}

func (d *ndvpDriver) volumeName(name string) string {
	prefixToUse := d.volumePrefix()
	if strings.HasPrefix(name, prefixToUse) {
		return name
	}
	return prefixToUse + name
}

func (d *ndvpDriver) snapshotPrefix() string {
	defaultPrefix := d.sd.DefaultSnapshotPrefix()
	prefixToUse := defaultPrefix
	snapshotPrefixRaw := d.config.SnapshotPrefixRaw // this is a raw version of the json value, we will get quotes in it
	if len(snapshotPrefixRaw) >= 2 {
		s := string(snapshotPrefixRaw)
		if s == "\"\"" || s == "" {
			prefixToUse = ""
			//log.Debugf("snapshotPrefix is specified as \"\", using no prefix")
		} else {
			// trim quotes from start and end of string
			prefixToUse = s[1 : len(s)-1]
			//log.Debugf("snapshotPrefix is specified, using prefix: %v", prefixToUse)
		}
	} else {
		prefixToUse = defaultPrefix
		//log.Debugf("snapshotPrefix is unspecified, using default prefix: %v", prefixToUse)
	}

	return prefixToUse
}

func (d *ndvpDriver) mountpoint(name string) string {
	return filepath.Join(d.root, name)
}

func NewNetAppDockerVolumePlugin(root string, config storage_drivers.CommonStorageDriverConfig, storage_d storage_drivers.StorageDriver) (*ndvpDriver, error) {
	// if root (volumeDir) doesn't exist, make it
	dir, err := os.Lstat(root)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(root, 0755); err != nil {
			return nil, err
		}
	}
	// if root (volumeDir) isn't a directory, error
	if dir != nil && !dir.IsDir() {
		return nil, fmt.Errorf("Volume directory '%v' exists and it's not a directory", root)
	}

	d := &ndvpDriver{
		root:   root,
		config: config,
		m:      &sync.Mutex{},
		sd:     storage_d,
	}
	return d, nil
}

func (d ndvpDriver) getMountPoint(requestName string) (string, error) {
	target := d.volumeName(requestName)
	m := d.mountpoint(target)
	log.Debugf("Getting path for volume '%s' as '%s'", target, m)

	fi, err := os.Lstat(m)
	if os.IsNotExist(err) {
		return "", err
	}
	if fi == nil {
		return "", fmt.Errorf("Could not stat %v", m)
	}

	return d.mountpoint(target), nil
}
