# Luzifer / FHEMHue

`fhemhue` is a bridge application able to expose FHEM devices to be controlled using Amazon Alexa enabled devices like [Amazon Echo](http://amzn.to/2sd7li0).

The application needs to be deployed in the same network as the Alexa enabled device as it uses [UPNP](https://en.wikipedia.org/wiki/Universal_Plug_and_Play) for device discovery. Additionally it requires a telnet connection to be configured in your FHEM instance.

## Setup

1. Enable a telnet connection in your FHEM instance:
```
defmod telnetPort telnet 7072 global
```
2. Compile a `config.yml` file containing mappings of your FHEM devices to speakable names in the language you set your Alexa to - you can take a look into the `example` folder for my configuration. (For example German speaking Alexa is not able to switch on the "Dashboard" device because it does understand "Alexa, Sport on"...)
```yaml
---

telnet:
  ip: 127.0.0.1
  port: 7072

devices:
# - id: <fhem device id>
#   name: <speakable name in Alexas language>
#   states:
#     on: <fhem state to set: "set <id> <states.on>" (defaults to "on")>
#     off: <fhem state to set: "set <id> <states.off>" (defaults to "off")>
```
3. Inside the configuration set the IP and port to be used to connect to the `telnetPort` opened in step 1
4. Set up a start script / systemd service (example unit file in `example` folder)
5. Go to your [Alexa control panel](https://alexa.amazon.com/) into the "Smart Home" section and let Alexa discover your devices: You now should see your devices specified in the `config.yml` file.

----

![](https://d2o84fseuhwkxk.cloudfront.net/fhemhue.svg)
