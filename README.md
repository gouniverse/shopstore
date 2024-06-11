# shopstore <a href="https://gitpod.io/#https://github.com/gouniverse/shopstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

[![Tests Status](https://github.com/gouniverse/shopstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/gouniverse/shopstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/shopstore)](https://goreportcard.com/report/github.com/gouniverse/shopstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/shopstore)](https://pkg.go.dev/github.com/gouniverse/shopstore)


## Installation

```ssh
go get -u github.com/gouniverse/shopstore
```

## Usage

```golang
ShopStore, err := shopstore.NewStore(shopstore.NewStoreOptions{
  DB:                 Database.DB(),
  DiscountTableName:  "shop_discount",
  OrderTableName:     "shop_order",
  AutomigrateEnabled: true,
})

if err != nil {
  panic(err)
}

if ShopStore == nil {
  panic("ShopStore is nil")
}
```
