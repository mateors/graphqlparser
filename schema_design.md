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




A static query is a query that does not change based on any variable, condition, or state of the program.

### What do we mean by transactions and GraphQL?
the current best practice is to look at designing these errors as part of our schema rather than treating them as exceptions/query level errors

A possibly better approach is for payload types to include something like a userErrors field:

```go
type SignUpPayload {
    userErrors: [UserError!]!
    account: Account
}

type UserError {
    # The error message
    message: String!
    # Indicates which field cause the error, if any
    #
    # Field is an array that acts as a path to the error
    #
    # Example:
    #
    # ["accounts", "1", "email"]
    #
    field: [String!]
    # An optional error code for clients to match on.
    code: UserErrorCode
}
```



## Union / Result Types
> errors as data

instead of using a specific field for errors, this approach uses union types to represent possible problematic states to the client. Let’s take the same `sign up` example, but design it using a result union:

```go
type Mutation {
    signUp(email: string!, password: String!): SignUpPayload
}
union SignUpPayload =
                        SignUpSuccess |
                        UserNameTaken |
                        PasswordTooWeak

mutation {
    signUp(
        email: "marc@productionreadygraphql.com",
        password: "P@ssword"
    ){
        ... on SignUpSuccess {
            account{
                id
            }
        }
        ... on UserNameTaken {
            message
            suggestedUsername
        }
        ... on PasswordTooWeak {
            message
            passwordRules
        }
     }
}
```
