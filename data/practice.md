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


## Resource
* [GraphQL Playground](https://www.youtube.com/watch?v=CHNAnGSmQeA)