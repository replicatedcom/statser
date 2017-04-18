package main

import (
	"fmt"
	"syscall"
	"time"
)

func Disk(mount string) {
	totalKey := fmt.Sprintf("disk.%s.total", cleanKey(mount))
	freeKey := fmt.Sprintf("disk.%s.free", cleanKey(mount))
	for {
		var stat syscall.Statfs_t
		if err := syscall.Statfs(mount, &stat); err != nil {
			fmt.Printf("statfs(%q) failed with: %v\n", mount, err)
		} else {
			// blocks * size per block = bytes
			sendGauge(freeKey, int64(stat.Bavail*uint64(stat.Bsize)))
			sendGauge(totalKey, int64(stat.Blocks*uint64(stat.Bsize)))
		}
		time.Sleep(1 * time.Second)
	}
}
