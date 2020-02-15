# STM32F303 based power measurement

This is a residential power use monitor based on the STM32f303 black pill. It supports one voltage channel and 3 current transformer channels.

Measurement circuit is a variation of https://learn.openenergymonitor.org/electricity-monitoring/ct-sensors/interface-with-arduino and https://learn.openenergymonitor.org/electricity-monitoring/voltage-sensing/measuring-voltage-with-an-acac-power-adapter

A client program exposes measurements as Prometheus metrics.

## Firmware

### Build

```
git submodule update --init # (Only needed once)
cd power-monitor
make -C libopencm3 # (Only needed once)
make
```

### Use

Flash the firmware and connect the STM32f303 with USB. Use the client to interface with the cdc-acm device.

## Client

### Build

```
cd pc_client
make
```

### Use

```
./powerclient -config=config.yaml
```

Config file contains dividers for different ADC channels to convert from ADC counts to current and voltage. See [pc_client/config_example.yaml](pc_client/config_example.yaml) for an example


## License

### Firmware

Since the  code is heavily based on libopencm3 examples it retains their LGPL3 license.

Uses nanopb `Copyright (c) 2011 Petteri Aimonen <jpa at nanopb.mail.kapsi.fi>`. See [power-monitor/nanopb/LICENCE.txt](power-monitor/nanopb/LICENSE.txt)

### Client

Licensed under the 3-Clause BSD License.
