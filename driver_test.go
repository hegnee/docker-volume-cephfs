package main

import (
	"fmt"
	"os"
	"testing"
	dkvolume "github.com/docker/go-plugins-helpers/volume"
	"github.com/stretchr/testify/assert"
)
var (
	testDriver cephfsDriver
)

func TestMain(m *testing.M){
	cephConf := os.Getenv("CEPH_CONF")
	testDriver = newCephfsDriver(
		"test",
		"ceph",
		dkvolume.DefaultDockerRootDirectory,
		cephConf,
	)
	code := m.Run()
	defer os.Exit(code)
}

func TestDirExists(t *testing.T) {
	_,f_bool,err := testDriver.DirExists("/t1")
//	fmt.Println(f_bool)
	assert.Equal(t,true,f_bool,fmt.Sprintf("%s", err))
}

func TestgetMds(t *testing.T){
	host,err := getMds()
	fmt.Println(host)
	assert.Equal(t,"ubuntu-1",host,fmt.Sprintf("%s",err))
}

