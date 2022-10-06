# GraphQL Parser
GraphQL has been released only as a specification. This means that GraphQL is in fact not more than a long document that describes in detail the behaviour of a GraphQL server.

## Schema 
> Schema contains all information about what a client can potentially do with a GraphQL API.

> Every GraphQL service defines a set of types which completely describe the set of possible data you can query on that service. Then, when queries come in, they are validated and executed against that schema.


## Arguments
Every field on a GraphQL `object type` can have zero or more arguments, for example the length field below:
```
type Starship {
  id: ID!
  name: String!
  length(unit: LengthUnit = METER): Float
}
```
Arguments can be either required or optional. When an argument is optional, we can define a default value - if the `unit` argument is not passed, it will be set to `METER` by default.

### Every GraphQL service has a query type and may or may not have a mutation type. 

### Default scalar types
* `Int`: A signed 32‐bit integer.
* `Float`: A signed double-precision floating-point value.
* `String`: A UTF‐8 character sequence.
* `Boolean`: `true` or `false`.
* `ID`: The ID scalar type represents a unique identifier


## Lists and Non-Null
Here, we're using a `String` type and marking it as Non-Null by adding an exclamation mark, `!` after the type name.

> `myField: [String!]`

### Object types, scalars, and enums are the only kinds of types you can define in GraphQL.

## GraphQL Client Libraries
GraphQL is particularly great for frontend developers since it completely eliminates many of the inconveniences and shortcomings that are experienced with REST APIs, such as over- and underfetching. Complexity is pushed to the server-side where powerful machines can take care of the heavy computation work. The client doesn't have to know where the data that it fetches is actually coming from and can use a single, coherent and flexible API.

> go get -u github.com/mateors/graphqlparser


## What is a resolver function?
A function on a GraphQL server that's responsible for fetching the data for a single field

## Resource
* https://graphql.org/learn/schema
* https://graphql.org/learn/queries
* https://graphql.org/learn/introspection
* https://www.howtographql.com
* https://www.howtographql.com/basics/2-core-concepts
* https://www.howtographql.com/basics/3-big-picture