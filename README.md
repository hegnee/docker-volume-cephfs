# docker-volume-cephfs
CephFS docker volume driver plugin.
MOUNT CEPH FS WITH THE KERNEL DRIVER
Usage:
./docker-volume-cephfs
docker run --name vu1  -d -P --volume-driver=cephfs -v test:/mnt/foo -it ubuntu bash
