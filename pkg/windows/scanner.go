//go:build windows

package windows

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func GetFileVersionInfo(filePath string) (FileVersionInfo, error) {
	info := FileVersionInfo{
		FilePath: filePath,
	}

	size, err := windows.GetFileVersionInfoSize(filePath, nil)
	if err != nil {
		return info, err
	}
	if size == 0 {
		return info, fmt.Errorf("GetFileVersionInfoSize returned 0 for file: %s", filePath)
	}

	buf := make([]byte, size)
	err = windows.GetFileVersionInfo(filePath, 0, uint32(size), unsafe.Pointer(&buf[0]))
	if err != nil {
		return info, err
	}

	queryValue := func(key string) string {
		var block *uint16
		var blockLen uint32
		err = windows.VerQueryValue(unsafe.Pointer(&buf[0]), key, unsafe.Pointer(&block), &blockLen)
		if err != nil || blockLen == 0 {
			return ""
		}
		return syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(block))[:blockLen])
	}

	var block *uint16
	var blockLen uint32
	err = windows.VerQueryValue(unsafe.Pointer(&buf[0]), `\VarFileInfo\Translation`, unsafe.Pointer(&block), &blockLen)
	if err != nil || blockLen == 0 {
		// 기본값 사용
		langID := "040904b0"
		info.ProductName = queryValue(`\StringFileInfo\` + langID + `\ProductName`)
		info.CompanyName = queryValue(`\StringFileInfo\` + langID + `\CompanyName`)
		info.FileVersion = queryValue(`\StringFileInfo\` + langID + `\FileVersion`)
		info.ProductVersion = queryValue(`\StringFileInfo\` + langID + `\ProductVersion`)
	} else {
		translations := unsafe.Slice((*uint32)(unsafe.Pointer(block)), blockLen/4)
		for _, rawValue := range translations {
			langID := (rawValue >> 16) & 0xFFFF
			codePage := rawValue & 0xFFFF
			strLangID := fmt.Sprintf("%04x%04x", codePage, langID)

			pName := queryValue(`\StringFileInfo\` + strLangID + `\ProductName`)
			if pName != "" {
				info.ProductName = pName
				info.CompanyName = queryValue(`\StringFileInfo\` + strLangID + `\CompanyName`)
				info.FileVersion = queryValue(`\StringFileInfo\` + strLangID + `\FileVersion`)
				info.ProductVersion = queryValue(`\StringFileInfo\` + strLangID + `\ProductVersion`)
				break
			}
		}
	}

	return info, nil
}

func ScanSystem(paths []string) ([]FileVersionInfo, error) {
	var products []FileVersionInfo
	for _, rootDir := range paths {
		filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if ext := filepath.Ext(path); ext != ".exe" && ext != ".dll" {
				return nil
			}

			fileInfo, err := GetFileVersionInfo(path)
			if err != nil {
				return nil
			}
			if fileInfo.ProductName != "" || fileInfo.CompanyName != "" ||
				fileInfo.FileVersion != "" || fileInfo.ProductVersion != "" {
				products = append(products, fileInfo)
			}
			return nil
		})
	}
	return products, nil
}
