# shopstore <a href="https://gitpod.io/#https://github.com/gouniverse/shopstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

[![Tests Status](https://github.com/gouniverse/shopstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/gouniverse/shopstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/shopstore)](https://goreportcard.com/report/github.com/gouniverse/shopstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/shopstore)](https://pkg.go.dev/github.com/gouniverse/shopstore)


## License

This project is licensed under the GNU General Public License version 3 (GPL-3.0). You can find a copy of the license at https://www.gnu.org/licenses/gpl-3.0.en.html

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.


## Installation

```ssh
go get -u github.com/gouniverse/shopstore
```

## Usage

```golang
ShopStore, err := shopstore.NewStore(shopstore.NewStoreOptions{
  DB:                     Database.DB(),
  CategoryTableName:      "shop_category",
  DiscountTableName:      "shop_discount",
  MediaTableName:         "shop_media",
  OrderTableName:         "shop_order",
  OrderLineItemTableName: "shop_order_line_item",
  ProductTableName:       "shop_product",
  AutomigrateEnabled: true,
})

if err != nil {
  panic(err)
}

if ShopStore == nil {
  panic("ShopStore is nil")
}
```
