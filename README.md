[![CI](https://github.com/mdb/gossboss/actions/workflows/ci.yml/badge.svg)](https://github.com/mdb/gossboss/actions/workflows/ci.yml)

# gossboss

Collect and view aggregated [goss](https://goss.rocks/) test results from
multiple remote goss servers.

## Usage

`gossboss` can be used as a CLI or as a server.

### CLI

`gossboss healthzs` collects `goss` test results from multiple `goss` servers
and reports their results:

```
gossboss healthzs \
  --servers "http://foo.com/healthz" \
  --servers "http://bar.com/healthz"
 ✘ http://foo.com/healthz
 ✘ http://bar.com/healthz
Error: Goss test failed
```

### Server

`gossboss serve` starts a server whose `/healthzs` endpoint returns aggregated
`goss` test results from multiple `goss` servers:

```
gossboss serve \
  --servers "http://foo.com/healthz" \
  --servers "http://bar.com/healthz"
2021/03/28 06:22:10 Starting server on :8085
```

View the `goss` test results from each server in aggregate via `gossboss`'s
`/healthzs` endpoint:

```
$ curl localhost:8085/healthzs | jq

{
  "Healthzs": [
    {
      "Result": {
        "results": [
          {
            "successful": false,
            "resource-id": "tcp://google.com:443",
            "resource-type": "Addr",
            "title": "",
            "meta": null,
            "test-type": 0,
            "result": 1,
            "property": "reachable",
            "err": null,
            "expected": [
              "true"
            ],
            "found": [
              "false"
            ],
            "human": "Expected\n    <bool>: false\nto equal\n    <bool>: true",
            "duration": 503511341,
            "summary-line": "Addr: tcp://google.com:443: reachable:\nExpected\n    <bool>: false\nto equal\n    <bool>: true"
          }
        ],
        "summary": {
          "test-count": 1,
          "failed-count": 1,
          "total-duration": 503801507
        },
        "summary-line": ""
      },
      "URL": "http://foo.com/healthz",
      "Error": null
    },
    {
      "Result": {
        "results": [
          {
            "successful": false,
            "resource-id": "tcp://google.com:443",
            "resource-type": "Addr",
            "title": "",
            "meta": null,
            "test-type": 0,
            "result": 1,
            "property": "reachable",
            "err": null,
            "expected": [
              "true"
            ],
            "found": [
              "false"
            ],
            "human": "Expected\n    <bool>: false\nto equal\n    <bool>: true",
            "duration": 503326591,
            "summary-line": "Addr: tcp://google.com:443: reachable:\nExpected\n    <bool>: false\nto equal\n    <bool>: true"
          }
        ],
        "summary": {
          "test-count": 1,
          "failed-count": 1,
          "total-duration": 503626800
        },
        "summary-line": ""
      },
      "URL": "http://bar.com/healthz",
      "Error": null
    }
  ],
  "Summary": {
    "failed-count": 2,
    "errored-count": 0
  }
}
```

## Development

To test and build `gossboss`:

```
make
```

## Releasing

`make tag` publishes a new git tag corresponding to the `Makefile`'s `VERSION`
variable. In turn, [GitHub Actions](https://github.com/mdb/gossboss/actions)
builds and publishes a new, corresponding
[release](http://github.com/mdb/gossboss/releases) and [Docker
image](https://hub.docker.com/r/clapclapexcitement/gossboss) using
[goreleaser](https://goreleaser.com/).

## TODO candidates

* service discovery (could `gossboss` be extensible to support the discovery of
goss server URLs via a cloud provider API?)
* failure notifications (Slack,
PagerDuty, SMS, webhook, etc.)
* configuration file support (could `gossboss` be
configured via a TOML or YAML file rather than command line flags?)
