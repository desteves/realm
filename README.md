# mongodb realm

Contains [mongoDB Realm](https://stitch.mongodb.com/) (formerly known as Stitch) [Authentication](pkg/auth) and a [GraphQL Client](pkg/graphql) go packages.

See [examples/](examples/)

## Usage 

```bash
go get https://github.com/desteves/realm
import  "github.com/desteves/realm/pkg/graphql"
```


## TODO (work in progress)

- Test Refresh Token with oauth package
- Test all auth mechanisms, so far only test anonymous login.
- Add `*_test` everywhere
- Implement `Disconnect` for auth packages.


## Atlas and Stitch Set up

- Create new project under an organization.
- Create new free cluster, call it `graphqlDatabase`.
- Load 'Sample Data Set'
- Create new stitch app, call it `graphqlServer`.   
  - Note down app id `graphqlserver-?????`
- Enable anonymous authentication
- Enable GraphQL
- Generate 'stitch schemas'
- Select any of the `sample-*` databases/collections and click 'Generate Schema'. 
  - Enable read/write access

### (Optional)  Set up webhook

This is an additional check to verify our client has network connectivty.

Follow the steps [listed here](https://docs.mongodb.com/stitch/reference/service-webhooks/#creating-a-webhook)

- Name the service as `ping`
- Name the webhook as `test` 
- Only allow the `GET` HTTP Method
- Paste the following in the function
```javascript
  // This function is the webhook's request handler.
  exports = function(payload, response) {
    response.setStatusCode(200);
    response.setBody(
    "{'message': 'pong'}"
  );
  };
```

See [examples/anonymousauth/main.go](examples/anonymousauth/main.go) for usage.


### (Optional)  Set up dummy graphql record

This is to verify our client can reach the graphql server and obtain a response.

Insert the following to your atlas cluster. This example is from the mongo shell:

```javascript

var appid = "graphqlserver-?????" // <------- UPDATE THIS!!!
use graphql
db.health.insertOne({
 "appid": appid,
 "description": "Checks GraphQL Server is reachable and operational",
 "status": "pass",
 "endpoint": "https://stitch.mongodb.com/api/client/v2.0/app/"+appid+"/graphql"
})

```

See [examples/graphqlhealthcheck/main.go](examples/graphqlhealthcheck/main.go) for usage.


## References

- [Mongo Go Driver](https://github.com/mongodb/mongo-go-driver)
- [MongoDB GraphQL Docs](https://docs-mongodbcom-staging.corp.mongodb.com/stitch/nick/graphql/graphql.html)