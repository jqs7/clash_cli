<h1 align="center">
  <img src="https://github.com/Dreamacro/clash/raw/master/docs/logo.png" alt="Clash" width="200">
  <br>
  Clash_CLI
  <br>
</h1>

<h4 align="center">A terminal client for <a href="https://github.com/Dreamacro/clash">Clash</a>.</h4>

![clash_cli](https://user-images.githubusercontent.com/12208686/54497096-91367080-4931-11e9-9851-93b09b91c161.gif)

## Install

You can build from source:

```sh
go get -u -v github.com/jqs7/clash_cli
```

Pre-built binaries are available: [release](https://github.com/jqs7/clash_cli/releases)

systemd oneshot service:
```toml
[Unit]
Description=Clash CLI
After=online.target clash.service
Requires=clash.service

[Service]
Type=oneshot
ExecStart=/usr/bin/clash_cli -q
WorkingDirectory=/home/user/.config/clash_cli

[Install]
WantedBy=multi-user.target
```

## Usage

`clash_cli [http://localhost:9090]`
