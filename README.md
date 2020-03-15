![](shine.png)
# Shine Engine World Service

> Emulated World service.


This project has dependencies on the modules: 

- [Networking](https://github.com/shine-o/shine.engine.networking)
- [Structs](https://github.com/shine-o/shine.engine.structs)
- [Protocol Buffers](https://github.com/shine-o/shine.engine.protocol-buffers)

If changes are needed on these modules, inside this project, use [gohack](https://github.com/rogpeppe/gohack) (a tool to go around using cached modules): 


    $ gohack get -vcs github.com/shine-o/shine.engine.networking 
    $ gohack get -vcs github.com/shine-o/shine.engine.structs
    $ gohack get -vcs github.com/shine-o/shine.engine.protocol-buffers


Open those projects, change code, commit, submit pull request, etc.. 

When done, **BEFORE COMMITING**, revert changes: 

    $ gohack undo github.com/shine-o/shine.engine.networking
    $ gohack undo github.com/shine-o/shine.engine.structs
    $ gohack undo github.com/shine-o/shine.engine.protocol-buffers
    
    
Obvious git practices like committing, not leaving garbage files, etc... are required to avoid problems. More info on [gohack](https://github.com/rogpeppe/gohack)