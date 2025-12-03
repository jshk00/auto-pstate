package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

const (
	defaultEppStateAC  = "balance_performance"
	defaultEppStateBat = "power"
	defaultGovernor    = "powersave"

	scalingDriverPath = "/sys/devices/system/cpu/cpu0/cpufreq/scaling_driver"
	governerPath      = "/sys/devices/system/cpu/cpu%d/cpufreq/scaling_governor"
	eppPath           = "/sys/devices/system/cpu/cpu%d/cpufreq/energy_performance_preference"
	preferences       = "/sys/devices/system/cpu/cpu0/cpufreq/energy_performance_available_preferences"
)

// main validates environment, then executes either auto or manual mode.
func main() {
	log.SetFlags(0)
	auto := flag.Bool("auto", false, "start the pstate auto setter based on power")
	flag.Parse()

	if err := isRoot(); err != nil {
		log.Fatalf("[ERROR] %v\n", err)
	}
	if err := isPState(); err != nil {
		log.Fatalf("[ERROR] %v\n", err)
	}

	if *auto {
		if err := runAuto(); err != nil {
			log.Fatalf("[ERROR] %v\n", err)
		}
	}

	if err := runManual(); err != nil {
		log.Fatalf("[ERROR] %v\n", err)
	}
}

// runAuto sets the governor, applies initial EPP state,
// and listens for AC adapter change events to update EPP dynamically.
func runAuto() error {
	if err := setGovernor(); err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	return setAutoState(ctx)
}

// runManual lists available EPP profiles and applies the user selection.
func runManual() error {
	b, err := os.ReadFile(preferences)
	if err != nil {
		return fmt.Errorf("failed to list preferences: %w", err)
	}

	entries := strings.Fields(strings.TrimSpace(string(b)))
	m := make(map[int]string, len(entries))
	for i, p := range entries {
		m[i+1] = p
		fmt.Printf("%d] %s\n", i+1, p)
	}

	var choice int
	fmt.Print("select one from above: ")
	if _, err := fmt.Scanf("%d", &choice); err != nil {
		return errors.New("invalid input")
	}
	profile, ok := m[choice]
	if !ok {
		return errors.New("option is not available")
	}

	return setEPP(profile)
}

// setAutoState listens to kernel uevents for AC state changes
// and adjusts EPP between AC and battery profiles.
func setAutoState(ctx context.Context) error {
	run := true
	currentProfile := ""
	go func() { <-ctx.Done(); run = false }()
	for run {
		c, err := charging()
		if err != nil {
			return err
		}
		if c && currentProfile != defaultEppStateAC {
			log.Println("[INFO] setting epp state to balance_performance")
			currentProfile = defaultEppStateAC
			_ = setEPP(defaultEppStateAC)
		}
		if !c && currentProfile != defaultEppStateBat {
			log.Println("[INFO] setting epp state to power")
			currentProfile = defaultEppStateBat
			_ = setEPP(defaultEppStateBat)
		}
		time.Sleep(5 * time.Second)
	}
	return errors.New("closed")
}

// setGovernor ensures all cores use the default CPU frequency governor.
func setGovernor() error {
	b, err := os.ReadFile(fmt.Sprintf(governerPath, 0))
	if err != nil {
		return errors.New("governor file does not exist")
	}

	if strings.TrimSpace(string(b)) != defaultGovernor {
		for i := 0; i < runtime.NumCPU(); i++ {
			err := os.WriteFile(
				fmt.Sprintf(governerPath, i),
				[]byte(defaultGovernor),
				os.ModePerm,
			)
			if err != nil {
				return fmt.Errorf("err: %w, while setting governor on core %d", err, i)
			}
		}
	}
	return nil
}

// setEPP sets the EPP policy on all CPU cores.
func setEPP(val string) error {
	for i := 0; i < runtime.NumCPU(); i++ {
		err := os.WriteFile(fmt.Sprintf(eppPath, i), []byte(val), os.ModePerm)
		if err != nil {
			return fmt.Errorf("err: %w, while setting epp_value %s on core %d", err, val, i)
		}
	}
	return nil
}

// isRoot returns an error if the program is not running as root.
func isRoot() error {
	if os.Geteuid() != 0 {
		return errors.New("script must be run with root")
	}
	return nil
}

// isPState checks whether the amd-pstate-epp driver is active.
func isPState() error {
	b, err := os.ReadFile(scalingDriverPath)
	if err != nil {
		return fmt.Errorf("file does not exist for scaling driver: %w", err)
	}
	if strings.TrimSpace(string(b)) != "amd-pstate-epp" {
		return errors.New("system is not running amd-pstate-epp")
	}
	return nil
}

// charging returns true if AC adapter is currently online.
func charging() (bool, error) {
	b, err := os.ReadFile("/sys/class/power_supply/AC/online")
	if err != nil {
		return false, fmt.Errorf("file not found for AC: %w", err)
	}
	return strings.TrimSpace(string(b)) == "1", nil
}
