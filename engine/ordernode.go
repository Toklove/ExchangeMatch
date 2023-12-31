package engine

import (
	"github.com/shopspring/decimal"
	"gome/api"
	"math"
)

type OrderNode struct {
	Action      int8    // 节点操作类型，1-add，2-del
	Uuid        string  // 用户唯一标识
	Oid         string  // 订单唯一标识
	Symbol      string  // 交易对
	Transaction int32   // 交易方向，buy/sale
	Price       float64 // 交易价格
	Volume      float64 // 交易数量
	Accuracy    int     // 计算精度
	NodeName    string  // 节点
	IsFirst     bool    // 是否是起始节点
	IsLast      bool    // 是否是结束节点
	PrevNode    string  // 前一个节点
	NextNode    string  // 后一个节点
	NodeLink    string  // 节点链标识

	// hash对比池标识.
	OrderHashKey   string
	OrderHashField string

	// 有序集合委托列表.
	OrderListSortSetKey  string
	OrderListSortSetRKey string // 相反的委托

	// hash委托深度.
	OrderDepthHashKey   string
	OrderDepthHashField string
}

func NewOrderNode(order *api.OrderRequest) *OrderNode {
	node := &OrderNode{}
	node.SetAccuracy()
	node.SetUuid(order)
	node.SetOid(order)
	node.SetSymbol(order)
	node.SetTransaction(order)
	node.SetVolume(order)
	node.SetPrice(order)
	node.SetOrderHashKey()
	node.SetListSortSetKey()
	node.SetDepthHashKey()
	node.SetNodeName()
	node.SetNodeLink()

	return node
}

func (node *OrderNode) SetAccuracy() {
	node.Accuracy = Conf.MeConf.Accuracy
}

func (node *OrderNode) SetUuid(order *api.OrderRequest) {
	node.Uuid = order.Uuid
}

func (node *OrderNode) SetOid(order *api.OrderRequest) {
	node.Oid = order.Oid
}

func (node *OrderNode) SetSymbol(order *api.OrderRequest) {
	node.Symbol = order.Symbol
}

func (node *OrderNode) SetTransaction(order *api.OrderRequest) {
	node.Transaction = int32(order.Transaction)
}

func (node *OrderNode) SetVolume(order *api.OrderRequest) {
	volume := decimal.NewFromFloat(order.Volume)
	mul := decimal.NewFromFloat(math.Pow10(node.Accuracy))

	node.Volume, _ = volume.Mul(mul).Float64()
}

func (node *OrderNode) SetPrice(order *api.OrderRequest) {
	volume := decimal.NewFromFloat(order.Price)
	mul := decimal.NewFromFloat(math.Pow10(node.Accuracy))
	node.Price, _ = volume.Mul(mul).Float64()
}

func (node *OrderNode) SetOrderHashKey() {
	node.OrderHashKey = node.Symbol + ":comparison"
	node.OrderHashField = node.Symbol + ":" + node.Uuid + ":" + node.Oid
}

func (node *OrderNode) SetListSortSetKey() {
	if api.TransactionType_value["SELL"] == node.Transaction {
		node.OrderListSortSetKey = node.Symbol + ":SELL"
		node.OrderListSortSetRKey = node.Symbol + ":BUY"
	} else {
		node.OrderListSortSetKey = node.Symbol + ":BUY"
		node.OrderListSortSetRKey = node.Symbol + ":SELL"
	}
}

func (node *OrderNode) SetDepthHashKey() {
	node.OrderDepthHashKey = node.Symbol + ":depth"
	priceStr := decimal.NewFromFloat(node.Price).String()
	node.OrderDepthHashField = node.Symbol + ":depth:" + priceStr
}

func (node *OrderNode) SetNodeName() {
	node.NodeName = node.Symbol + ":node:" + node.Oid
}

func (node *OrderNode) SetNodeLink() {
	priceStr := decimal.NewFromFloat(node.Price).String()
	node.NodeLink = node.Symbol + ":link:" + priceStr
}
