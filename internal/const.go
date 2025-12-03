package internal

const (
	DefaultEppStateAC  = "balance_performance"
	DefaultEppStateBat = "power"
	DefaultGovernor    = "powersave"

	ScalingDriverPath = "/sys/devices/system/cpu/cpu0/cpufreq/scaling_driver"
	GovernerPath      = "/sys/devices/system/cpu/cpu%d/cpufreq/scaling_governor"
	EppPath           = "/sys/devices/system/cpu/cpu%d/cpufreq/energy_performance_preference"
	Preferences       = "/sys/devices/system/cpu/cpu0/cpufreq/energy_performance_available_preferences"
)
