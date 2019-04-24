# eos_balance_server
### build
```
$ git clone https://github.com/gaozhengxin/eos_balance_tracker.git
$ cd eos_balance_tracker
$ ./build.sh
```
### config
Edit config/config.go.
* `EOS_ACCOUNT`: owner's account
* `NODEOS`: api path of a nodeos configurated `filter_on=*` (allowing to retrieve transactions)
Rebuild program.
### run
```
$ ./balance_server
```
For the first time, run with `--reinit=true`.
### flag
* `reinit`: initiate or rebuild database, default false
* `dbpath`: absolute or relative path of database
* `port`: listening port, default 7000
### http api
```
http://0.0.0.0:7000/get_balance?user_key=drefigvhv1ywn2yvjpgnu53extxb1xi4l1
```
