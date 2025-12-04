package internal

import (
	"bytes"
	"log"
	"sync/atomic"
	"syscall"
)

type AutoEPPSetter struct {
	mode    atomic.Value
	eventCh chan bool
}

func (as *AutoEPPSetter) Start() {
	as.eventCh = make(chan bool)
	as.SetMode(auto)
	go as.chargeEvent()
	go as.run()
}

func (as *AutoEPPSetter) run() {
	for e := range as.eventCh {
		if as.GetMode() == auto {
			if e {
				log.Println("[INFO] setting epp state to balance_performance")
				if err := SetEPP(defaultEppStateAC); err != nil {
					log.Println(err)
				}
				continue
			}
			log.Println("[INFO] setting epp state to power")
			if err := SetEPP(defaultEppStateBat); err != nil {
				log.Println(err)
			}
		}
	}
	log.Println("closed auto epp setter")
}

func (as *AutoEPPSetter) chargeEvent() {
	fd, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW, syscall.NETLINK_KOBJECT_UEVENT)
	if err != nil {
		log.Println(err)
		return
	}
	defer syscall.Close(fd)

	if err := syscall.Bind(fd, &syscall.SockaddrNetlink{
		Family: syscall.AF_NETLINK,
		Groups: 1,
	}); err != nil {
		log.Println(err)
		return
	}

	buf := make([]byte, 4096)
	for {
		n, _, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			log.Println(err)
			continue
		}
		parts := bytes.Split(buf[:n], []byte{0})
		uevent := make(map[string]string)

		for _, part := range parts {
			kv := bytes.SplitN(part, []byte{'='}, 2)
			if len(kv) == 2 {
				uevent[string(kv[0])] = string(kv[1])
			}
		}
		if uevent["DEVTYPE"] == "power_supply" &&
			uevent["ACTION"] == "change" &&
			uevent["POWER_SUPPLY_NAME"] == "AC" {
			as.eventCh <- uevent["POWER_SUPPLY_ONLINE"] == "1"
		}
	}
}

func (as *AutoEPPSetter) Close() {
	close(as.eventCh)
}

type mode int

const (
	auto   mode = 0
	manual mode = 1
)

func (m mode) String() string {
	return [...]string{"auto", "manual"}[m]
}

func (as *AutoEPPSetter) SetMode(m mode) {
	as.mode.Store(m)
}

func (as *AutoEPPSetter) GetMode() mode { //nolint
	return as.mode.Load().(mode) //nolint
}
