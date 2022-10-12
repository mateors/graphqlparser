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

### Data paging
We can use GraphQL arguments to control the amount of data that is returned from our queries. This process is called data paging because a specific number of records are returned to represent one page of data.

```js
type Query {
  ...
  allUsers(first: Int=50 start: Int=0): [User!]!
  allPhotos(first: Int=25 start: Int=0): [Photo!]!
}
```
In the preceding example, we have added optional arguments for first and start. If the client does not supply these arguments with the query, we will use the default values provided. By default, the allUsers query returns only the first 50 users, and the allPhotos query returns only the first 25 photos.

```js
query {
  allUsers(first: 10 start: 90) {
    name
    avatar
  }
}
```

### Sorting
When querying a list of data, we might also want to define how the returned list of data should be sorted. We can use arguments for this, as well.

```js
enum SortDirection {
ASCENDING
DESCENDING
} 

enum SortablePhotoField {
  name
  description
  category
  created
} 

Query {
  allPhotos(
    sort: SortDirection = DESCENDING
    sortBy: SortablePhotoField = created
  ): [Photo!]!

}
```

Clients can now control how their photos are sorted when they issue an allPhotos query:
```js
query {
  allPhotos(sortBy: name)
}
```

So far, we've added arguments only to fields of the Query type, but it is important to note that you can add arguments to any field.

```js
type User {

  postedPhotos(
    first: Int = 25
    start: Int = 0
    sort: SortDirection = DESCENDING
    sortBy: SortablePhotoField = created
    category: PhotoCategory
  ): [Photo!]

}

```

## Mutations
Mutations must be defined in the schema. Just like queries, mutations also are defined in their own custom object type and added to the schema. Technically, there is no difference between how a mutation or query is defined in your schema. The difference is in intent. 

> We should create mutations only when an action or event will change something about the state of our application.

When designing your GraphQL service, make a list of all of the actions that a user can take with your application. Those are most likely your mutations.

```js
type Mutation {
  postPhoto(
  name: String!
  description: String
  category: PhotoCategory=PORTRAIT
  ): Photo!
} 

schema {
  query: Query
  mutation: Mutation
}
```

a user can post a photo by sending the following mutation:
```js
mutation {
  postPhoto(name: "Sending the Palisades") {
      id
      url
      created
      postedBy {
        name
      }
  }
}

mutation postPhoto(
  $name: String!
  $description: String
  $category: PhotoCategory
) 

{
  postPhoto(name: $name,  description: $description,  category: $category) {
    id
    name
    email
  }
}
```

### Input Types
As you might have noticed, the arguments for a couple of our queries and mutations are getting quite lengthy. There is a better way to organize these arguments using input types. An input type is similar to the GraphQL object type except it is used only for input arguments.

Let's improve the postPhoto mutation using an input type for our arguments:
```js
input PostPhotoInput {
  name: String!
  description: String
  category: PhotoCategory=PORTRAIT
} 

type Mutation {
  postPhoto(input: PostPhotoInput!): Photo!
}
```

```js
mutation newPhoto($input: PostPhotoInput!) {
  postPhoto(input: $input) {
      id
      url
      created
  }
}
```

When we send the mutation, we need to supply the new photo data in our query variables nested under the input field:
```json
{
  "input": {
  "name": "Hanging at the Arc",
  "description": "Sunny on the deck of the Arc",
  "category": "LANDSCAPE"
  }
}
```
Our input is grouped together in a JSON object and sent along with the mutation in the query variables under the “input” key. 

Because the query variables are formatted as JSON, the category needs to be a string that matches one of the categories from the PhotoCategory type.

Input types help us organize our schema and reuse arguments. They also improve the schema documentation that GraphiQL or GraphQL Playground automatically generates.

### Return Types
```
type AuthPayload {
  user: User!
  token: String!
} 

type Mutation {
  ...
  githubAuth(code: String!): AuthPayload!
}
```
You can use custom return types on any field for which we need more than simple payload data. 

## Resource
* [GraphQL Playground](https://www.youtube.com/watch?v=CHNAnGSmQeA)
* https://spec.graphql.org/October2016/#index
* [GraphQL Document-Syntax](https://spec.graphql.org/draft/#sec-Document-Syntax)