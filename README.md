# srv-announcer

## Building
This is a fairly straightforward Go application.

 - Obtain a recent Go version.
   Look at `.github/workflows/ci.yml` for the one used in CI.
 - run `go test ./...` for the tests
 - run `go build` to build a `./srv-announcer` static binary

## Usage
```
NAME:
   srv-announcer - Sidecar managing DNS records in an SRV record set (RFC2782), a poormans alternative to proper service discovery

USAGE:
   srv-announcer [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dry-run                    Don't actually update DNS, only log what would be done (default: false) [$SRV_ANNOUNCER_DRY_RUN]
   --zone-name value            Name of the Route53 Zone the records to manage are in [$SRV_ANNOUNCER_ZONE_NAME]
   --srv-record-name value      RFC2782 Name (_service._proto.name) [$SRV_ANNOUNCER_SRV_RECORD_NAME]
   --srv-record-ttl value       TTL of the RFC2782 SRV Record Set in seconds (default: 60) [$SRV_ANNOUNCER_SRV_RECORD_TTL]
   --srv-record-priority value  Priority of the RFC2782 SRV Record (default: 10) [$SRV_ANNOUNCER_SRV_RECORD_PRIORITY]
   --srv-record-weight value    Weight of the RFC2782 SRV Record (default: 10) [$SRV_ANNOUNCER_SRV_RECORD_WEIGHT]
   --srv-record-port value      Port of the RFC2782 SRV Record (default: 443) [$SRV_ANNOUNCER_SRV_RECORD_PORT]
   --srv-record-target value    Target of the RFC2782 SRV Record. Usualy something resembling your hostname. [$SRV_ANNOUNCER_SRV_RECORD_TARGET]
   --check-target value         hostname:port to check. Will be $srv-record-target:$srv-record-port if unspecified [$SRV_ANNOUNCER_CHECK_TARGET]
   --check-interval value       Interval between checks (default: 10s) [$SRV_ANNOUNCER_CHECK_INTERVAL]
   --check-timeout value        Timeout for each check (default: 1s) [$SRV_ANNOUNCER_CHECK_TIMEOUT]
   --help, -h                   show help (default: false)
```
