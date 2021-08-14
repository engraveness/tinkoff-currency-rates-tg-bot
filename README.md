Tinkoff currency rates telegram bot
===
This project is a simplified version of the telegram bot for displaying the exchange rate of the online bank [Tinkoff](https://tinkoff.ru)  made specifically for hosting on raspberry pi.
The bot will respond to any of your messages with its message containing the current exchange rate. Only `EUR` and `USD` are available.

### Disclaimer

All product and company names are the registered trademarks of their original owners. The use of any trade name or trademark is for identification and reference purposes only and does not imply any association with the trademark holder of their product brand.

## Interaction scenarios
There are two modes of interaction with bot:
- passive update - default mode. Bot will edit it's last message to show actual exchange rates without distracting notifications;
- threshold mode - you can set the threshold, and if the change in the new exchange rates exceeds it, you will receive a new message. For more information, see the commands section.

## Menu commands
- `/start` subscribe to the currency rate updates
- `/stop` unsubscribe from the bot
- `/threshold` - show current rate threshold. If the change of one of the new exchange rates exceeds this limit, you will receive a new message with the current exchange rate.
- `/threshold value` - sets up the new threshold

## Command line arguments
The only supported command line argument is the path to the configuration file. If it is not installed, it will try to open the batch .json file in the current directory.

## Config.json
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
  "usersFilePath": "map.json",
  "saveTheSentRates": false
}
```

All refresh rates are specified in seconds.
`saveTheSentRates` flag is necessary to reduce the usage of the raspberry pi memory card. The sent exchange rates will be stored in memory and may (or may not) eventually be flushed to disk.

## Compiling for raspberry pi
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
