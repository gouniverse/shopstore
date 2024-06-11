package shopstore

// type LineItem struct {
// 	ID       string
// 	OrdeID   string
// 	Name     string
// 	Price    float64
// 	Quantity int64
// }

// func (order *Order) LineItemAdd(lineItem LineItem) {
// 	order.lineItems = append(order.lineItems, lineItem)
// }

// func (order *Order) LineItemList() []LineItem {
// 	return order.lineItems
// }

// func (order *Order) LineItemRemove(lineItemId string) error {

// 	index := order.findLineItemIndex(lineItemId)
// 	if index != -1 {
// 		// As I'd like to keep the items ordered, in Golang we have to shift all of the elements at
// 		// the right of the one being deleted, by one to the left.
// 		order.lineItems = append(order.lineItems[:index], order.lineItems[index+1:]...)
// 	}

// 	return nil
// }

// func (order *Order) LineItemsRemoveAll() {
// 	order.lineItems = []LineItem{}
// }

// func (order *Order) findLineItemIndex(lineItemId string) int {
// 	for index, lineItem := range order.lineItems {
// 		if lineItemId == lineItem.ID {
// 			return index
// 		}
// 	}

// 	return -1
// }
