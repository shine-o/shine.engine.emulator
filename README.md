![](assets/shine.png)

# Shine Engine Emulator

---

[![Go Report Card](https://goreportcard.com/badge/github.com/shine-o/shine.engine.emulator)](https://goreportcard.com/report/github.com/shine-o/shine.engine.emulator)


Videos showcase: 

- [tools - packet sniffer](https://www.youtube.com/watch?v=Y08oHJucHRI)
- [world - character creation](https://www.youtube.com/watch?v=GF7cUkPe6BI&t=16s)
- [zone  - player movements](https://www.youtube.com/watch?v=WPR9IcppmkI)
- [zone  - entity interaction range](https://www.youtube.com/watch?v=cSnldVbl2wA&feature=youtu.be)
- [zone  - entity interaction range 2](https://www.youtube.com/watch?v=roSZNHxg7o4)
- [zone  - monsters!!](https://www.youtube.com/watch?v=f7nPVcIaKfw)


## Development setup

    git clone https://github.com/shine-o/shine.engine.emulator
    
    cd shine.engine.emulator
    
    go mod vendor
    
    cp .env.dist .env
    
    docker-compose up --build
        
    # if you made any change to a service:
    docker-compose up -d --force-recreate --no-deps <service-name>
    
    
If everything is okay, you should see something like this when using **docker-compose ps**:


![](assets/docker-services.PNG)    

## Metrics
   
For metrics I use the following services:    
    - Prometheus
    - Loki
    - Grafana

The services are configured and ready to use in the **docker-compose.yml** file. You can get something like this:

![](assets/grafana.PNG)
    
## Event logic for login, world, zone services

#### From tcp connection to network command to logic handler

![](docs/zone-logic.PNG)


#### Processes and events example

![](docs/process-events.PNG)    