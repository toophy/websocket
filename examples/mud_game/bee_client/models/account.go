package models

import(
	"sync"
	"github.com/gorilla/websocket"
)

// AccountConn 帐号的网络连接
type AccountConn struct {
	Account string
	Pwd     string
	Mt      int
	Temp    bool // 临时网络连接
	C       *websocket.Conn
	C1      chan string
}

var(
	gAccounts map[string]*AccountConn
	gAccountLock *sync.Mutex
)

func init(){
	gAccounts = make(map[string]*AccountConn,0)
	gAccountLock = new(sync.Mutex)
}

func NewAccount(account string) *AccountConn{
	gAccountLock.Lock()
	defer gAccountLock.Unlock()

	if _,ok:=gAccounts[account];ok{
		return gAccounts[account];
	}
	
	gAccounts[account] = &AccountConn{
		Account : account,
		Pwd :"",
		Mt : 0,
		Temp:true,
		C:nil,
		C1:make(chan string)}

	return gAccounts[account]
}

func GetAccount(account string) *AccountConn{
	gAccountLock.Lock()
	defer gAccountLock.Unlock()

	if _,ok:=gAccounts[account];ok{
		return gAccounts[account]
	}
	return nil
}

func LeaveAccount(account string) {
	gAccountLock.Lock()
	defer gAccountLock.Unlock()

	delete(gAccounts,account)
}