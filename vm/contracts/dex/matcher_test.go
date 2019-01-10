package dex

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"github.com/vitelabs/go-vite/common/types"
	"math/big"
	"testing"
	"time"
)

var (
	VITE = types.TokenTypeId{'V', 'I', 'T', 'E', ' ', 'T', 'O', 'K', 'E', 'N'}
	ETH = types.TokenTypeId{'E', 'T', 'H', ' ', ' ', 'T', 'O', 'K', 'E', 'N'}
	NANO = types.TokenTypeId{'N', 'A', 'N', 'O', ' ', 'T', 'O', 'K', 'E', 'N'}
)

func TestMatcher(t *testing.T) {
	localStorage := NewMapStorage()
	st := BaseStorage(&localStorage)
	mc := NewMatcher(getAddress(), &st)
	// buy
	buy1 := newOrderInfo(101, ETH, VITE, false, Limited, "100.02", 10, time.Now().UnixNano()/1000)
	buy2 := newOrderInfo(102, ETH, VITE, false, Limited, "100.03", 30, time.Now().UnixNano()/1000)
	buy3 := newOrderInfo(103, ETH, VITE, false, Limited, "100.02", 70, time.Now().UnixNano()/1000)
	buy4 := newOrderInfo(104, ETH, VITE, false, Limited, "100.04", 45, time.Now().UnixNano()/1000)
	buy5 := newOrderInfo(105, ETH, VITE, false, Limited, "100.01", 20, time.Now().UnixNano()/1000)
	mc.MatchOrder(buy1)
	mc.MatchOrder(buy2)
	mc.MatchOrder(buy3)
	mc.MatchOrder(buy4)
	mc.MatchOrder(buy5)
	assert.Equal(t, 5, len(localStorage.logs))

	bookNameToMakeForBuy := getBookIdToMakeForOrder(buy5)
	assert.Equal(t,104, fromOrderIdToInt(mc.books[bookNameToMakeForBuy].header))

	sell1 := newOrderInfo(201, ETH, VITE, true, Limited, "100.1", 100, time.Now().UnixNano()/1000)
	sell2 := newOrderInfo(202, ETH, VITE, true, Limited, "100.02", 50, time.Now().UnixNano()/1000)
	mc.MatchOrder(sell1)
	assert.Equal(t, 6, len(localStorage.logs))
	mc.MatchOrder(sell2)

	bookNameToMakeForSell := getBookIdToMakeForOrder(sell1)
	assert.Equal(t,int32(1), mc.books[bookNameToMakeForSell].length)
	assert.Equal(t,201, fromOrderIdToInt(mc.books[bookNameToMakeForSell].header))
	assert.Equal(t,int32(4), mc.books[bookNameToMakeForBuy].length)
	assert.Equal(t,102, fromOrderIdToInt(mc.books[bookNameToMakeForBuy].header))

	assert.Equal(t, 11, len(localStorage.logs))
	pl, _, _, _  := mc.books[bookNameToMakeForSell].getByKey(mc.books[bookNameToMakeForBuy].header)
	od, _ := (*pl).(Order)
	assert.True(t, big.NewInt(5).Cmp(new(big.Int).SetBytes(od.ExecutedQuantity)) == 0)

	log := localStorage.logs[9]
	txEvent := TransactionEvent{}
	txEvent = txEvent.fromBytes(log.Data).(TransactionEvent)
	assert.Equal(t, 202, fromOrderIdBytesToInt(txEvent.TakerId))
	assert.Equal(t, 104, fromOrderIdBytesToInt(txEvent.MakerId))
	assert.True(t, CheckBigEqualToInt(45, txEvent.Quantity))
	assert.True(t, priceEqual(txEvent.Price, "100.04"))
	assert.True(t, CheckBigEqualToInt(4502, txEvent.Amount))

	log = localStorage.logs[10]
	txEvent = TransactionEvent{}
	txEvent = txEvent.fromBytes(log.Data).(TransactionEvent)
	assert.Equal(t, 202, fromOrderIdBytesToInt(txEvent.TakerId))
	assert.Equal(t, 102, fromOrderIdBytesToInt(txEvent.MakerId))
	assert.True(t, CheckBigEqualToInt(5,txEvent.Quantity))
	assert.True(t, priceEqual(txEvent.Price, "100.03"))
	assert.True(t, CheckBigEqualToInt(500, txEvent.Amount))

	buy6 := newOrderInfo(106, ETH, VITE, false, Limited, "100.3", 101, time.Now().UnixNano()/1000)
	mc.MatchOrder(buy6)
	bookNameForBuy := getBookIdToMakeForOrder(buy6)
	bookNameForSell := getBookIdToTake(buy6)
	assert.Equal(t,106, fromOrderIdToInt(mc.books[bookNameForBuy].header))
	assert.Equal(t,0, fromOrderIdToInt(mc.books[bookNameForSell].header))
	assert.Equal(t,int32(5), mc.books[bookNameForBuy].length)
	assert.Equal(t,int32(0), mc.books[bookNameForSell].length)

	assert.Equal(t, 14, len(localStorage.logs))

	log = localStorage.logs[11]
	odEvent := OrderUpdateEvent{}
	odEvent = odEvent.fromBytes(log.Data).(OrderUpdateEvent)
	assert.Equal(t, 106, fromOrderIdBytesToInt(odEvent.Id))
	assert.Equal(t, PartialExecuted, int(odEvent.Status))
	assert.True(t, CheckBigEqualToInt(100, odEvent.ExecutedQuantity))
	assert.True(t, CheckBigEqualToInt(10010, odEvent.ExecutedAmount))

	log = localStorage.logs[12]
	odEvent = OrderUpdateEvent{}
	odEvent = odEvent.fromBytes(log.Data).(OrderUpdateEvent)
	assert.Equal(t, 201, fromOrderIdBytesToInt(odEvent.Id))
	assert.Equal(t, FullyExecuted, int(odEvent.Status))
	assert.True(t, CheckBigEqualToInt(100, odEvent.ExecutedQuantity))
	assert.True(t, CheckBigEqualToInt(10010, odEvent.ExecutedAmount))

	log = localStorage.logs[13]
	txEvent = TransactionEvent{}
	txEvent = txEvent.fromBytes(log.Data).(TransactionEvent)
	assert.Equal(t, 106, fromOrderIdBytesToInt(txEvent.TakerId))
	assert.Equal(t, 201, fromOrderIdBytesToInt(txEvent.MakerId))
	assert.True(t, CheckBigEqualToInt(100, txEvent.Quantity))
	assert.True(t, priceEqual(txEvent.Price, "100.1"))
	assert.True(t, CheckBigEqualToInt(10010, txEvent.Amount))
}

