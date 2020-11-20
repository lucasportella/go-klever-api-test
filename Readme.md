# Klever: Backend Developer

### Challenge

The Technical Challenge consists of creating an API with Golang using gRPC with stream pipes that exposes an upvote service endpoints.

### Techinical Requeriments
  
  - Keep the code in Github
  
  - The API must guarantee the typing of user inputs. If an input is expected as a string, it can only be received as a string.
  
  - The structs used with your mongo model should support Marshal/Unmarshal with bson, json and struct
  
  - The API should contain unit test of methods it uses

### Extra

  - Deliver the whole solution running in some free cloud service

## Postman

[Collection Link](https://www.getpostman.com/collections/3ef276a55fe65e4857f5)


## How to run

1. Clone Repository

```bash
git clone https://github.com/roneycharles/klever
```

2. Run Server

```bash
  make run
```

3. Run Client

```bash
  make run_client
```