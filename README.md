# GraphQL Parser | A Query Language for APIs
GraphQL has been released only as a specification. This means that GraphQL is in fact not more than a long document that describes in detail the behaviour of a GraphQL server.


* GraphQL is a query language for APIs - not databases.
* A more efficient Alternative to REST
* API defines how a client can load data from a server.
* GraphQL uses the concept of resolver functions to collect the data that's requested by a client.
* GraphQL APIs typically only expose a single endpoint

One of the major advantages of GraphQL is that it allows for naturally querying nested information. 


> At its core, GraphQL enables `declarative data fetching` where a client can specify exactly what data it needs from an API. Instead of multiple endpoints that return fixed data structures, a GraphQL server only exposes a single endpoint and responds with precisely the data a client asked for.


```
QueryDocument()
Operation
Query
Mutation
Subscription

parseField
parseArguments
parseFragment
parseObject
```

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

## AST

### Document Syntax
Document Syntax listed using [GraphQL Spec](https://spec.graphql.org/October2021/#index)

```
Document::
	Definition

Definition::
	OperationDefinition ->ExecutableDefinition
	FragmentDefinition -> ExecutableDefinition

	#TypeSystemDefinition -> TypeSystemDefinitionOrExtension
	SchemaDefinition ->TypeSystemDefinition
	TypeDefinition -> TypeSystemDefinition
		ScalarTypeDefinition
		ObjectTypeDefinition
		InterfaceTypeDefinition
		UnionTypeDefinition
		EnumTypeDefinition
		InputObjectTypeDefinition

	DirectiveDefinition -> TypeSystemDefinition

	#TypeSystemExtension -> TypeSystemDefinitionOrExtension
	SchemaExtension -> TypeSystemExtension
	TypeExtension -> TypeSystemExtension


####################################################################################
OperationDefinition::
	OperationType Name[opt] VariablesDefinition[opt] Directives[opt] SelectionSet
	SelectionSet

SelectionSet:
{ Selection[list] }

Selection:
	Field
	FragmentSpread
	InlineFragment

Field:
Alias[opt] Name Arguments[opt] Directives[opt] SelectionSet[opt]

Argument:
Name : Value

FragmentSpread:
... FragmentName Directives[opt]

InlineFragment:
... TypeCondition[opt] Directives[opt] SelectionSet

TypeCondition:
on NamedType

VariableDefinition:
Variable : Type DefaultValue[opt] Directives[opt]

DirectiveDefinition:
Description[opt] directive @ Name ArgumentsDefinition[opt] repeatable[opt] on DirectiveLocations

ArgumentsDefinition:
( InputValueDefinition[list] )

InputValueDefinition:
Description[opt] Name : Type DefaultValue[opt] Directives[opt]


Type:
	NamedType
	ListType
	NonNullType
-----------------------------------

FragmentDefinition::
	fragment FragmentName TypeCondition Directives[opt] SelectionSet

SchemaDefinition::
	Description[opt] schema Directives[opt] { RootOperationTypeDefinition[list] }

RootOperationTypeDefinition: 
    OperationTypeDefinition:
OperationType : NamedType

ScalarTypeDefinition::
Description[opt] scalar Name Directives[opt]


ObjectTypeDefinition::
Description[opt] type Name ImplementsInterfaces[opt] Directives[opt] { FieldsDefinition }
Description[opt] type Name ImplementsInterfaces[opt] Directives[opt]


FieldDefinition::
Description[opt] Name ArgumentsDefinition[opt] : Type Directives[opt]

InterfaceTypeDefinition::
Description[opt] interface Name ImplementsInterfaces[opt] Directives[opt] FieldsDefinition
Description[opt] interface Name ImplementsInterfaces[opt] Directives[opt]


UnionTypeDefinition::
Description[opt] union Name Directives[opt] UnionMemberTypes[opt]

example: union SearchResult = Lift | Trail

EnumTypeDefinition::
Description[opt] enum Name Directives[opt] EnumValuesDefinition
Description[opt] enum Name Directives[opt]

EnumValueDefinition:
Description[opt] EnumValue Directives[opt]

example:
enum Direction {
  NORTH
  EAST
  SOUTH
  WEST
}


InputObjectTypeDefinition::
Description[opt] input Name Directives[opt] InputFieldsDefinition
Description[opt] input Name Directives[opt]


Variable:
$ Name

DefaultValue:
= Value


Value:
	Variable
	IntValue
	FloatValue
	StringValue
	BooleanValue
	NullValue
	EnumValue
	ListValue
	ObjectValue

Name:
 NameStart Letter Digit
####################################################################################

RootOperationTypeDefinition::
 OperationType : NamedType


Type::
	1. NamedType
	2. ListType
	3. NonNullType

NamedType
	Name

ListType
	[Type]

NonNullType
	NamedType!
	ListType!


QUERY
MUTATION
SUBSCRIPTION
FIELD
FRAGMENT_DEFINITION
FRAGMENT_SPREAD
INLINE_FRAGMENT
VARIABLE_DEFINITION

SCHEMA
SCALAR
OBJECT
FIELD_DEFINITION
ARGUMENT_DEFINITION
INTERFACE
UNION
ENUM
ENUM_VALUE
INPUT_OBJECT
INPUT_FIELD_DEFINITION
-

OperationTypeDefinition
ScalarDefinition

extend scalar Name Directives

extend type Name ImplementsInterfaces[opt] Directives [Cons][opt] FieldsDefinition

```

### ObjectType

* Document
	* Definition
		2. TypeSystemDefinitionOrExtension
			1. TypeSystemDefinition

### TypeSystemDefinition
1. SchemaDefinition
2. **TypeDefinition**
3. DirectiveDefinition

### TypeDefinition
1. ScalarTypeDefinition
2. **ObjectTypeDefinition**
3. InterfaceTypeDefinition
4. UnionTypeDefinition
5. EnumTypeDefinition
6. InputObjectTypeDefinition

### 2. ObjectTypeDefinition
1. Description`opt` **type** Name ImplementsInterfaces`opt` Directives`opt` FieldsDefinition
2. Description`opt` **type** Name ImplementsInterfaces`opt` Directives`opt`


### ImplementsInterfaces
1. ImplementsInterfaces & NamedType
2. **implements** &`opt` NamedType

### FieldsDefinition
`{` FieldDefinition`list` `}`

### FieldDefinition
Description`opt` Name ArgumentsDefinition `opt` `:` Type Directives`opt`

### Type
1. NamedType
2. ListType
3. NonNullType

### NamedType
* Name -> NameStart/Letter + NameContinue/Letter+Digit

### ListType
* `[` Type `]`

### NonNullType
* NamedType`!`
* ListType`!`

### ArgumentsDefinition
> `(` InputValueDefinition`list` `)`

### InputValueDefinition
Description`opt` Name `:` Type DefaultValue`opt` Directives`opt`

### Directives
Directive`list`

### Directive
`@`Name Arguments`opt`

### Arguments
`(` Argument`list` `)`

### Argument
Name `:` Value

### Value
1. Variable
2. IntValue
3. FloatValue
4. StringValue
5. BooleanValue
6. NullValue
7. EnumValue
8. ListValue
9. **ObjectValue**

### ObjectValue
* `{` `}`
* `{` ObjectField`list` `}`

### ObjectField
Name `:` Value

```
type Lift {
  id: ID!
  name: String!
}
```

## Type Extensions

### TypeSystemExtension:
1. SchemaExtension
2. TypeExtension

### TypeExtension::
1. ScalarTypeExtension
2. ObjectTypeExtension
3. InterfaceTypeExtension
4. UnionTypeExtension
5. EnumTypeExtension
6. InputObjectTypeExtension

### 1. ScalarTypeExtension:
> `extend scalar Name Directives`

### 2. ObjectTypeExtension:
* `extend type Name ImplementsInterfaces[opt] Directives[opt] FieldsDefinition`
* `extend type Name ImplementsInterfaces[opt] Directives`
* `extend type Name ImplementsInterfaces`

### 3. InterfaceTypeExtension:
* `extend interface Name ImplementsInterfaces[opt] Directives[opt] FieldsDefinition`
* `extend interface Name ImplementsInterfaces[opt] Directives`
* `extend interface Name ImplementsInterfaces`

### 4. UnionTypeExtension:
* `extend union Name Directives[opt] UnionMemberTypes`
* `extend union Name Directives`

### 5. EnumTypeExtension:
* `extend enum Name Directives[opt] EnumValuesDefinition`
* `extend enum Name Directives`

### 6. InputObjectTypeExtension:
* `extend input Name Directives[opt] InputFieldsDefinition`
* `extend input Name Directives`

## ObjectDefinition Input
```
type Person {
  id: ID!
  adult: Boolean!
}
```
### Output steps / Manual tracing:
* parseDocument> 29
* parseObjectDefinition->START {29 2 2 6 type}
* parseDescription {29 2 2 6 type} {20 2 7 13 Person}
* expectToken 29
* fieldDefinition {20 3 18 20 id}
* fieldDefinition {20 4 28 33 adult}
* parseObjectDefinition->DONE
```
type Person {
id: ID!
adult: Boolean!
}
```

## Tips for parsing ast.Node
* Take closer look in your token/lexer package
* Parsing feels like you are walking in a street.
* You should not/or avoid useing builtin package/libarary
* Every parser method should have starting with a validation check and `return nil`
* Optional fields parser method must have a starting validation checker which `return nil` if wrong
* validation checker is a current token checker, ex: `p.curTokenIs(token.IDENT)`
* If any error or bug make sure token/lexer package producing the correct tokens.
* Make sure a break point for loop `for { }`.

### ObjectTypeDefinition Input
```go
	"""Object definition"""
	type Person {
		"Description for id" id: ID!
		"Description for age" age: Int
		length("Yes" unit: LengthUnit = METER, "No" corner: Int = 50): Float
	}
```

### ObjectTypeDefinition Parser Output
```go
"""
Object definition
"""
type Person {
"Description for id"
id: ID!
"Description for age"
age: Int
length("Yes" unit: LengthUnit = METER, "No" corner: Int = 50): Float
}
```

* git tag
* git tag -a v0.0 -m "parseObjectDefinition complete"
* git show v0.0
* git push origin v0.0


### Return nil in the following function if any error
* parseName
* parseNamed
* parseType

### Parse Hierarchy
```
parseDocument::
	parseObjectDefinition::
		parseFieldsDefinition
			parseFieldDefinition
				parseArgumentDefinition
					parseInputValueDefinition
```
### TypeSystemDefinition
* SchemaDefinition
* TypeDefinition -> DONE
* DirectiveDefinition

### TypeDefinition
* ScalarTypeDefinition -> DONE
* ObjectTypeDefinition -> DONE
* InterfaceTypeDefinition -> DONE
* UnionTypeDefinition -> DONE
* EnumTypeDefinition --> DONE
* InputObjectTypeDefinition --> DONE

### ExecutableDefinition
* OperationDefinition -> DONE
* FragmentDefinition -> DONE

### SelectionSet
 * Selection
	1. Field -> DONE
	2. FragmentSpread -> DONE
	3. InlineFragment -> DONE

### Query shorthand -> Done

## Resource
* https://graphql.org/learn/schema
* https://graphql.org/learn/queries
* https://graphql.org/learn/introspection
* https://www.howtographql.com
* https://www.howtographql.com/basics/2-core-concepts
* https://www.howtographql.com/basics/3-big-picture
* https://www.apollographql.com/docs/apollo-server/schema/directives
* http://snowtooth.herokuapp.com