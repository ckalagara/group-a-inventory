# group-a-inventory

## Setup

### Docker
```
docker network create app-network
docker run --network app-network --name mongodb -p 27017:27017 -d mongodb/mongodb-community-server:latest
docker run --network app-network -p 50052:50052 group-a-inventory
```

### protoc
```
 protoc --go_out=. --go-grpc_out=. proto/inventory.proto
```

### Grpcurl

```
Last login: Wed Sep  3 19:48:14 on ttys004
i@narwhal group-a-manager % grpcurl -d '{
  "item": {
    "id": "Apple_laptop_mackbookprom3",
    "name": "Laptop",
    "description": "A high-performance laptop",
    "quantity": 10
  }
}' -plaintext localhost:50052 inventory.Service/AddItem
{
  "item": {
    "id": "Apple_laptop_mackbookprom3",
    "name": "Laptop",
    "description": "A high-performance laptop",
    "quantity": 10
  }
}
i@narwhal group-a-manager % grpcurl -d '{
  "item": {
    "id": "Apple_laptop_mackbookairm3",
    "name": "Laptop",
    "description": "A lightweight laptop",
    "quantity": 10
  }
}' -plaintext localhost:50052 inventory.Service/AddItem
{
  "item": {
    "id": "Apple_laptop_mackbookairm3",
    "name": "Laptop",
    "description": "A lightweight laptop",
    "quantity": 10
  }
}
i@narwhal group-a-manager % grpcurl -d '{}' -plaintext localhost:50052 inventory.Service/ListItems
{
  "items": [
    {
      "id": "Apple_laptop_mackbookprom3",
      "name": "Laptop",
      "description": "A high-performance laptop",
      "quantity": 10
    },
    {
      "id": "Apple_laptop_mackbookairm3",
      "name": "Laptop",
      "description": "A lightweight laptop",
      "quantity": 10
    }
  ]
}

```