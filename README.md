![](shine.png)
# Shine Online Packet Sniffer

> Tool that captures all packets between the client and services and show them in real time.
> The program also saves the streams between the flows to a pcapng file.
>
 [![Go Report Card](https://goreportcard.com/badge/github.com/shine-o/shine.engine.packet-sniffer)](https://goreportcard.com/report/github.com/shine-o/shine.engine.packet-sniffer)
 
#### The road so far: [video showcase](https://www.youtube.com/watch?v=Y08oHJucHRI)


#### Configuration

Adjust the **config/.sniffer.yml** file to your needs


## build

 ``` 
 $ go build -o sniffer.exe
 ```

## sniffer capture

Start capturing and decoding packets

### Synopsis

Start capturing and decoding packets

```
e.g, server side capture:

$ sniffer capture --config "config/.sniffer.server.default.yml"
```

### Options

```
  -h, --help   help for capture
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.sniffer.yaml)
```


#### Packet info


![](packet-flow-draw.png)