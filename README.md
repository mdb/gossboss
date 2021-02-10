# gossboss

Collect aggregated [goss](https://goss.rocks/) test results from multiple remote goss servers.

## Usage

`gossboss` can be used as a CLI or as a server.

### CLI

Collect and view `goss` test results from multiple `goss` servers:

```
gossboss healthzs \
  --server "http://foo.com/healthz" \
  --server "http://bar.com/healthz"
```

### Server

Start a server whose `/healthzs` endpoint returns aggregated `goss` test results from multiple `goss` servers:

```
gossboss serve \
  --server "http://foo.com/healthz" \
  --server "http://bar.com/healthz"
```

## Development

To test and build `gossboss`:

```
make
```

## TODO candidates

* service discovery (could `gossboss` be extensible to support the discovery of goss server URLs via a cloud provider API?)
* failure notifications (Slack, PagerDuty, SMS, webhook, etc.)
* configuration file support (could `gossboss` be configured via a TOML or YAML file rather than command line flags?)
