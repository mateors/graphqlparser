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


### Union Type
In GraphQL, a union type is a type that we can use to return one of several different types.

```js
query schedule {

 agenda {
  ...on Workout {
      name
      reps
  }
  ...on StudyGroup {
      name
      subject
      students
  }
 }
}
```
we could handle this by creating a union type called AgendaItem:

```js
union AgendaItem = StudyGroup | Workout

type StudyGroup {
  name: String!
  subject: String
  students: [User!]!
} 

type Workout {
  name: String!
  reps: Int!
} 

type Query {
  agenda: [AgendaItem!]!
}
```
AgendaItem combines study groups and workouts under a single type. When we add the agenda field to our Query, we are defining it as a list of either workouts or study groups.

It is possible to join as many types as we want under a single union. Simply separate each type with a pipe:
> `union = StudyGroup | Workout | Class | Meal | Meeting | FreeTime`


### Interfaces
Another way of handling fields that could contain multiple types is to use an interface. Interfaces are abstract types that can be implemented by object types.

An interface defines all of the fields that must be included in any object that implements it

```js
scalar DataTime

interface AgendaItem {
  name: String!
  start: DateTime!
  end: DateTime!
} 

type StudyGroup implements AgendaItem {
  name: String!
  start: DateTime!
  end: DateTime!
  participants: [User!]!
  topic: String!
} 

type Workout implements AgendaItem {
  name: String!
  start: DateTime!
  end: DateTime!
  reps: Int!
} 

type Query {
  agenda: [AgendaItem!]!
}
```

In this example, we create an interface called AgendaItem. This interface is anabstract type that other types can implement. When another type implements an interface, it must contain the fields defined by the interface.


### Which one should i use?
Both union types and interfaces are tools that you can use to create fields that contain different object types. It’s up to you to decide when to use one or the other. 

In general, if the objects contain completely different fields, it is a good idea to use union types. They are very effective. If an object type must contain specific fields in order to interface with another type of object, you will need to user an interface rather than a union type.



### Arguments
Arguments can be added to any field in GraphQL. They allow us to send data that can affect outcome of our GraphQL operations.


```js
type Query {
  User(githubLogin: ID!): User!
  Photo(id: ID!): Photo!
}
```
Just like a field, an argument must have a type. That type can be defined using any of the scalar types or object types that are available in our schema.

To select a specific user, we need to send that user’s unique githubLogin as an argument.
```js
query {
  User(githubLogin: "Mostain") {
    name
    avatar
  }
}
```

To select details about an individual photo, we need to supply that photo’s ID:
```js
query {
  Photo(id: "14TH5B6NS4KIG3H4S") {
    name
    description
    url
  }
}
```
In both cases, arguments were required to query details about one specific record. Because these arguments are required, they are defined as non-nullable fields.

> If we do not supply the id or githubLogin with these queries, the GraphQL parser will return an error.


### Filtering Data
Arguments do not need to be non-nullable. We can add optional arguments using nullable fields. This means that we can supply arguments as optional parameters when we execute query operations.

```js
type Query {
  ...
  allPhotos(category: PhotoCategory): [Photo!]!
}
```

if a category is supplied, we should get a filtered list of photos in the same category:
```js
query {
  allPhotos(category: "SELFIE") {
    name
    description
    url
  }
}
```

## Resource
* [GraphQL Playground](https://www.youtube.com/watch?v=CHNAnGSmQeA)
* https://spec.graphql.org/October2016/#index
* [GraphQL Document-Syntax](https://spec.graphql.org/draft/#sec-Document-Syntax)