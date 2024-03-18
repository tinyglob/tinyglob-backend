### URL

https://tinyglob-backend-production.up.railway.app/

### Tech

Serverless Postgres - https://neon.tech/ <br />
Deployments & Scale - https://railway.app/ <br />
Back-end Language - https://go.dev/ <br />

### API

```console
router.Get("/", getRootHandler)
router.Get("/jobs", getJobCountByContinentHandler)
router.Get("/jobs/continent/{continent}", getJobsByContinentHandler)
router.Get("/jobs/id/{id}", getJobByIDHandler)
```
