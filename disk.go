package main

import (
  "fmt"
  "time"
  "syscall"
)

func Disk(mount string) {
  totalKey := fmt.Sprintf("disk.%s.total", cleanKey(mount))
  freeKey := fmt.Sprintf("disk.%s.free", cleanKey(mount))
  for {
    var stat syscall.Statfs_t
    if err := syscall.Statfs(mount, &stat); err != nil {
      fmt.Println("statfs(%q) failed with: %v", mount, err)
    } else {
      // blocks * size per block = bytes
      sendGauge(freeKey, int64(stat.Bavail) * stat.Bsize)
      sendGauge(totalKey, int64(stat.Blocks) * stat.Bsize)
    }
    time.Sleep(1 * time.Second)
  }
}
