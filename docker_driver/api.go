// Copyright 2016 NetApp, Inc. All Rights Reserved.

package docker_driver

import (
  "fmt"
  "os"
  "strings"
  "path/filepath"

  "github.com/docker/go-plugins-helpers/volume"
  "github.com/netapp/netappdvp/utils"

  log "github.com/Sirupsen/logrus"
)

func (d ndvpDriver) Create(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("Create(%v)", r)

	opts := r.Options
	target := d.volumeName(r.Name)
  log.Debugf("target: %v", target) //Added
	m := d.mountpoint(target)
  log.Debugf("mountpoint: %v", m) //Added

	fi, err := os.Lstat(m)
	if os.IsNotExist(err) {
    log.Debugf("os.IsNotExist = true") //Added
		if err := os.MkdirAll(m, 0755); err != nil {
			return volume.Response{Err: err.Error()}
		}
	} else if err != nil {
		return volume.Response{Err: err.Error()}
	}

	if fi != nil && !fi.IsDir() {
		return volume.Response{Err: fmt.Sprintf("%v already exists and it's not a directory", m)}
	}

	var createErr error

	// If 'from' is specified, create a snapshot and a clone rather than a new empty volume
	if from, ok := opts["from"]; ok {
		source := d.volumeName(from)
    log.Debugf("source: %v", source) //Added

		// If 'fromSnapshot' is specified, we use the existing snapshot instead
		snapshot := opts["fromSnapshot"]
		createErr = d.sd.CreateClone(target, source, snapshot, d.snapshotPrefix())
    log.Debugf("Calling sd.CreateClone with d.snapshotPrefix(), target: %v, source %v, snapshot: %v", target, source, snapshot) //Added
	} else {
		createErr = d.sd.Create(target, opts)
	}

	if createErr != nil {
		os.Remove(m)
    log.Debugf("createErr != nil and removing m: %v", m) //Added
		return volume.Response{Err: fmt.Sprintf("Error creating storage: %v", createErr)}
	}

	return volume.Response{}
}

func (d ndvpDriver) List(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("List(%v)", r)

	// open directory ...
	volumeDir := d.root
	dir, err := os.Open(volumeDir)
	if err != nil {
		return volume.Response{Err: fmt.Sprintf("Problem opening directory %v, error: %v", volumeDir, err)}
	}
	defer dir.Close()

	// stat the directory
	fi, err := dir.Stat()
	if err != nil {
		return volume.Response{Err: fmt.Sprintf("Problem stating directory %v, error: %v", volumeDir, err)}
	}
	if !fi.IsDir() {
		return volume.Response{Err: fmt.Sprintf("%v is not a directory!", volumeDir)}
	}

	// finally, we spin through all the subdirectories (if any) and return them in our List response
	var vols []*volume.Volume
	dirs := make([]string, 0)   // lint complains to switch to this, but it doens't work -> var dirs []string
	fis, err := dir.Readdir(-1) // -1 means return all the FileInfos
	if err != nil {
		return volume.Response{Err: fmt.Sprintf("Problem reading directory %v, error: %v", volumeDir, err)}
	}
	for _, fileinfo := range fis {
		if fileinfo.IsDir() {
			dirs = append(dirs, fileinfo.Name())

			// removes the prefix based on prefix length, for instance [10:] to remove 'netappdvp_' from start of name
			volumePrefix := d.volumePrefix()

			// only trim if it matches the prefix
			if strings.HasPrefix(fileinfo.Name(), volumePrefix) {
				volumeName := fileinfo.Name()[len(volumePrefix):]
				log.Debugf("List() adding volume: %v from: %v", volumeName, filepath.Join(volumeDir, fileinfo.Name()))

				v := &volume.Volume{Name: volumeName, Mountpoint: filepath.Join(volumeDir, fileinfo.Name())}
				vols = append(vols, v)
			} else {
				log.Debugf("wrong prefix, skipping fileinfo.Name: %v", fileinfo.Name())
			}
		}
	}

	return volume.Response{Volumes: vols}
}

