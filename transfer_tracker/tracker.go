package tracker

import (
	"../dao"
	"../config"
	"../utils"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"github.com/gaozhengxin/cryptocoins/src/go/rpcutils"
)

var dbPath string
var err error

func Run () {

	dbpath := flag.String("dbpath", "./data", "database path")
	reinit := flag.Bool("reinit", false, "reinit")

	flag.Parse()

	if *reinit {
		initDb()
	}

	dbPath, err = filepath.Abs(*dbpath)
	if err != nil {
		panic(err)
	}
	config.DbPath = dbPath

	exists, err := PathExists(dbPath)
	if err != nil {
		panic(err)
	}
	if !exists {
		initDb()
	}
	log.Printf("database path is %v\n", config.DbPath)

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	} ()

	for {
		cursor := dao.Get("cursor")
		pos, e := strconv.Atoi(cursor)
		if e != nil {
			err = fmt.Errorf("cursor not found\ntry add reinit\n")
			log.Fatal(err)
		}
		transfers := scan(pos)
		updateBalances(transfers)
		time.Sleep(time.Duration(60)*time.Second)
	}
}

func scan(pos int) []transfer {
	api := "v1/history/get_actions"
	offset := 99
	var transfers []transfer
	var almostutd = false
	for !almostutd {
		data := `{"pos":` + strconv.Itoa(pos) + `,"offset":` + strconv.Itoa(offset) + `,"account_name":"` + config.EOS_ACCOUNT + `"}`
		ret := rpcutils.DoCurlRequest(config.NODEOS, api, data)
		trs, err := parseResult(ret)
		if err != nil {
			panic(err)
		}
		transfers = append(transfers, trs...)
		l := len(transfers)
		pos += l
		dao.Put("cursor", strconv.Itoa(pos))
		if l < offset + 1 {
			almostutd = true
		}
		time.Sleep(time.Duration(2000000))
	}
	return transfers
}

type transfer struct {
	seq float64;
	blockNum float64;
	txid string;
	memo string;
	quantity string;
}

func parseResult(ret string) ([]transfer, error) {
	var tfs []transfer
	var result map[string]interface{}
	err = json.Unmarshal([]byte(ret),&result)
	if err != nil {
		return nil, fmt.Errorf("failed parsing rpc response")
	}
	actions := result["actions"].([]interface{})
	for _, action := range actions {
		receipt := action.(map[string]interface{})["action_trace"].(map[string]interface{})["receipt"]
		receiver := receipt.(map[string]interface{})["receiver"].(string)
		if receiver != config.EOS_ACCOUNT {
			continue
		}
		act := action.(map[string]interface{})["action_trace"].(map[string]interface{})["act"]
		name := act.(map[string]interface{})["name"].(string)
		if name != "transfer" {
			continue
		}
		data := act.(map[string]interface{})["data"].(map[string]interface{})
		from := data["from"].(string)
		to := data["to"].(string)
		quantity := data["quantity"].(string)
		memo := data["memo"].(string)
		if !utils.IsUserKey(memo) {
			continue
		}

		quantity = utils.ParseQuantity(quantity)
		if from == config.EOS_ACCOUNT {
			quantity = "-" + quantity
			// withdraw
		} else if to == config.EOS_ACCOUNT{
			// deposit
		} else {
			continue  // 不可能
		}

		tf := new(transfer)
		tf.quantity = quantity
		tf.memo = memo
		tf.seq = action.(map[string]interface{})["account_action_seq"].(float64)
		tf.blockNum = action.(map[string]interface{})["block_num"].(float64)
		tf.txid = action.(map[string]interface{})["action_trace"].(map[string]interface{})["trx_id"].(string)
		tfs = append(tfs, *tf)
	}
	return tfs, nil
}

func updateBalances(transfers []transfer) {
	if len(transfers) == 0 {
		return
	}
	for _, tf := range transfers {
		log.Printf("updating balance for %s...\n", tf.memo)
		log.Printf("balance change is %v\n", tf.quantity)
		err = dao.UpdateBalance(tf.memo, tf.quantity)
		if err != nil {
			log.Fatal(err)
		}
		//bal := dao.GetBalance(tf.memo)
		//log.Printf("%v's balance is %v.", tf.memo, bal)
	}
}

func PathExists(path string) (bool, error) {
	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func initDb() {
	err = os.RemoveAll(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	dao.Put("cursor","0")
}
