Tinkoff currency rates telegram bot
===
This project is a simplified version of the telegram bot for displaying the exchange rate of the online bank [Tinkoff](https://tinkoff.ru)  made specifically for hosting on raspberry pi.
The bot will respond to any of your messages with its message containing the current exchange rate. With the receipt of new data, it will edit its message. So, this is something close to a widget, but executed in a telegram. The only supported currencies are `EUR` and `USD`.

### Disclaimer

All product and company names are the registered trademarks of their original owners. The use of any trade name or trademark is for identification and reference purposes only and does not imply any association with the trademark holder of their product brand.

# Menu commands
The only available commands are `/start` and `/stop`

# Command line arguments
The only supported command line argument is the path to the configuration file. If it is not installed, it will try to open the batch .json file in the current directory.

# Config.json
The main directory contains [an example of a config](config.example.json)
If not, here is desirable format:
```
{
  "refreshRate": {
    "currency": 300,
    "messages": 10,
    "mainCycle": 5
  },
  "telegramToken": "YOUR_TOKEN_HERE",
  "usersFilePath": "map.json"
}
```

All refresh rates are specified in seconds.

# Compiling for raspberry pi
### From windows
Create and execute from the root directory a bat script with the following contents:
```
set GOARCH=arm
set GOOS=linux
set GOARM=5
go build
```

### From linux
Run specified command from a root directory:
```
env GOOS=linux GOARCH=arm GOARM=5 go build
```
