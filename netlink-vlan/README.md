Manage vlan's using netlink
---                   

Add vlan and assign ip address
```shell
  ./netlink-vlan add -i eth0 -vlan 10 -addr 10.0.0.1/24
```

Remove vlan
```shell
  ./netlink-vlan del -i eth0 -vlan 10
```
