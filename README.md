# eos_balance_server
### clone
```
$ git clone https://github.com/gaozhengxin/eos_balance_tracker.git
```
### config
Edit config/config.go.
* `EOS_ACCOUNT`: owner's account
* `NODEOS`: api path of a nodeos configurated `filter_on=*` (allowing to retrieve transactions)
Rebuild program.
### build
```
$ cd eos_balance_tracker
$ go build ./server/balance_server.go
```
### run
```
$ ./balance_server
```
For the first time, add flag `--reinit=true`.
### flag
* `reinit`: initiate or rebuild database, default false
* `dbpath`: absolute or relative path of database
* `port`: listening port, default 7000
### http api
```
http://0.0.0.0:7000/get_balance?user_key=drefigvhv1ywn2yvjpgnu53extxb1xi4l1
```
