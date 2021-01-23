# gossboss

Collect aggregated [goss](https://goss.rocks/) test results from multiple remote goss servers.

## Usage

`gossboss` can be used as a CLI or as a server.

### CLI

```
gossboss healthzs \
  --server "http://foo.com/healthz" \
  --server "http://bar.com/healthz"
```

### Server - TODO

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
