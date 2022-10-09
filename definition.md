# GraphQL Terminology

Clients use the GraphQL query language to make requests to a GraphQL service.

We refer to these request sources as `documents`. 

## Document
A document may contain operations (queries, mutations, and subscriptions) as well as fragments, a common unit of composition allowing for query reuse.

## NAMES
GraphQL Documents are full of named things:
operations, fields, arguments, types, directives, fragments, and variables.
Names in GraphQL are case‐sensitive. 

A document contains multiple definitions, either executable or representative of a GraphQL type system.

Documents are only executable by a GraphQL service if they contain an OperationDefinition and otherwise only contain ExecutableDefinition

OperationDefinition:
OperationType Name[opt] VariableDefinitions[opt] Directives[opt] SelectionSet

SelectionSet:
{ Selection}

Selection:
    Field
    FragmentSpread
    InlineFragment


## Grammar Notation

### Optionality and Lists
A subscript suffix Symbol ~list~ is shorthand for a list of one or more of that symbol, represented as an additional recursive production.

