package net

import (
	"fmt"

	"github.com/rivo/tview"
)

func getIfaceTree(iface Iface) *tview.TreeNode {
	root := tview.NewTreeNode(fmt.Sprintf("%s (%s)", iface.Name, iface.Type))
	root.AddChild(tview.NewTreeNode(fmt.Sprintf("MTU: %d", iface.MTU)))
	root.AddChild(tview.NewTreeNode(fmt.Sprintf("State: %s", iface.State)))

	if len(iface.IPv4.Addresses) > 0 {
		IPv4Node := tview.NewTreeNode("IPv4 Addresses")
		for _, address := range iface.IPv4.Addresses {
			IPv4Node.AddChild(tview.NewTreeNode(address.String()))
		}
		root.AddChild(IPv4Node)
	}
	if len(iface.IPv6.Addresses) > 0 {
		IPv6Node := tview.NewTreeNode("IPv6 Addresses")
		for _, address := range iface.IPv6.Addresses {
			IPv6Node.AddChild(tview.NewTreeNode(address.String()))
		}
		root.AddChild(IPv6Node)
	}

	return root
}

func TreeView(netState NetState, pages *tview.Pages) (*tview.TreeView, error) {
	if pages == nil {
		return nil, fmt.Errorf("Can't make a NetState treeView page for nil pages")
	}

	root := tview.NewTreeNode(fmt.Sprintf("%s Network State", netState.Hostname.Running))
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	interfaces := tview.NewTreeNode("Interfaces")
	root.AddChild(interfaces)

	defaultIface := netState.GetDefaultNextHopIface()
	if defaultIface != nil {
		interfaces.AddChild(getIfaceTree(*defaultIface))
	}
	for _, iface := range netState.Ifaces {
		if defaultIface != nil && defaultIface.Name == iface.Name {
			continue // Skip defaultRouteIface, since we always display it first
		}
		interfaces.AddChild(getIfaceTree(iface))
	}

	return tree, nil
}
