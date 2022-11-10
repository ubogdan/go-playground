package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/vishvananda/netlink"
)

func addVlan(masterInterface string, vlanID int, ipNet *net.IPNet) error {
	master, err := netlink.LinkByName(masterInterface)
	if err != nil {
		log.Fatalf("failed to lookup master %q: %v", master, err)
	}

	vlan := &netlink.Vlan{
		LinkAttrs: netlink.LinkAttrs{
			Name:        fmt.Sprintf("%s.%d", master.Attrs().Name, vlanID),
			MTU:         master.Attrs().MTU,
			ParentIndex: master.Attrs().Index,
		},
		VlanId: vlanID,
	}

	err = netlink.LinkAdd(vlan)
	if err != nil {
		return fmt.Errorf("failed to create vlan: %v", err)
	}

	link, err := netlink.LinkByName(vlan.Name)
	if err != nil {
		return fmt.Errorf("failed to lookup master %q: %v", vlan.Name, err)
	}

	err = netlink.LinkSetUp(link)
	if err != nil {
		return fmt.Errorf("failed to set %q UP: %v", link.Attrs().Name, err)
	}

	if ipNet != nil {
		addr := &netlink.Addr{IPNet: ipNet, Label: ""}

		err = netlink.AddrAdd(link, addr)
		if err != nil {
			return fmt.Errorf("failed to add IP addr %v to %q: %v", ipNet.IP.String(), link.Attrs().Name, err)
		}
	}

	return nil
}

// removeVlan godoc.
func removeVlan(masterInterface string, vlanID int) error {
	linkName := fmt.Sprintf("%s.%d", masterInterface, vlanID)

	link, err := netlink.LinkByName(linkName)
	if err != nil {
		return fmt.Errorf("failed to lookup vlan %s: %v", linkName, err)
	}

	err = netlink.LinkSetDown(link)
	if err != nil {
		return fmt.Errorf("failed to set %q DOWN: %v", link.Attrs().Name, err)
	}

	err = netlink.LinkDel(link)
	if err != nil {
		return fmt.Errorf("failed to delete %q: %v", link.Attrs().Name, err)
	}

	return nil
}

var masterIF, vlanIP string

func main() {
	//masterIF := "enp0s31f6"
	vlanID := 10
	//vlanIP := "10.0.0.1/24"

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addCmd.StringVar(&masterIF, "i", "eth0", "Master interface")
	addCmd.IntVar(&vlanID, "vlan", 0, "VLAN ID")
	addCmd.StringVar(&vlanIP, "addr", "", "VLAN IP address")

	rmCmd := flag.NewFlagSet("rm", flag.ExitOnError)
	rmCmd.StringVar(&masterIF, "i", "eth0", "master interface")
	rmCmd.IntVar(&vlanID, "vlan", 0, "VLAN ID")

	flag.Parse()

	args := flag.Args()

	switch args[0] {
	case addCmd.Name():
		addCmd.Parse(args[1:])

		var ipNet *net.IPNet

		if vlanIP != "" {
			net, err := netlink.ParseIPNet(vlanIP)
			if err != nil {
				log.Fatalf("failed to parse IPNet %q: %v", vlanIP, err)
			}
			ipNet = net
		}

		err := addVlan(masterIF, vlanID, ipNet)
		if err != nil {
			log.Fatalf("failed to add vlan: %v", err)
		}

	case rmCmd.Name():
		rmCmd.Parse(args[1:])

		err := removeVlan(masterIF, vlanID)
		if err != nil {
			log.Fatalf("failed to remove vlan: %v", err)
		}
	}

}
