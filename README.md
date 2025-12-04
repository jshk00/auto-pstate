# auto-pstate
**auto-pstate** is a lightweight tool for automatically switching AMD EPP P-state profiles based on system events such as charging. It helps optimize performance and power consumption seamlessly.

## Features
- Automatically switches AMD P-states based on charging events.
- Sets `balance_performance` profile when charging and `power` profile when on battery.
- Provides manual override and profile management via `pdctl`.

## Prerequisites
- AMD processor with `amd-pstate` support.
- Kernel parameter `amd_pstate=active` must be enabled.

## Installation

### One-line Installation
```sh
curl -sSL https://github.com/jshk00/auto-pstate/releases/download/0.0.2/install | sudo bash
```

### Build from source
- Ensure you have Go 1.22+ installed.
```sh
git clone https://github.com/jshk00/auto-pstate.git
cd auto-pstate
sudo ./install.sh
```

## Uninstall
```sh
sudo systemctl disable --now auto-pstate.service
sudo rm -rf /run/pstated
sudo rm -rf /usr/bin/pstated /usr/bin/pdctl
```

## Usage
The auto-pstate systemd service automatically applies profiles based on charging events:
- Charging: balance_performance profile is set
- On battery: power profile is set

## Commands
- List all available profiles
```sh
sudo pdctl list-prefs
```

- Set a profile manually (switches pstated to manual mode):
```sh
sudo pdctl set <profile_name>
```

- Re-enable automatic mode:
```sh
sudo pdctl auto on
```

- Disable automatic mode:
```sh
sudo pdctl auto off
```

- Fetch current status:
```sh
sudo pdctl status
```

## Examples
```sh
# List all profiles
$ sudo pdctl list-prefs
balance_performance
power
performance
default

# Set manual profile
$ sudo pdctl set balance_performance
p-state set to balance_performance(manual mode)

# Re-enable auto mode
$ sudo pdctl auto on
auto mode enabled

# Fetch current status
$ sudo pdctl status 
profile --> balance_performance
mode --> manual
```
