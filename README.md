# dp-frontend-articles-controller
Controller for handling articles on ONS website

### Getting started

* Run `make debug`

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable         | Default                   | Description
| ---------------------------- | ------------------------- | -----------
| BIND_ADDR                    | :26500                    | The host and port to bind to
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s                        | The graceful shutdown timeout in seconds (`time.Duration` format)
| DEBUG                        | false                     | Enable debug mode
| API_ROUTER_URL               | http://localhost:23200/v1 | The URL of the [dp-api-router](https://github.com/ONSdigital/dp-api-router)
| SITE_DOMAIN                  | localhost                 |
| HEALTHCHECK_INTERVAL         | 30s                       | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s                       | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2021, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

