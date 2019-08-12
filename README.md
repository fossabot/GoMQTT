# GoMQTT

[![Travis CI](https://img.shields.io/travis/com/Hexawolf/GoMQTT.svg?style=flat-square&logo=Linux)](https://travis-ci.com/Hexawolf/GoMQTT)
[![CodeCov](https://img.shields.io/codecov/c/github/Hexawolf/GoMQTT.svg?style=flat-square)](https://codecov.io/gh/Hexawolf/GoMQTT)
[![stability-experimental](https://img.shields.io/badge/stability-experimental-orange.svg?style=flat-square)](https://github.com/emersion/stability-badges#experimental)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FHexawolf%2FGoMQTT.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2FHexawolf%2FGoMQTT?ref=badge_shield)

GoMQTT is an experimental broker implementation for MQTT-SN protocol.
Any client that implements it properly can use this broker for sending and receiving messages.

**WARNING:** GoMQTT already plans to violate some parts of MQTT and MQTT-SN standard!

## Getting started

GoMQTT is mainly tested against Linux and must work on any common distribution that can run a Go
compiler. Building should be simple:

```bash
go build
```

After building, rename `broker.example.cfg` into `broker.cfg` and change values according to your needs.
Hopefully, now you are ready to run it.

## License

The code is under MIT license. See [LICENSE](LICENSE) for more information.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FHexawolf%2FGoMQTT.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2FHexawolf%2FGoMQTT?ref=badge_large)