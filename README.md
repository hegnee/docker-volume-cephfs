# docker-volume-cephfs<br>
##CephFS docker volume driver plugin.<br>
<br>
##Create CephFS
ceph osd pool create <pool_name> <pg_num>  <br>
ceph osd pool create cephfs_data 64 <br>
ceph osd pool create cephfs_metadata 64 <br>
ceph fs new <fs_name> <metadata_pool> <data_pool>  <br>
ex:ceph fs new mycephfs cephfs_metadata cephfs_data <br>  
##MOUNT CEPH FS WITH THE KERNEL DRIVER<br>
Usage:<br>
./docker-volume-cephfs<br>
docker run --name vu1  -d -P --volume-driver=cephfs -v test:/mnt/foo -it ubuntu bash<br>