func TestDust(t *testing.T) {
	localStorage := NewMapStorage()
	st := BaseStorage(&localStorage)
	mc := NewMatcher(getAddress(), &st)
	// buy quantity = origin * 100,000,000
	buy1 := newOrderInfo(301, VITE, ETH, false, Limited, "0.0012345", 100000000, time.Now().UnixNano()/1000)
	mc.MatchOrder(buy1)
	// sell
	sell1 := newOrderInfo(401, VITE, ETH,true, Limited, "0.0012342", 100000200, time.Now().UnixNano()/1000)
	mc.MatchOrder(sell1)

	bookNameToMakeForBuy := getBookIdToMakeForOrder(buy1)
	bookNameToMakeForSell := getBookIdToMakeForOrder(sell1)
	assert.Equal(t,int32(0), mc.books[bookNameToMakeForBuy].length)
	assert.Equal(t,int32(0), mc.books[bookNameToMakeForSell].length)
	assert.Equal(t,0, fromOrderIdToInt(mc.books[bookNameToMakeForBuy].header))
	assert.Equal(t,0, fromOrderIdToInt(mc.books[bookNameToMakeForSell].header))

	assert.Equal(t, 4, len(localStorage.logs))
	log := localStorage.logs[1]
	orderEvent := OrderUpdateEvent{} //taker
	orderEvent = orderEvent.fromBytes(log.Data).(OrderUpdateEvent)
	assert.Equal(t, 401, fromOrderIdBytesToInt(orderEvent.Id))
	assert.Equal(t, FullyExecuted, int(orderEvent.Status))
	assert.True(t, CheckBigEqualToInt(100000000, orderEvent.ExecutedQuantity))
	assert.True(t, CheckBigEqualToInt(123450, orderEvent.ExecutedAmount))
	assert.Equal(t, VITE.Bytes(), orderEvent.RefundToken)
	assert.True(t, CheckBigEqualToInt(200, orderEvent.RefundQuantity))

	log = localStorage.logs[2]
	orderEvent = OrderUpdateEvent{} // maker
	orderEvent = orderEvent.fromBytes(log.Data).(OrderUpdateEvent)
	assert.Equal(t, 301, fromOrderIdBytesToInt(orderEvent.Id))
	assert.Equal(t, FullyExecuted, int(orderEvent.Status))
	assert.True(t, CheckBigEqualToInt(100000000, orderEvent.ExecutedQuantity))
	assert.True(t, CheckBigEqualToInt(123450, orderEvent.ExecutedAmount))
	assert.Equal(t, ETH.Bytes(), orderEvent.RefundToken)
	assert.True(t, CheckBigEqualToInt(0, orderEvent.RefundQuantity))

	log = localStorage.logs[3]
	txEvent := TransactionEvent{} // maker
	txEvent = txEvent.fromBytes(log.Data).(TransactionEvent)
	assert.Equal(t, true, txEvent.TakerSide)
	assert.Equal(t, 401, fromOrderIdBytesToInt(txEvent.TakerId))
	assert.Equal(t, 301, fromOrderIdBytesToInt(txEvent.MakerId))
	assert.True(t, priceEqual("0.0012345", txEvent.Price))
	assert.True(t, CheckBigEqualToInt(100000000, txEvent.Quantity))
	assert.True(t, CheckBigEqualToInt(123450, txEvent.Amount))
}

func TestMarket(t *testing.T) {

}

func newOrderInfo(id int, tradeToken types.TokenTypeId, quoteToken types.TokenTypeId, side bool, orderType int32, price string, quantity uint64, ts int64) Order {
	order := Order{}
	order.Id = orderIdBytesFromInt(id)
	order.TradeToken = tradeToken.Bytes()
	order.QuoteToken = quoteToken.Bytes()
	order.Side = side // buy
	order.Type = orderType
	order.Price = price
	order.Quantity = new(big.Int).SetUint64(quantity).Bytes()
	order.Timestamp = int64(ts)
	order.Status = Pending
	order.Amount = CalculateAmount(order.Quantity, order.Price)
	order.ExecutedQuantity = big.NewInt(0).Bytes()
	order.ExecutedAmount = big.NewInt(0).Bytes()
	order.RefundToken = []byte{}
	order.RefundQuantity = big.NewInt(0).Bytes()
	return order
}

func CheckBigEqualToInt(expected int, value []byte) bool {
	return new(big.Int).SetUint64(uint64(expected)).Cmp(new(big.Int).SetBytes(value)) == 0
}

func fromOrderIdToInt(orderId nodeKeyType) int {
	return fromOrderIdBytesToInt(orderId.(OrderId).bytes())
}

func fromOrderIdBytesToInt(orderIdBytes []byte) int {
	res := int(binary.LittleEndian.Uint32(orderIdBytes[16:20]))
	return res
}