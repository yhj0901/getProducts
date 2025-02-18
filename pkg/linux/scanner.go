//go:build linux

package linux

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetELFInfo(filePath string) (FileVersionInfo, error) {
	info := FileVersionInfo{
		FilePath: filePath,
	}

	// 파일 타입 확인
	fileCmd := exec.Command("file", "-b", filePath)
	output, err := fileCmd.Output()
	if err != nil {
		return info, fmt.Errorf("file command failed: %v", err)
	}
	info.FileType = strings.TrimSpace(string(output))

	// dpkg로 버전 정보 확인 시도
	dpkgCmd := exec.Command("dpkg", "-S", filePath)
	if output, err := dpkgCmd.Output(); err == nil {
		parts := strings.Split(string(output), ":")
		if len(parts) > 0 {
			packageName := strings.TrimSpace(parts[0])
			if vOutput, err := exec.Command("dpkg", "-s", packageName).Output(); err == nil {
				lines := strings.Split(string(vOutput), "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "Version:") {
						info.Version = strings.TrimSpace(strings.TrimPrefix(line, "Version:"))
						break
					}
				}
			}
		}
	}

	// rpm으로 버전 정보 확인 시도
	if info.Version == "" {
		rpmCmd := exec.Command("rpm", "-qf", filePath)
		if output, err := rpmCmd.Output(); err == nil {
			info.Version = strings.TrimSpace(string(output))
		}
	}

	// readelf로 ELF 헤더 정보 읽기
	readelfCmd := exec.Command("readelf", "-n", "-d", filePath)
	output, err = readelfCmd.Output()
	if err != nil {
		return info, fmt.Errorf("readelf command failed: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	var deps []string
	for _, line := range lines {
		switch {
		case strings.Contains(line, "Build ID:"):
			info.BuildID = strings.TrimSpace(strings.Split(line, "Build ID:")[1])
		case strings.Contains(line, "SONAME"):
			parts := strings.Split(line, "[")
			if len(parts) > 1 {
				info.SONAME = strings.Trim(parts[1], "[]")
			}
		case strings.Contains(line, "NEEDED"):
			parts := strings.Split(line, "[")
			if len(parts) > 1 {
				dep := strings.Trim(parts[1], "[]")
				deps = append(deps, dep)
			}
		}
	}
	info.Dependencies = deps

	return info, nil
}

// ScanSystem scans the specified paths for ELF files
func ScanSystem(paths []string) ([]FileVersionInfo, error) {
	var products []FileVersionInfo
	for _, rootDir := range paths {
		filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}

			elfInfo, err := GetELFInfo(path)
			if err != nil {
				return nil
			}
			products = append(products, elfInfo)
			return nil
		})
	}
	return products, nil
}
