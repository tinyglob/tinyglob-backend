### Production URL

https://tinyglob-backend-production.up.railway.app/

### Tech

Serverless Postgres - https://neon.tech/ <br />
Deployments & Scale - https://railway.app/ <br />
Back-end Language - https://go.dev/ <br />

### Quick Start

```
railway login
```

```
railway link
```

```
railway run go run main.go
```

``` (alternative)
make run prod
```

```
curl http://localhost:8080/
```

### APIs

`GET` `/` <br/>
`GET` `/jobs` <br/>
`GET` `/jobs/continent/{continent}` <br/>
`GET` `/jobs/country/{country}` <br/>
`GET` `/jobs/id/{id}` <br/>
