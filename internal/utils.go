package internal

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// setGovernor ensures all cores use the default CPU frequency governor.
func SetGovernor() error {
	b, err := os.ReadFile(fmt.Sprintf(GovernerPath, 0))
	if err != nil {
		return errors.New("governor file does not exist")
	}
	if strings.TrimSpace(string(b)) != DefaultGovernor {
		for i := 0; i < runtime.NumCPU(); i++ {
			err := os.WriteFile(
				fmt.Sprintf(GovernerPath, i),
				[]byte(DefaultGovernor),
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
func SetEPP(val string) error {
	for i := 0; i < runtime.NumCPU(); i++ {
		err := os.WriteFile(fmt.Sprintf(EppPath, i), []byte(val), os.ModePerm)
		if err != nil {
			return fmt.Errorf("err: %w, while setting epp_value %s on core %d", err, val, i)
		}
	}
	return nil
}

func GetEPP() (string, error) {
	b, err := os.ReadFile(fmt.Sprintf(EppPath, 0))
	if err != nil {
		return "", fmt.Errorf("unable to get current_profile: %w", err)
	}
	return string(b), nil
}

// isRoot returns an error if the program is not running as root.
func IsRoot() error {
	if os.Geteuid() != 0 {
		return errors.New("script must be run with root")
	}
	return nil
}

// isPState checks whether the amd-pstate-epp driver is active.
func IsPState() error {
	b, err := os.ReadFile(ScalingDriverPath)
	if err != nil {
		return fmt.Errorf("file does not exist for scaling driver: %w", err)
	}
	if strings.TrimSpace(string(b)) != "amd-pstate-epp" {
		return errors.New("system is not running amd-pstate-epp")
	}
	return nil
}

// charging returns true if AC adapter is currently online.
func Charging() (bool, error) {
	b, err := os.ReadFile("/sys/class/power_supply/AC/online")
	if err != nil {
		return false, fmt.Errorf("file not found for AC: %w", err)
	}
	return strings.TrimSpace(string(b)) == "1", nil
}

// List the available preference in epp profiles
func GetPreferences() ([]string, error) {
	b, err := os.ReadFile(Preferences)
	if err != nil {
		return nil, fmt.Errorf("failed to list preferences: %w", err)
	}
	return strings.Fields(strings.TrimSpace(string(b))), nil
}
