package internal

import (
	"log"
	"sync/atomic"
	"time"
)

type AutoEPPSetter struct {
	mode    atomic.Value
	profile string
	stop    bool
}

func (as *AutoEPPSetter) Run() {
	for !as.stop {
		if as.GetMode() == auto {
			c, err := Charging()
			if err != nil {
				log.Println(err)
				continue
			}
			if c && as.profile != DefaultEppStateAC {
				log.Println("[INFO] setting epp state to balance_performance")
				as.profile = DefaultEppStateAC
				if err := SetEPP(DefaultEppStateAC); err != nil {
					log.Println(err)
				}
			}
			if !c && as.profile != DefaultEppStateBat {
				log.Println("[INFO] setting epp state to power")
				as.profile = DefaultEppStateBat
				if err := SetEPP(DefaultEppStateBat); err != nil {
					log.Println(err)
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func (as *AutoEPPSetter) Close() {
	as.stop = true
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
