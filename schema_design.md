# GraphQL Schema Design

Fields should often do one thing, and do it really well. A good indication that a field might be trying to do more than one thing is a boolean argument.

## Cursor pagination
cursor pagination is a good choice for GraphQL APIs. Today, most GraphQL APIs use cursor-based pagination, and that is mostly due to Relay’s connection pattern

## Relay Connections

```
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