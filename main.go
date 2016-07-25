package main

import (
	"fmt"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	dkvolume "github.com/docker/go-plugins-helpers/volume"
)


var (
	pluginName = flag.String("name","cephfs","Docker plugin name for use on --volume-driver option")
	cephcluter = flag.String("cluster","ceph","Ceph cluster")
	rootMountDir = flag.String("mount",dkvolume.DefaultDockerRootDirectory, "Mount directory for volumes on host")
	pluginDir          = flag.String("plugins", "/run/docker/plugins", "Docker plugin directory for socket")
	logDir             = flag.String("logdir", "/var/log", "Logfile directory")
	cephConfigFile     = flag.String("config", "", "Ceph cluster config")
)

func init(){
	flag.Parse()
}

func socketPath() string{
	socket_name := *pluginName+".sock"
	//return filepath.Join(*pluginDir,*pluginName)
	return filepath.Join(*pluginDir,socket_name)
}

func logfilePath() string{
	return filepath.Join(*logDir,*pluginName+"-docker-cephfs.log")
}

func isDebugEnable() bool{
	return os.Getenv("CEPHFS_DOCKER_DEBUG") == "1"
}

func setupLogging() (*os.File, error){
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix(*pluginName + "-volume-plugin: ")
	logfileName := logfilePath()
	if !isDebugEnable() && logfileName != ""{
		logFile,err := os.OpenFile(logfileName,os.O_RDWR|os.O_CREATE|os.O_APPEND,0666)
		if err != nil {
			if os.IsPermission(err){
				log.Printf("WARN:logging fallback to STDERR: %v",err)
			}else{
				return nil,err
			}
		}else{
			log.Printf("INFO: setting log file: %s",logfileName)
			log.SetOutput(logFile)
			return logFile,nil
		}
	}
	return nil,nil
}

func shutdownLogging(logFile *os.File){
	if logFile != nil {
		log.Println("INFO: closing log file")
		logFile.Sync()
		logFile.Close()
	}
}

func reloadLogging(logFile *os.File) (*os.File,error) {
	log.Println("INFO: reloading log")
	if logFile != nil {
		shutdownLogging(logFile)
	}
	return setupLogging()
}

func main(){
	logFile,err := setupLogging()
	if err != nil {
		log.Panicf("Unable to setup logging: %s",err)
	}
	defer shutdownLogging(logFile)
	d := newCephfsDriver(
		*pluginName,
		*cephcluter,
		*rootMountDir,
		*cephConfigFile,
	)
	log.Println("INFO: Creating Docker VolumeDriver Handler")
	h := dkvolume.NewHandler(d)
	socket := socketPath()
	log.Printf("INFO: Opening Socket for Docker to Connect: %s",socket)
	err = os.MkdirAll(filepath.Dir(socket),os.ModeDir)
	if err != nil {
		log.Panicf("Error Creating socket directory: %s",err)
	}
	signalChannel := make(chan os.Signal,2)
	signal.Notify(signalChannel,syscall.SIGTERM,syscall.SIGKILL,syscall.SIGHUP)
	go func() {
		for sig:= range signalChannel {
			switch sig {
				case syscall.SIGTERM,syscall.SIGKILL:
					log.Printf("INFO: received TERM or KILL signal: %s",sig)
					shutdownLogging(logFile)
					os.Exit(0)
				case syscall.SIGHUP:
					log.Printf("INFO: received HUG signal: %s",sig)
					logFile,err = reloadLogging(logFile)
					if err != nil {
						log.Printf("Unable to reload log:%s",err)
					}
					d.connect()
			}
		}
	}()
	log.Printf("Start ServeUnix..")
	//fmt.Println(h.ServeUnix("root", "cephfs"))
	fmt.Println(h.ServeUnix("root", socket))
	/*
	err = h.ServeUnix("",socket)
	if err != nil {
		log.Printf("ERROR: Unable to create Unix socket: %v",err)
	}
	*/
}
