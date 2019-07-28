1. Update `filesMap.json` with the list of files -> Copy you want to backup
2. Copy `com.eddie.filesWatcherDaemon.plist` to `~/Library/LaunchAgents/com.eddie.filesWatcherDaemon.plist` 
to load it to launchd do:
    `launchctl load .../com.eddie.filesWatcherDaemon.plist`

```Usage of ./filesWatcher:
  -dst src
        the destination file path used only with src
  -f string
        a json file with the list of files (default "./filesMap.json")
  -src string
        the source file path```

Using this for [dotFiles](https://github.com/edwardIshaq/dotFiles)

