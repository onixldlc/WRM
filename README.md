# WRM - Windows Resolution Manager
WRM is a cli tool to change resolution without the need of, going to desktop -> click left -> resolution -> change resolution, with WRM you can define your own configuration, want 1920x1080 with 75hz but windows keep defaulting to 60hz everytime you change it in the settings ? well this tool can help you with preconfigured config, all you need to do is just create a preset in json, then followed by running `./wrm config <config-name>` and it will automatically change your desktop resolution to the new resolution you have specified!

## build
to build the project all you need is golang then do 
```
go build -o wrm.exe
```
and boom you are done now all you need is to preconfigure your setup

## TODO:
1. ~EVERYTHING! (still working on listing!)~ well... to a certain degree
2. Error Logging
3. Config backup incase of crashes
4. a way to add config via cli