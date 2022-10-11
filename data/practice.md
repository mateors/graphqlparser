# GraphQL Practice with Playground

## Playground API Endpoint
* http://snowtooth.herokuapp.com
* http://snowtooth.moonhighway.com



```go
mutation GiveAnyName{ 
  setLiftStatus(id: "astra-express"	status: OPEN)
  {
    id
    name
  }
}

subscription SubscriptionName{
  liftStatusChange{
    id
    name
    status
  }
}


mutation GiveAnyName{ 
  m1:setLiftStatus(id: "astra-express"	status: OPEN)
  {
    id
    name
  }
  
  m2:setLiftStatus(id: "jazz-cat"	status: CLOSED)
  {
    id
    name
  }
}
```


## Introspection
One of the most powerful features of GraphQL is introspection. Introspection is the ability to query details about the current API's schema. Introspection is how those nifty ছিমছাম GraphQL documents are added to the GraphiQL Playground interface

```
query Introspection{
  __schema {
    types {
    name
    description
    }
  }
}
```
When we run this query, we see every type available on the API, including root types, custom types, and even scalar types.


If we want to see the details of a particular type, we can run the __type query and send the name of the type that
we want to query as an argument:

```js
query liftDetails {
  __type(name:"Lift") {
    name
    fields {
      name
      description
      type {
        name
        description
      }
    }
  }
}
```

### What fields are available on the root types:
```js
query roots {
  __schema {
    queryType {
      ...typeFields
    }
    mutationType {
      ...typeFields
    }
    subscriptionType {
      ...typeFields
    }
  }
} 

fragment typeFields on __Type {
  name
  fields {
  	name
  }
}
```
The redundancy of the preceding query has been reduced by using a fragment.


## Abstract Syntax Trees - AST
An abstract syntax tree, or AST, is a hierarchical object that represents our query. The AST is an object that contains nested fields that represent the details of a GraphQL query.

An `object hierarchy` is a concept from computer programming. It references descendants of objects acting as properties of an object.


A document contains at least one definition, but it can also contain a list of definitions.

Definitions are only one of two types: OperationDefinition or FragmentDefinition. 

## Resource
* [GraphQL Playground](https://www.youtube.com/watch?v=CHNAnGSmQeA)
* https://spec.graphql.org/October2016/#index
* [GraphQL Document-Syntax](https://spec.graphql.org/draft/#sec-Document-Syntax)