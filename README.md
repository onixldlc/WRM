# WRM - Windows Resolution Manager
WRM is a cli tool to change resolution without the need of, going to desktop -> click left -> resolution -> change resolution, with WRM you can define your own configuration, want 1920x1080 with 75hz but windows keep defaulting to 60hz everytime you change it in the settings ? well this tool can help you with preconfigured config, all you need to do is just create a preset in json, then followed by running 
```
./wrm config <config-name>
```
and it will automatically change your desktop resolution to the new resolution you have specified!

## build
to build the project all you need is golang then do 
```
go build -o wrm.exe
```
and boom you are done now all you need is to preconfigure your setup

## config
For configuration you can use the id when you do a `./wrm list` or if your monitor id keep changing you can use the model name instead, although if you have 2 monitor with the same brand and model this might be an issue for you and best thing you can do i just to use id instead of your monitor model name, for configuration you can also use both the id from `./wrm list` or the model name of that monitor for [example](https://github.com/onixldlc/WRM/blob/main/config.json):
```json
{
    "configurations": [
      {
        "name": "Gaming Setup",
        "monitor": 1,
        "resolution": "1920x1080",
        "frequency": 180
      },
      {
        "name": "Work Setup",
        "monitor_name": "27G2G5",
        "resolution": "2560x1440",
        "frequency": 60
      }
    ]
  }
```

as you can see you can either use key "monitor" to refere to monitor id, or "monitor_name" to the monitor model name, without both "monitor" and "monitor_name" WRM would run just fine, and if the configuration is loaded it will be fine until you try to apply it, when you apply it it will print out `Monitor index in configuration is out of range.` since the id default to 0 if both "monitor" and "monitor_name" doesn't exist!

> [!NOTE]
> If configuration uses space in between the name, you will need to add " to apply it, for example `./WRM config "Gaming Setup"` 


## TODO:
1. ~EVERYTHING! (still working on listing!)~ well... to a certain degree
2. Error Logging
3. Config backup incase of crashes
4. a way to add config via cli