func (d ndvpDriver) Get(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("Get(%v)", r)

	// Gather the target volume name as the storage sees it
	target := d.volumeName(r.Name)

	path, err := d.getMountPoint(r.Name)
	if err != nil {
		return volume.Response{Err: err.Error()}
	}

	// Ask the storage driver for the list of snapshots associated with the volume
	snaps, err := d.sd.SnapshotList(target)

	// If we don't get any snapshots, that's fine. We'll return an empty list.
	status := map[string]interface{}{
		"Snapshots": snaps,
	}

	v2 := &volume.Volume{
		Name:       r.Name,
		Mountpoint: path,
		Status:     status, // introduced in Docker 1.12, earlier versions ignore
	}

	return volume.Response{
		Volume: v2,
	}
}

func (d ndvpDriver) Remove(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("Remove(%v)", r)

	target := d.volumeName(r.Name)

	// allow user to completely disable volume deletion
	if d.config.DisableDelete {
		log.Infof("Skipping removal of %s because of user preference to disable volume deletion", target)
		return volume.Response{}
	}

	log.Debugf("Removing docker volume %s", target)

	m := d.mountpoint(target)

	fi, err := os.Lstat(m)
	if os.IsNotExist(err) {
		return volume.Response{} // nothing to do
	} else if err != nil {
		return volume.Response{Err: err.Error()}
	}

	if fi != nil && !fi.IsDir() {
		return volume.Response{Err: fmt.Sprintf("%v is not a directory", m)}
	}

	// use the StorageDriver to destroy the storage objects
	destroyErr := d.sd.Destroy(target)
	if destroyErr != nil {
		return volume.Response{Err: fmt.Sprintf("Problem removing docker volume: %v error: %v", target, destroyErr)}
	}

	log.Debugf("rmdir(%s)", m)
	err3 := os.Remove(m)
	if err3 != nil {
		return volume.Response{Err: err3.Error()}
	}

	return volume.Response{}
}

func (d ndvpDriver) Path(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("Path(%v)", r)

	path, err := d.getMountPoint(r.Name)
	if err != nil {
		return volume.Response{Err: err.Error()}
	}

	return volume.Response{
		Mountpoint: path,
	}
}

func (d ndvpDriver) Mount(r volume.MountRequest) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("Mount(%v)", r)

	target := d.volumeName(r.Name)

	m := d.mountpoint(target)
	log.Debugf("Mounting volume %s on %s", target, m)

	fi, err := os.Lstat(m)

	if os.IsNotExist(err) {
		if err := os.MkdirAll(m, 0755); err != nil {
			return volume.Response{Err: err.Error()}
		}
	} else if err != nil {
		return volume.Response{Err: err.Error()}
	}

	if fi != nil && !fi.IsDir() {
		return volume.Response{Err: fmt.Sprintf("%v already exists and it's not a directory", m)}
	}

	// check if already mounted before we do anything...
	dfOuput, dfOuputErr := utils.GetDFOutput()
	if dfOuputErr != nil {
		return volume.Response{Err: fmt.Sprintf("Error checking if %v is already mounted: %v", m, dfOuputErr)}
	}
	for _, e := range dfOuput {
		if e.Target == m {
			log.Debugf("%v already mounted, returning existing mount", m)
			return volume.Response{Mountpoint: m}
		}
	}

	// use the StorageDriver to attach the storage objects, place any extra options in this map
	attachOptions := make(map[string]string)

	attachErr := d.sd.Attach(target, m, attachOptions)
	if attachErr != nil {
		log.Error(attachErr)
		return volume.Response{Err: fmt.Sprintf("Problem attaching docker volume: %v mountpoint: %v error: %v", target, m, attachErr)}
	}

	return volume.Response{Mountpoint: m}
}

func (d ndvpDriver) Unmount(r volume.UnmountRequest) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("Unmount(%v)", r)

	target := d.volumeName(r.Name)

	m := d.mountpoint(target)
	log.Debugf("Unmounting docker volume %s", target)

	// use the StorageDriver to unmount the storage objects
	detachErr := d.sd.Detach(target, m)
	if detachErr != nil {
		return volume.Response{Err: fmt.Sprintf("Problem unmounting docker volume: %v error: %v", target, detachErr)}
	}

	return volume.Response{}
}

func (d ndvpDriver) Capabilities(r volume.Request) volume.Response {
	d.m.Lock()
	defer d.m.Unlock()

	log.Debugf("Capabilities(%v)", r)

	return volume.Response{Capabilities: volume.Capability{Scope: "global"}}
}
