//Exercise 9.1: Add a function Withdraw(amount int) bool to the gopl.io/ch9/bank1
// program. The result should indicate whether the transaction succeeded or failed due to insuf-
// ficient funds. The message sent to the monitor goroutine must contain both the amount to
// withdraw and a new channel over which the monitor goroutine can send the boolean result
// back to Withdraw.

package bank

type WithdrawInfo struct {
	amount int
	result chan bool
}

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance
var withdraws = make(chan WithdrawInfo)

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }
func Withdraw(amount int) bool {
	info := WithdrawInfo{amount, make(chan bool)}
	withdraws <- info
	return <-info.result
}

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case info := <-withdraws:
			info.result <- (balance > info.amount)
			if balance > info.amount {
				balance -= info.amount
			}
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}

//!-
