<img src="https://github.com/user-attachments/assets/0e70decd-79c3-4a4e-bc69-824ad51bfb56">
<div>
  <img src="https://img.shields.io/badge/go-000000?style=for-the-badge&logo=go"> 
  <img src="https://img.shields.io/badge/linux-000000?style=for-the-badge&logo=linux">
  <img src="https://img.shields.io/badge/windows-000000?style=for-the-badge&logo=windows">
</div>
<b>This is basically tripwire, but for your Linux server. The idea is that you set up Syscanary to monitor changes in things that would indicate a system compromise or a malfunction.</b>

## Features
Syscanary can currently monitor changes in:
- File system
- USB devices
- Open ports (currently only Linux support)

## Installation
```
sudo apt update && sudo apt install golang  
go install github.com/spoofimei/syscanary
export PATH=$PATH:~/go/bin
```

## Usage
`syscanary`

## Configuration
Open syscanary.json and make your configurations:
```
{
    "loglevel": 1, # <-- 0=debug 1=info 2=error
    "logfile": "alerts.log", # <-- remove if you want to have console output 
    "detections": ["usb", "integrity", "ports", "internet"], # <-- remove or add modules to enable or disable them
    "modules": { # <-- DO NOT REMOVE ANY MODULES FROM HERE AND DON'T LEAVE SETTINGS EMPTY!!!
        "integrity": {
            "interval": 1, # <-- how many seconds it will wait before checking again
            "paths": ["/var/log"]
        },
        "usb": {
            "interval": 1
        },
        "ports": {
            "interval": 1,
            "ignorelocal": true
        },
        "internet": {
            "interval": 1,
            "domain": "example.com"
        }
    }
}
```
