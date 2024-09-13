## Volume Bitmap Utility

### Description

This project provides a utility to interact with the volume bitmap of a Windows filesystem. It includes functions to retrieve and analyze the volume bitmap, which can be used to determine the allocation status of clusters on the disk.

### Features

- Retrieve volume bitmap from a specified volume (vss volume)
- Check if a range of clusters is allocated (i.e., contains data)

### Installation

```sh
go get github.com/8ugMak1r/go-win-bitmap
```

```go
import "github.com/8ugMak1r/go-win-bitmap"
```

### Usage

#### Retrieve Volume Bitmap

To retrieve the volume bitmap for a specified volume:

```go
package main

import (
    "fmt"
    "log"
    "github.com/8ugMak1r/go-win-bitmap"
)

func main() {
    volume := `C:\`
    volsize := uint64(500 * 1024 * 1024 * 1024) // Example volume size: 500GB
    clustersize := uint32(4096) // Example cluster size: 4KB

    vbb, err := bitmap.GetVolumeBitmapBuffer(volume, volsize, clustersize)
    if err != nil {
        log.Fatalf("Failed to get volume bitmap: %v", err)
    }

}
```

#### Retrieve Volume Bitmap for VSS Volume

To retrieve the volume bitmap for a VSS (Volume Shadow Copy Service) volume:

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/8ugMak1r/go-win-bitmap"
)

func main() {
	vssVolume := `\\?\GLOBALROOT\Device\HarddiskVolumeShadowCopy1\`
	volsize := uint64(500 * 1024 * 1024 * 1024) // Example volume size: 500GB
	clustersize := uint32(4096)                 // Example cluster size: 4KB

	bm, err := bitmap.GetVolumeBitmapBuffer(vssVolume, volsize, clustersize)
	if err != nil {
		log.Fatalf("Failed to get VSS volume bitmap: %v", err)
	}

	// read volume according bitmap
	v, err := os.Open(vssVolume)
	if err != nil {
		log.Fatal(err)
	}
	defer v.Close()
	start := 0
	end := 8
	if bm.Allocated(start, end) {
		blockSize := int(clustersize) * 8 * (end - start)
		// read data from vss volume
		offset := start * int(clustersize) * 8
		data := make([]byte, blockSize)
		n, err := v.ReadAt(data, int64(offset))
		if err != nil {
			//
		}

	}
}
```

### References

- [VOLUME\_BITMAP\_BUFFER structure](https://learn.microsoft.com/en-us/windows/win32/api/winioctl/ns-winioctl-volume_bitmap_buffer)
- [FSCTL\_GET\_VOLUME\_BITMAP control code](https://learn.microsoft.com/en-us/windows/win32/api/winioctl/ni-winioctl-fsctl_get_volume_bitmap)

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.