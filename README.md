# Tapogo

![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/achetronic/tapogo)
![GitHub](https://img.shields.io/github/license/achetronic/tapogo)
[![Go Reference](https://pkg.go.dev/badge/github.com/achetronic/tapogo.svg)](https://pkg.go.dev/github.com/achetronic/tapogo)
[![Go Report Card](https://goreportcard.com/badge/github.com/achetronic/tapogo)](https://goreportcard.com/report/github.com/achetronic/tapogo)

![YouTube Channel Subscribers](https://img.shields.io/youtube/channel/subscribers/UCeSb3yfsPNNVr13YsYNvCAw?label=achetronic&link=http%3A%2F%2Fyoutube.com%2Fachetronic)
![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/achetronic?style=flat&logo=twitter&link=https%3A%2F%2Ftwitter.com%2Fachetronic)

A golang library (and CLI) to control your Tapo SmartPlugs `P100 / P110` with latest firmware versions

> At this moment, the library is not covering the whole API as it is discovered by reverse engineering
> If you want to cover more things, consider [contributing](#how-to-contribute)

## Motivation

I own several Tapo P100/P110 at home to control some appliances using custom automations 
[for example for the heater](https://github.com/achetronic/autoheater/). 
Honestly, I completely forgot to disable firmware updates, and before [Matter](https://csa-iot.org/all-solutions/matter/) 
disabling updates is almost a requirement as automation manufacturers tends to modify their closed APIs in a way they 
change almost everything from time to time, forcing you to do some reverse engineering.

As you can imagine, this project was created because of that, when some days ago, TPLink decided to switch from 
`securePassthrough` protocol to `KLAP` on `P100 / P110` plugs, and the library I use under the hood [does not seem
to be actively maintained](https://github.com/richardjennings/tapo/issues/4#issuecomment-1840902314), 
so I decided to research and craft a new one on my own.

> This library does not pretend to cover all devices, not even all the protocol versions, but always the latest ones.
> As the project will use releases, you can select which version fit your needs

## Library

```go
    import (
        "github.com/achetronic/tapogo/pkg/tapogo"
        "github.com/achetronic/tapogo/api/types"
    )

    var tapoClient *tapogo.Tapo
    var response   *types.ResponseSpec
    var err error

	tapoClient, err = tapogo.NewTapo("192.168.0.100", "username", "password")
    response, err = tapoClient.TurnOn()
    response, err = tapoClient.TurnOff()
    response, err = tapoClient.GetEnergyUsage()
    response, err = tapoClient.DeviceInfo()
```

## CLI

`go install github.com/achetronic/tapogo`

### Usage
```
tapogo <ip-address> <username> <password> [on, off, energy-usage, device-info]
```

For example:

```
tapogo 192.168.0.100 email@address thepassword energy-usage
{
  "error_code": 0,
  "result": {
    "current_power": 0,
    ...
    "month_energy": 10000,
    "month_runtime": 10000,
    "today_energy": 400,
    "today_runtime": 300
  }
}
```

## How to contribute

Of course, we are open to external collaborations for this project. For doing it you must:

* Open an issue to discuss what is needed and the reason
* Fork the repository
* Make your changes to the code
* Open a PR. The code will be reviewed and tested (always)

> We are developers and hate bad code. For that reason we ask you the highest quality on each line of code to improve
> this project on each iteration.

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

## Special mention

This project was done using IDEs from JetBrains. They helped us to develop faster, so we recommend them a lot! 🤓

<img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.png" alt="JetBrains Logo (Main) logo." width="150">
