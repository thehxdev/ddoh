# ddoh
ddoh is a simple and easy to use DNS-over-HTTPS client that acts like a normal dns server on `127.0.0.1:53` and
sends DNS requests to specified DoH server.


## Build

### Linux / macOS
```bash
CGO_ENABLED=0 go build -ldflags='-s -w -buildid=' .
```

### Windows
```powershell
$env:CGO_ENABLED=0
go build -ldflags='-s -w -buildid=' .
```

### Makefile
```bash
make
```


## Usage
Before starting ddoh, make sure that port 53 is open and not used by another process. Then:

> [!WARNING]
> Starting a server on port 53 needs superuser access. You may need to start `ddoh` with `sudo` or as root user.

```bash
./ddoh -c config.json
```

## Config
Config file is in JSON format.
```json
{
    "local_resolver": "9.9.9.9",
    "doh_server": "https://max.rethinkdns.com/rec",
    "doh_ip": "137.66.7.89",
    "udp_buffer_size": 512
}
```

- `local_resolver`: Used to resolve the DoH server hostname. Must be an IP address.
- `doh_server`: DoH URL
- `doh_ip` (optional): Sometimes DoH servers are limited due to censorship. You can specify DoH hostname's IP address.
If IP address is sprecified, ddoh will not use `local_resolver`.
- `udp_buffer_size`: UDP buffer size. Higher values will increase memory usage. (default value `512`)


## Contribution
If you can improve the source code or make this software better, feel free to send a PR :)
