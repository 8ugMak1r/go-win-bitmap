//go:build windows

package bitmap

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	FSCTL_GET_VOLUME_BITMAP = 0x9006f
)

type STARTING_LCN_INPUT_BUFFER struct {
	StartingLcn uint64
}

type VOLUME_BITMAP_BUFFER struct {
	StartingLcn uint64
	BitmapSize  uint64
	Buffer      []byte
}

// Allocated checks if there is any data in the cluster range [start, end).
func (b *VOLUME_BITMAP_BUFFER) Allocated(start, end int) bool {
	if start < 0 || end < 0 || start > end || end > len(b.Buffer) {
		return false
	}
	for _, value := range b.Buffer[start:end] {
		if value > 0 {
			return true
		}
	}
	return false
}

func calculateRequiredBufferSize(diskCapacity uint64, clusterSize uint32) uint64 {
	totalClusters := diskCapacity / uint64(clusterSize)
	requiredBufferSize := (totalClusters + 7) / 8
	return requiredBufferSize + 16
}

func GetVolumeBitmapBuffer(volume string, volsize uint64, clustersize uint32) (*VOLUME_BITMAP_BUFFER, error) {
	handle, err := windows.CreateFile(
		windows.StringToUTF16Ptr(volume),
		windows.GENERIC_READ,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		0,
		0,
	)
	if err != nil {
		return nil, fmt.Errorf("error opening volume: %w", err)
	}
	defer windows.CloseHandle(handle)

	requiredBufferSize := calculateRequiredBufferSize(volsize, clustersize)
	output := make([]byte, requiredBufferSize+16)
	input := STARTING_LCN_INPUT_BUFFER{StartingLcn: 0}

	var bytesReturned uint32
	err = windows.DeviceIoControl(
		handle,
		FSCTL_GET_VOLUME_BITMAP,
		(*byte)(unsafe.Pointer(&input)),
		uint32(unsafe.Sizeof(input)),
		&output[0],
		uint32(len(output)),
		&bytesReturned,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("DeviceIoControl error: %w", err)
	}

	outputBuffer := (*VOLUME_BITMAP_BUFFER)(unsafe.Pointer(&output[0]))
	bitmapSizeInBytes := (outputBuffer.BitmapSize + 7) / 8

	if uint64(len(output)) < 16+bitmapSizeInBytes {
		return nil, fmt.Errorf("buffer size is too small")
	}
	bitmap := output[16 : 16+bitmapSizeInBytes]

	return &VOLUME_BITMAP_BUFFER{
		StartingLcn: outputBuffer.StartingLcn,
		BitmapSize:  outputBuffer.BitmapSize,
		Buffer:      bitmap,
	}, nil
}
