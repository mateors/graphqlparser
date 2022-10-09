# GraphQL Schema Design

Fields should often do one thing, and do it really well. A good indication that a field might be trying to do more than one thing is a boolean argument.

## Cursor pagination
cursor pagination is a good choice for GraphQL APIs. Today, most GraphQL APIs use cursor-based pagination, and that is mostly due to Relay’s connection pattern

## Relay Connections

```go
type ProductConnection {
    edges: [ProductEdge]
    pageInfo: PageInfo!
}

type ProductEdge {
    cursor: String!
    node: Product!
}

type PageInfo {
    endCursor: String
    hasNextPage: Boolean!
    hasPreviousPage: Boolean!
    startCursor: String
}

type Product {
    name: String!
}
```

If your use case can be supported by `cursor-based pagination`, I highly recommend choosing the `connection pattern` when designing your GraphQL API:

* Cursor pagination is generally a great choice for accuracy and performance.
* It lets Relay clients seamlessly integrate with your API.
* It is probably the most common pattern in GraphQL at the moment and lets us be consistent with other APIs in the space.
* The connection pattern lets us design more complex use cases, thanks to the Connection and Edge types.







