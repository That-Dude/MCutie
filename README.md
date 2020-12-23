# MCutie

MCutie is a cross platform agent that publishes device metrics to your MQTT broker. It can also receive published commands to execute locally.

I've (lightly) tested it on macOS Big Sur and Windows 10 (x64) , it complies fine for Linux  but I've not had a chance to test it  yet.

Whilst it works fine for me, this should be considered <u>**Alpha**</u> software because:

- I'm a rookie coder! Golang is my first experience at coding in a 'proper' language

- I've not read up on how to test code yet!

- I basically use StackExchange to solve all of my problems :-)

## Motivation

I developed this agent specifically so that the computers in my local network could publish stats that show up in the HomeAssistant Dashboard:

![Home-Assistant-Dashboard-screen-shot.png](/Users/dan/Nextcloud/programming/golang/mcutie/docs/Home-Assistant-Dashboard-screen-shot.png)

You can also execute commands from Home Assistant using the Scripts module like this:

![Home-Assistant-Action-screenshot.png](/Users/dan/Nextcloud/programming/golang/mcutie/docs/Home-Assistant-Action-screenshot.png)

## Installation

Clone the repository or download the ZIP file somewhere.

Rename `config.yaml.sample` to `config.yaml`

Edit `config.yaml` with your MQTT broker settings

```yaml
url: ssl://ha.yourdomain.com:8883
username: mqttuser
password: Your-Passowrd
updateinterval: "9"
```

Update internal is how frequently you want the agent to publish your computer stats to the MQTT server. I guess in a larger network you might want to increate this to 60 seconds or more to prevent load on your MQTT broker.

### MacOS

I've writen a Bash script to install MCutie as a local user service:

```bash
chmod +x osx-install.sh

./osx-install.sh

```

It will now run at boot and continually attempt to re-start itself after 60 seconds after failure , e.g. you disconnect from your network or your MQTT broker goes offline.

### Windows

The program runs and works as expected but I've not decided how best to run it as a user land service yet.

## Notes

1. Did I mention this is alpha software under active development from a rookie coder? 

2. I decided early on that I didn't want MCutie to run as root, it's a big security burden that I didn't want to be responsible for and 99% of what I wanted to achieve can be done running in user space.

### Credits

https://stackexchange.com

[GitHub - shirou/gopsutil: psutil for golang](https://github.com/shirou/gopsutil)

[github.com/sirupsen/logrus]()

[github.com/eclipse/paho.mqtt.golang]()

[github.com/zpatrick/go-config]()

[github.com/deckarep/gosx-notifier]()