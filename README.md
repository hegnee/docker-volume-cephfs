# docker-volume-cephfs<br>
CephFS docker volume driver plugin.<br>
MOUNT CEPH FS WITH THE KERNEL DRIVER<br>
Usage:<br>
./docker-volume-cephfs<br>
docker run --name vu1  -d -P --volume-driver=cephfs -v test:/mnt/foo -it ubuntu bash<br>
