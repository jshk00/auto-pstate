package internal

const (
	defaultEppStateAC  = "balance_performance"
	defaultEppStateBat = "power"
	defaultGovernor    = "powersave"

	scalingDriverPath = "/sys/devices/system/cpu/cpu0/cpufreq/scaling_driver"
	governerPath      = "/sys/devices/system/cpu/cpu%d/cpufreq/scaling_governor"
	eppPath           = "/sys/devices/system/cpu/cpu%d/cpufreq/energy_performance_preference"
	preferences       = "/sys/devices/system/cpu/cpu0/cpufreq/energy_performance_available_preferences"
	SockPath          = "/run/pstated/pstated.sock"
)
