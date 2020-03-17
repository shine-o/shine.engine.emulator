![](shine.png)
                     # Shine Engine Login Service
                     
                     > Emulated login service.
                     
                     ---
                     <p align="left">
                         <a href="https://circleci.com/gh/shine-o/shine.engine.login/tree/dev"><img src="https://circleci.com/gh/shine-o/shine.engine.login/tree/dev.svg?style=shield" alt="Build Status"></a>
                     </p>
                     
                     This project has dependencies on the modules: 
                     
                     - [Networking](https://github.com/shine-o/shine.engine.networking)
                     - [Structs](https://github.com/shine-o/shine.engine.structs)
                     - [Protocol Buffers](https://github.com/shine-o/shine.engine.protocol-buffers)
                     
                     If changes are needed on these modules, append to the file go.mod:
                            
                         replace github.com/shine-o/shine.engine.networking => C:\Users\marbo\go\src\github.com\shine-o\shine.engine.networking
                         replace github.com/shine-o/shine.engine.structs => C:\Users\marbo\go\src\github.com\shine-o\shine.engine.structs
                         replace github.com/shine-o/shine.engine.protocol-buffers => C:\Users\marbo\go\src\github.com\shine-o\shine.engine.protocol-buffers
                     
                     Obvious git practices like committing, not leaving garbage files, etc... are required to avoid problems. 
                     
                     
                     With any other module, you can use [gohack](https://github.com/rogpeppe/gohack)
                     
                     # Overview
                     ![](login-scheme.png)