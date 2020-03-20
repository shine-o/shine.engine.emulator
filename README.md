![](shine.png)
# Shine Online Packet Sniffer

> Tool that captures all packets between the client and services and show them in real time.
> The program also saves the streams between the flows to a pcapng file.
>
 [![Go Report Card](https://goreportcard.com/badge/github.com/shine-o/shine.engine.packet-sniffer)](https://goreportcard.com/report/github.com/shine-o/shine.engine.packet-sniffer)
#### Configuration

Adjust the **config.yml** file to your needs

#### Default ports for 2016 services. 

| Service     | ALL  | Client | Zone | w_Manager |
| ----------- | ---- | ------ | ---- | --------- |
| Login       |      | 9010   |      | 9011      |
| w_Manager_0 |      | 9110   | 9111 |           |
| Zone_0      |      | 9210   | 9218 |           |
| Zone_1      |      | 9212   | 9219 |           |
| Zone_2      |      | 9214   | 9220 |           |
| Zone_3      |      | 9216   | 9221 |           |
| Zone_4      |      | 9218   | 9222 |           |
| AccountLog  | 9311 |        |      |           |
| Character   | 9411 |        |      |           |
| GameLog     | 9511 |        |      |           |

#### Packet info

**NOTE**: current implementation doesn't represent the flow defined in the Stream graphic.

![](packet-flow-draw.png)
