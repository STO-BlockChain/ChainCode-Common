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
