package btc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

//UTXO represents unspent transaction outputs.
type UTXO struct {
	Addr   string
	Hash   []byte
	Amount uint64
	Index  uint32
	Script []byte
	Age    uint64
	Key    *Key
}

//UTXOs is for sorting UTXO
type UTXOs []*UTXO

var cacheUTXO = make(map[string]UTXOs)

//Service is for getting UTXO or sending transactions , basically by using WEB API.
type Service interface {
	GetServiceName() string
	GetUTXO(string, *Key) (UTXOs, error)
	SendTX([]byte) ([]byte, error)
}

//TestServices is an array containing generator of Services for testnet
var TestServices = []func() (Service, error){
	NewBlockrServiceForTest,
}

//TestServices is an array containing generator of Services
var Services = []func() (Service, error){
	NewBlockrService,
}

//to sort UTXO

//Len returns length of UTXO
func (us UTXOs) Len() int {
	return len(us)
}

//Swap swaps UTXO
func (us UTXOs) Swap(i, j int) {
	us[i], us[j] = us[j], us[i]
}

//Less returns true is age is smaller.
func (us UTXOs) Less(i, j int) bool {
	return us[i].Amount < us[j].Amount
}

//SelectService returns a service randomly.
func SelectService(isTestnet bool) (Service, error) {
	n := rand.Int() % len(Services)
	if isTestnet {
		return TestServices[n]()
	}
	return Services[n]()
}

//SetTXSpent sets  tx hash is already spent.
func SetUTXOSpent(hash []byte) {
	for k, v := range cacheUTXO {
		for i, utxo := range v {
			if bytes.Compare(hash, utxo.Hash) == 0 {
				v = append(v[0:i], v[i+1:]...)
				cacheUTXO[k] = v
				return
			}
		}
	}
}

type unspent struct {
	Status string
	Data   struct {
		Address string
		Unspent []struct {
			Tx            string
			Amount        string
			N             int
			Confirmations int
			Script        string
		}
	}
	Code    int
	Message string
}

type sendtx struct {
	Status  string
	Data    string
	Code    int
	Message string
}

//BlockrService is a service using Blockr.io.
type BlockrService struct {
	isTestnet bool
}

//NewBlockrServiceForTest creates BlockrService struct for test.
func NewBlockrServiceForTest() (Service, error) {
	b := &BlockrService{isTestnet: true}
	return b, nil
}

//NewBlockrService creates BlockrService struct for not test.
func NewBlockrService() (Service, error) {
	return &BlockrService{isTestnet: false}, nil
}

//GetServiceName return service name.
func (b *BlockrService) GetServiceName() string {
	return "BlockrService"
}

//SendTX send a transaction using Blockr.io.
func (b *BlockrService) SendTX(data []byte) ([]byte, error) {
	var btc string

	if b.isTestnet {
		btc = "tbtc"
	} else {
		btc = "btc"
	}

	resp, err := http.PostForm("http://"+btc+".blockr.io/api/v1/tx/push",
		url.Values{"hex": {hex.EncodeToString(data)}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	logging.Println(string(body))
	var u sendtx
	err = json.Unmarshal(body, &u)
	if err != nil {
		return nil, err
	}
	if u.Status != "success" {
		return nil, errors.New("blockr returns " + u.Message)

	}
	return hex.DecodeString(u.Data)
}

//GetUTXO gets unspent transaction outputs by using Blockr.io.
func (b *BlockrService) GetUTXO(addr string, key *Key) (UTXOs, error) {
	if cacheUTXO[addr] != nil {
		return cacheUTXO[addr], nil
	}
	var btc string

	if b.isTestnet {
		btc = "tbtc"
	} else {
		btc = "btc"
	}
	//http://btc.blockr.io/api/v1/address/unspent/18BjFdiThEtu7D8hF3yURRPyPh9gNkRcBB
	resp, err := http.Get("http://" + btc + ".blockr.io/api/v1/address/unspent/" + addr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var u unspent
	err = json.Unmarshal(body, &u)
	if err != nil {
		return nil, err
	}
	if u.Status != "success" {
		return nil, errors.New("blockr returns " + u.Message)
	}

	utxos := make(UTXOs, 0, len(u.Data.Unspent))
	for _, tx := range u.Data.Unspent {
		utxo := UTXO{}
		utxo.Addr = addr
		amount, err := strconv.ParseFloat(tx.Amount, 64)
		if err != nil {
			return nil, err
		}
		utxo.Amount = uint64(amount * BTC)
		utxo.Hash, err = hex.DecodeString(tx.Tx)
		if err != nil {
			return nil, err
		}
		utxo.Index = uint32(tx.N)
		utxo.Script, err = hex.DecodeString(tx.Script)
		if err != nil {
			return nil, err
		}
		utxo.Age = uint64(tx.Confirmations)
		utxo.Key = key
		utxos = append(utxos, &utxo)
	}
	if key != nil {
		cacheUTXO[addr] = utxos
	}
	return utxos, nil
}
