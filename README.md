![](shine.png)
# Shine Engine World Service

> Emulated World service.
---
[![CircleCI](https://circleci.com/gh/shine-o/shine.engine.world/tree/master.svg?style=shield)](https://circleci.com/gh/shine-o/shine.engine.world/tree/master.svg?style=shield)
[![Go Report Card](https://goreportcard.com/badge/github.com/shine-o/shine.engine.world)](https://goreportcard.com/report/github.com/shine-o/shine.engine.world)

This project has dependencies on the modules: 

- [Networking](https://github.com/shine-o/shine.engine.networking)
- [Structs](https://github.com/shine-o/shine.engine.networking/structs)
- [Protocol Buffers](https://github.com/shine-o/shine.engine.protocol-buffers)


If quick changes for testing are needed on these modules, append to the file go.mod:
       
    replace github.com/shine-o/shine.engine.networking => C:\Users\marbo\go\src\github.com\shine-o\shine.engine.networking
    replace github.com/shine-o/shine.engine.protocol-buffers => C:\Users\marbo\go\src\github.com\shine-o\shine.engine.protocol-buffers

