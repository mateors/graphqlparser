# GraphQL Parser | A Query Language for APIs
GraphQL has been released only as a specification. This means that GraphQL is in fact not more than a long document that describes in detail the behaviour of a GraphQL server.


GraphQL is a query language for APIs - not databases.

A more efficient Alternative to REST

API defines how a client can load data from a server.

GraphQL uses the concept of resolver functions to collect the data that's requested by a client.

GraphQL APIs typically only expose a single endpoint

One of the major advantages of GraphQL is that it allows for naturally querying nested information. 


At its core, GraphQL enables declarative data fetching where a client can specify exactly what data it needs from an API. Instead of multiple endpoints that return fixed data structures, a GraphQL server only exposes a single endpoint and responds with precisely the data a client asked for.

## Mutations
> Making changes to the data that’s currently stored in the backend. With GraphQL, these changes are made using so-called mutations. There generally are three kinds of mutations:

1. creating new data
2. updating existing data
3. deleting existing data

Mutations follow the same `syntactical structure` as queries, but they always need to start with the `mutation keyword`.

```
mutation {
  createPerson(name: "Mostain", age: 36) {
    name
    age
  }
}
```

## Schema
It specifies the capabilities of the API and defines how clients can request the data. It is often seen as a contract between the server and client.

Generally, a `schema is` simply a collection of `GraphQL types`. However, when writing the schema for an API, there are some special root types:

```
type Query { ... }
type Mutation { ... }
type Subscription { ... }
```

## Benefits of a Schema & Type System
GraphQL uses a strong type system to define the capabilities of an API. All the types that are exposed in an API are written down in a schema using the GraphQL Schema Definition Language (SDL). This schema serves as the contract between the client and the server to define how a client can access the data.

## SDL - Schema Definition Language
GraphQL has its own type system that’s used to define the schema of an API. The syntax for writing schemas is called Schema Definition Language (SDL).

## Schema 
> Schema contains all information about what a client can potentially do with a GraphQL API.

> Every GraphQL service defines a set of types which completely describe the set of possible data you can query on that service. Then, when queries come in, they are validated and executed against that schema.

* queries and mutations
* schema 
* query
* type system
* field

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

```
type Character {
  name: String!
  appearsIn: [Episode]!
}
```

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