package data

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// Disk structure for lsblk JSON parsing
type BlockDevice struct {
	Name     string        `json:"name"`
	Size     string        `json:"size"`
	Type     string        `json:"type"`
	Children []BlockDevice `json:"children,omitempty"`
}

type LsblkOutput struct {
	BlockDevices []BlockDevice `json:"blockdevices"`
}

// GetDisks returns a list of formatted disk strings for dropdowns
func GetDisks() []string {
	cmd := exec.Command("lsblk", "-d", "-n", "-o", "NAME,SIZE,TYPE", "--json")
	output, err := cmd.Output()
	if err != nil {
		return []string{"/dev/sda (Test)", "/dev/nmve0n1 (Test)"}
	}

	var data LsblkOutput
	if err := json.Unmarshal(output, &data); err != nil {
		return []string{"Error parsing disks"}
	}

	var disks []string
	for _, dev := range data.BlockDevices {
		if dev.Type == "disk" {
			disks = append(disks, fmt.Sprintf("/dev/%s (%s)", dev.Name, dev.Size))
		}
	}
	if len(disks) == 0 {
		disks = append(disks, "No disks found")
	}
	return disks
}

// GetPartitions returns partitions for a given disk (or all if empty)
func GetPartitions() []string {
	// We want all partitions to let user select Root/EFI
	cmd := exec.Command("lsblk", "-l", "-n", "-o", "NAME,SIZE,TYPE", "--json")
	output, err := cmd.Output()
	if err != nil {
		return []string{"/dev/sda1 (Test)", "/dev/sda2 (Test)"}
	}

	var data LsblkOutput
	if err := json.Unmarshal(output, &data); err != nil {
		return []string{"Error parsing partitions"}
	}

	var parts []string
	for _, dev := range data.BlockDevices {
		// lsblk -l output flat list, sometimes 'part' type
		if dev.Type == "part" {
			parts = append(parts, fmt.Sprintf("/dev/%s (%s)", dev.Name, dev.Size))
		}
	}
	if len(parts) == 0 {
		return []string{"No partitions found"}
	}
	return parts
}
