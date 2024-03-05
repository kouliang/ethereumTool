package main

import (
	"fmt"

	"github.com/kouliang/ethereumtool/account"
	"github.com/kouliang/ethereumtool/email"
)

func main() {
	fmt.Println(account.IsAvailableAddress("0xA6e7Ce1c292E5d52508b58e2EC52E3D741793679"))
}

func emailTest() {
	content :=
		`Task success:
dividingTime:1659600000
stakeListLength:5251
totalPower:764804
contract address:0x6fce9980e5527BBc2E925c5af28B5da26480786F
nonce:69 gasPrice:5000000000 gasLimit:47364
tx broadcast:0x8cbd839b3f34534b1f46d8f629ca248a3cbec2cfbab08a2d5d544cbf8f4ddaf3
receipted - status:1, blockNumber:20140996`

	msg, err := email.SenEmail(" ", content, []string{"coderkl@qq.com", "kouliangg@gmail.com"})
	if err != nil {
		fmt.Println("SendEmail failed:", err.Error())
	} else {
		fmt.Println(msg)
	}
}
