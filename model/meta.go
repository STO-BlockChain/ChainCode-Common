package model

// TransferMeta is Multi Transfer를 이용하기 위한 데이터 구조체
type TransferMeta struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount,omitempty"`
}

type TransferMetaN struct {
	FromAddress string `json:"fromaddress"`
	ToAddress   string `json:"toaddress"`
	Amount      uint64 `json:"amount,omitempty"`
}

/*
// WalletMeta is 지갑 데이터 구조체
type WalletMeta struct {
	Publickey  string `json:"publickey,omitempty"`
	Txtime     string `json:"txtime,omitempty"`
	Nowtime    int64  `json:"nowtime,omitempty"`
	Transdata  string `json:"transdata,omitempty"`
	Transjdata string `json:"transjdata,omitempty"`
	Sigmsg     string `json:"sigmsg,omitempty"`
}
*/
