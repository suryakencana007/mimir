# Mimir

<p align="center"><img src="doc/mimir.jpg" width="360"></p>

* [Description](#description)
* [Installation](#installation)
* [Usage](#usage)
* [Documentation](#documentation)
* [FAQ](#faq)
* [License](#license)

## Description

Mimir is an Package Library tools that helps your application at scale. With Mimir, you can:

- Breaker
- Logging
- Sql Pagination
- AWS lambda handler for http go standard  


## Installation
```
go get -u github.com/suryakencana007/mimir
```

## Usage


- Use a Breaker

```go
package main
import (
    "net/http"
    "github.com/hashicorp/go-cleanhttp"
    "github.com/suryakencana007/mimir/breaker"
    "github.com/suryakencana007/mimir/log"
)

func main() {
    endpoint := `www.google.com`
    cb := breaker.NewBreaker(
        "",
        100,
        10,
    )
    var res *http.Response
    err := cb.Execute(func() (err error) {
        client := cleanhttp.DefaultClient()
        req, _ := http.NewRequest(http.MethodGet,
            endpoint, nil)
        res, err = client.Do(req)
        return err
    })
    if err != nil {
        log.Error("Error",
            log.Field("error", err.Error()),
        )
        panic(err)
    }
}
```

- Use a Log

```go
package main
import (
    "github.com/suryakencana007/mimir/log"
)

func main(){
    log.ZapInit()
    
    log.Error("Error",
        log.Field("error", "error for log"),
    )
    log.Debug("Debug",
        log.Field("debug", "debug for log"),
        log.Field("message", "message debug for log"),
    )
    log.Info("Info",
        log.Field("info", "info for log"),
        log.Field("message", "message info for log"),
    )
    log.Warn("Warn",
        log.Field("warning", "warning for log"),
        log.Field("message", "message warning for log"),
    )
}

``` 

- Use sql Pagination

```go
package main

import (
    "net/http"
    "github.com/suryakencana007/mimir/sql"
)

func main() {
    handler := func() http.HandlerFunc { 
        return func(w http.ResponseWriter, r *http.Request) {
            pagination := &sql.Pagination{
                Params: r.URL.Query(),
            }
            orders := service.All(pagination)
            body := response.NewResponse()
            body.SetData(orders)
            body.Pagination = &response.Pagination{
                Page:  pagination.Page,
                Size:  pagination.Limit,
                Total: pagination.Total,
            }
            response.WriteJSON(w, r, body)
        }
    }
    http.ListenAndServe(":8080", handler())
}
``` 

- AWS lambda http handler

```go
package main

import (
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/suryakencana007/mimir/request"
)

func main() {
    handler := func() http.HandlerFunc { 
        return func(w http.ResponseWriter, r *http.Request) {
            pagination := &sql.Pagination{
                Params: r.URL.Query(),
            }
            orders := service.All(pagination)
            body := response.NewResponse()
            body.SetData(orders)
            body.Pagination = &response.Pagination{
                Page:  pagination.Page,
                Size:  pagination.Limit,
                Total: pagination.Total,
            }
            response.WriteJSON(w, r, body)
        }
    }
    lambda.Start(
        request.HandleEvent(
            handler(),
        ),
    )
}
```

### Importing the package

This package can be used by adding the following import statement to your `.go` files.

```go
    import "github.com/suryakencana007/mimir"
```

## FAQ

[Please do!](https://github.com/suryakencana007/mimir/blob/master/CONTRIBUTING.md) We are looking for any kind of contribution to improve and add more library function helper Mimir core funtionality and documentation. When in doubt, make a PR!

## License
 
 ```
 Copyright 2019, Nanang Suryadi (https://nulis.dev)
 
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
 
     http://www.apache.org/licenses/LICENSE-2.0
 
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
 ```
 
