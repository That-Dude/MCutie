package getstats

import (
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"strings"

	"github.com/distatus/battery"
	externalip "github.com/glendc/go-external-ip"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func dealwithErr(err error) {
	if err != nil {
		fmt.Println(err)
		//os.Exit(-1)
	}
}

func bytesToSize(bytes uint64) string {
	sizes := []string{"", "", "", "", ""}
	if bytes == 0 {
		return fmt.Sprint(float64(0), "bytes")
	} else {
		var bytes1 = float64(bytes)
		var i = math.Floor(math.Log(bytes1) / math.Log(1024))
		var count = math.Round(bytes1 / math.Pow(1024, i))
		var j = int(i)
		return fmt.Sprint(float64(count), sizes[j])
	}
}

// TimeStr - Convert seconds to day,weeks....etc
func TimeStr(sec int) (res string) {
	wks, sec := sec/604800, sec%604800
	ds, sec := sec/86400, sec%86400
	hrs, sec := sec/3600, sec%3600
	mins, sec := sec/60, sec%60
	CommaRequired := false
	if wks != 0 {
		res += fmt.Sprintf("%dw", wks)
		CommaRequired = true
	}
	if ds != 0 {
		if CommaRequired {
			res += " "
		}
		res += fmt.Sprintf("%dd", ds)
		CommaRequired = true
	}
	if hrs != 0 {
		if CommaRequired {
			res += " "
		}
		res += fmt.Sprintf("%dh", hrs)
		CommaRequired = true
	}
	if mins != 0 {
		if CommaRequired {
			res += " "
		}
		res += fmt.Sprintf("%dm", mins)
		CommaRequired = true
	}
	if sec != 0 {
		if CommaRequired {
			res += " "
		}
		res += fmt.Sprintf("%ds", sec)
	}
	return
}

// BatteryLevel x
func BatteryLevel() string {
	s := "none"
	batteries, err := battery.GetAll()
	if err, isFatal := err.(battery.ErrFatal); isFatal {
		fmt.Fprintln(os.Stderr, err)
	}
	if len(batteries) == 0 {
	}
	errs, partialErrs := err.(battery.Errors)
	for i, bat := range batteries {
		if partialErrs && errs[i] != nil {
			fmt.Fprintf(os.Stderr, "Error getting info for BAT%d: %s\n", i, errs[i])
			continue
		}
		s = fmt.Sprintf("%.0f", bat.Current/bat.Full*100)
	}
	return s
}

// ExternalIP x
func ExternalIP() string {
	s := "x"
	consensus := externalip.DefaultConsensus(nil, nil)

	ip, err := consensus.ExternalIP()
	if err == nil {
		s = (ip.String())
	}
	return s
}

// CurrentUser x
func CurrentUser() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	return user.Username
}

// MemTotal x
func MemTotal() string {
	vmStat, err := mem.VirtualMemory()
	dealwithErr(err)
	s := fmt.Sprintf(bytesToSize(vmStat.Total))
	return s
}

// MemUsed x
func MemUsed() string {
	vmStat, err := mem.VirtualMemory()
	dealwithErr(err)
	s := fmt.Sprintf(bytesToSize(vmStat.Used))
	return s
}

// MemUsedPercent x
func MemUsedPercent() string {
	vmStat, err := mem.VirtualMemory()
	dealwithErr(err)
	s := fmt.Sprintf(strconv.FormatFloat(vmStat.UsedPercent, 'f', 0, 64))
	return s
}

// MemFree x
func MemFree() string {
	vmStat, err := mem.VirtualMemory()
	dealwithErr(err)
	s := fmt.Sprintf(bytesToSize(vmStat.Free))
	return s
}

// DiskTotal x
func DiskTotal() string {
	s := "0"
	if runtime.GOOS == "darwin" {
		diskStat, err := disk.Usage("/System/Volumes/Data") //Macos point point
		dealwithErr(err)
		s = fmt.Sprintf(bytesToSize(diskStat.Total))
	}

	if runtime.GOOS == "windows" {
		diskStat, err := disk.Usage("\\") // mount point for Windows
		dealwithErr(err)
		s = fmt.Sprintf(bytesToSize(diskStat.Total))
	}
	return s
}

// DiskUsed x
func DiskUsed() string {
	s := "0"
	if runtime.GOOS == "darwin" {
		diskStat, err := disk.Usage("/System/Volumes/Data") //Macos point point
		dealwithErr(err)
		s = fmt.Sprintf(bytesToSize(diskStat.Used))
	}

	if runtime.GOOS == "windows" {
		diskStat, err := disk.Usage("\\") // mount point for Windows
		dealwithErr(err)
		s = fmt.Sprintf(bytesToSize(diskStat.Used))
	}
	return s
}

// DiskFree x
func DiskFree() string {
	s := "0"
	if runtime.GOOS == "darwin" {
		diskStat, err := disk.Usage("/System/Volumes/Data") //Macos point point
		dealwithErr(err)
		s = fmt.Sprintf(bytesToSize(diskStat.Free))
	}

	if runtime.GOOS == "windows" {
		diskStat, err := disk.Usage("\\") // mount point for Windows
		dealwithErr(err)
		s = fmt.Sprintf(bytesToSize(diskStat.Free))
	}
	return s
}

// UpTime x
func UpTime() string {
	hostStat, err := host.Info()
	dealwithErr(err)
	x := strconv.FormatUint(hostStat.Uptime, 10)
	i, err := strconv.Atoi(x)
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	s := fmt.Sprintf(TimeStr(i))
	return s
}

// LocalIP x
func LocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// localAddr := conn.LocalAddr().(*net.UDPAddr)
	temp := conn.LocalAddr().String() //this is the local IP in the format 1.2.3.4:55555
	delimiter := ":"
	strlocalAddr := strings.Split(temp, delimiter)[0] //remove everything after and including the colon
	s := strlocalAddr
	return s
}

// CPUUsage x
func CPUUsage() string {
	percentage, err := cpu.Percent(0, false)
	dealwithErr(err)
	s := ""
	for idx, cpupercent := range percentage {
		s = ("Current CPU utilization: [" + strconv.Itoa(idx) + "] " + strconv.FormatFloat(cpupercent, 'f', 2, 64))
		s = fmt.Sprintf(strconv.FormatFloat(cpupercent, 'f', 0, 64))

	}
	return s
}

// HostName x
func HostName() string {
	hostStat, err := host.Info()
	dealwithErr(err)
	//fmt.Println("HostnameA: " + hostStat.Hostname)
	s := hostStat.Hostname
	return s
}
