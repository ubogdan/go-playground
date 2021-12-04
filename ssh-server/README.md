SSH server example
------------------

### Build ssh server
```shell
go build -o ssh_server
 ```

### Generate ssh server key
```shell
ssh-keygen -f id_rsa
```

### Run ssh server
```shell
sudo su -c ./ssh_server
```