# fs
JunOS BGP Flow Spec routes provisioning tool

[![Go Report](https://img.shields.io/badge/go%20report-A%2B-blue?style=flat-square&color=00c9ff&labelColor=bec8d2)](https://goreportcard.com/report/github.com/horseinthesky/fs)
[![License: MIT](https://img.shields.io/badge/License-MIT-blueviolet.svg?style=flat-square)](https://opensource.org/licenses/MIT)

---

## Config
Tool config file shoud look somewhat like this:
```yaml
source: <url_to_download_flow_rules_yaml_file>

inventory: <url_to_netbox_api_device/virtual-machines>

creds:
  username: <your_device_username>
  password: <your_device_password>
  key: <path_to_your_public_ssh_key>

interval: 10
```
#### Options:
- __source__: something like `https://<url>/repos/<reponame>/raw/flows.yml?raw`
  You can use local FS routes file with `-f` flag.
  If Web source is used you must provide Netbox token with `-t` flag.
- __inventory__: Netbox URL similar to `https://<netbox_url>/api/virtualization/virtual-machines/?role=<myfavoriterole>`
- __creds__: username, password and/or ssh key to reach the devices.
  If both password and key are provided key is used
- __interval__: sleep time (seconds) between deploying rules to devices

**NOTE**: path to config file is passed with `-c` flag.

**NOTE2**: log filepath and log level are passed with `-l` and `-d` flag respectively.

## Data
Flows data must have the following format:
```yaml
flows:
  - name: BOT-2251-1
    destination: 103.3.62.64/32
    protocol:
      - tcp
    destination-port:
      - 14433
      - 14444
  - name: INC-960-1-accept
    destination: 84.201.174.174/32
    protocol:
      - udp
    destination-port:
      - 5055
    action: accept
  - name: SUPPORT-82690-1
    destination: 178.154.244.169/32
    protocol:
      - udp
    source-port:
      - 53
      - 389
      - 11211
  - name: SUPPORT-82690-2
    destination: 178.154.244.169/32
    protocol:
      - icmp
  - name: DUTY-9631-1
    destination: 84.252.135.75/32
    protocol:
      - udp
    source-port:
      - 389
    action: discard
  - name: REGULAR-94800-1
    destination: 84.201.181.26/32
    protocol:
      - udp
    action: discard
  - name: REGULAR-94800-2
    destination: 84.201.171.239/32
    protocol:
      - udp
    action: discard
```
#### Supported options:
- __name__: rule name
- __destination__: destination prefix for this traffic flow
- __source__: source prefix for this traffic flow
- __protocol__: `tcp`/`udp`/`icmp`
- __destination-port__: destination TCP/UDP port
- __source-port__: source TCP/UDP port
- __action__: `accept`/`discard` (default)

YAML Flow Spec routes data is parsed to XML format to build a NETCONF payload and deployed to devices via NETCONF.

Tool supports JunOS devices only.
#### Result:
```bash
[edit routing-options flow]
root@vMX8# show
route BOT-2251-1 {
    match {
        destination-port [ 14433 14444 ];
        destination 103.3.62.64/32;
    }
    then discard;
}
```
