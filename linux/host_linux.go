// +build linux

package linux

import (
	"github.com/shirou/gopsutil/utils"
	"github.com/shirou/gopsutil/structs"

	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

type LSB struct {
	ID          string
	Release     string
	Codename    string
	Description string
}

func HostInfo() (*structs.HostInfoStat, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	ret := &structs.HostInfoStat{
		Hostname: hostname,
		OS:       runtime.GOOS,
	}

	platform, family, version, err := getPlatformInformation()
	if err == nil {
		ret.Platform = platform
		ret.PlatformFamily = family
		ret.PlatformVersion = version
	}
	uptime, err := BootTime()
	if err == nil {
		ret.Uptime = uptime
	}

	return ret, nil
}

func BootTime() (uint64, error) {
	sysinfo := &syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(sysinfo); err != nil {
		return 0, err
	}
	return uint64(sysinfo.Uptime), nil
}

func Users() ([]structs.UserStat, error) {
	utmpfile := "/var/run/utmp"

	file, err := os.Open(utmpfile)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	u := utmp{}
	entrySize := int(unsafe.Sizeof(u))
	count := len(buf) / entrySize

	ret := make([]structs.UserStat, 0, count)

	for i := 0; i < count; i++ {
		b := buf[i*entrySize : i*entrySize+entrySize]

		var u utmp
		br := bytes.NewReader(b)
		err := binary.Read(br, binary.LittleEndian, &u)
		if err != nil {
			continue
		}
		user := structs.UserStat{
			User:     utils.ByteToString(u.UtUser[:]),
			Terminal: utils.ByteToString(u.UtLine[:]),
			Host:     utils.ByteToString(u.UtHost[:]),
			Started:  int(u.UtTv.TvSec),
		}
		ret = append(ret, user)
	}

	return ret, nil

}

func getLSB() (*LSB, error) {
	ret := &LSB{}
	if utils.PathExists("/etc/lsb-release") {
		contents, err := utils.ReadLines("/etc/lsb-release")
		if err != nil {
			return ret, err // return empty
		}
		for _, line := range contents {
			field := strings.Split(line, "=")
			if len(field) < 2 {
				continue
			}
			switch field[0] {
			case "DISTRIB_ID":
				ret.ID = field[1]
			case "DISTRIB_RELEASE":
				ret.Release = field[1]
			case "DISTRIB_CODENAME":
				ret.Codename = field[1]
			case "DISTRIB_DESCRIPTION":
				ret.Description = field[1]
			}
		}
	} else if utils.PathExists("/usr/bin/lsb_release") {
		out, err := exec.Command("/usr/bin/lsb_release").Output()
		if err != nil {
			return ret, err
		}
		for _, line := range strings.Split(string(out), "\n") {
			field := strings.Split(line, ":")
			if len(field) < 2 {
				continue
			}
			switch field[0] {
			case "Distributor ID":
				ret.ID = field[1]
			case "Release":
				ret.Release = field[1]
			case "Codename":
				ret.Codename = field[1]
			case "Description":
				ret.Description = field[1]
			}
		}

	}

	return ret, nil
}

func getPlatformInformation() (string, string, string, error) {
	platform := ""
	family := ""
	version := ""

	lsb, _ := getLSB()

	if utils.PathExists("/etc/oracle-release") {
		platform = "oracle"
		contents, err := utils.ReadLines("/etc/oracle-release")
		if err == nil {
			version, _ = getRedhatishVersion(contents)
		}
	} else if utils.PathExists("/etc/enterprise-release") {
		platform = "oracle"
		contents, err := utils.ReadLines("/etc/enterprise-release")
		if err == nil {
			version, _ = getRedhatishVersion(contents)
		}
	} else if utils.PathExists("/etc/debian_version") {
		if lsb.ID == "Ubuntu" {
			platform = "ubuntu"
			version = lsb.Release
		} else if lsb.ID == "LinuxMint" {
			platform = "linuxmint"
			version = lsb.Release
		} else {
			if utils.PathExists("/usr/bin/raspi-config") {
				platform = "raspbian"
			} else {
				platform = "debian"
			}
			contents, err := utils.ReadLines("/etc/debian_version")
			if err == nil {
				version = contents[0]
			}
		}
	} else if utils.PathExists("/etc/redhat-release") {
		contents, err := utils.ReadLines("/etc/redhat-release")
		if err == nil {
			version, _ = getRedhatishVersion(contents)
			platform, _ = getRedhatishPlatform(contents)
		}
	} else if utils.PathExists("/etc/system-release") {
		contents, err := utils.ReadLines("/etc/system-release")
		if err == nil {
			version, _ = getRedhatishVersion(contents)
			platform, _ = getRedhatishPlatform(contents)
		}
	} else if utils.PathExists("/etc/gentoo-release") {
		platform = "gentoo"
		contents, err := utils.ReadLines("/etc/gentoo-release")
		if err == nil {
			version, _ = getRedhatishVersion(contents)
		}
		// TODO: suse detection
		// TODO: slackware detecion
	} else if utils.PathExists("/etc/arch-release") {
		platform = "arch"
		// TODO: exherbo detection
	} else if lsb.ID == "RedHat" {
		platform = "redhat"
		version = lsb.Release
	} else if lsb.ID == "Amazon" {
		platform = "amazon"
		version = lsb.Release
	} else if lsb.ID == "ScientificSL" {
		platform = "scientific"
		version = lsb.Release
	} else if lsb.ID == "XenServer" {
		platform = "xenserver"
		version = lsb.Release
	} else if lsb.ID != "" {
		platform = strings.ToLower(lsb.ID)
		version = lsb.Release
	}

	switch platform {
	case "debian", "ubuntu", "linuxmint", "raspbian":
		family = "debian"
	case "fedora":
		family = "fedora"
	case "oracle", "centos", "redhat", "scientific", "enterpriseenterprise", "amazon", "xenserver", "cloudlinux", "ibm_powerkvm":
		family = "rhel"
	case "suse":
		family = "suse"
	case "gentoo":
		family = "gentoo"
	case "slackware":
		family = "slackware"
	case "arch":
		family = "arch"
	case "exherbo":
		family = "exherbo"
	}

	return platform, family, version, nil

}

func getRedhatishVersion(contents []string) (string, error) {
	return "", nil
}

func getRedhatishPlatform(contents []string) (string, error) {
	return "", nil
}
