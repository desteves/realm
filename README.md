# mongodb realm go package

Contains a [mongoDB Realm](https://stitch.mongodb.com/) [Authentication package](pkg/auth) and a [GraphQL Client package](pkg/graphql) which authenticates and runs queries and mutations against [Realm's GraphQL Server](https://docs.mongodb.com/stitch/graphql/).

See [examples/](examples/)

Note: Realm is formerly known as Stitch. 

## Usage 

Note: This requires to have Atlas and Realm set up. See below for details. 

```bash
go get https://github.com/desteves/realm
import  "github.com/desteves/realm/pkg/graphql"
```

## TODO (work in progress)

- BUG: Don't show errors array when no error.
- Add Atlas API && Stitch-CLI commands for the Atlas+Realm Set up

## Atlas Setup 

- Create new project under an organization. Register [here](https://www.mongodb.com/cloud/atlas/register)
- Create new free cluster, call it `graphqlDatabase`.
- Load 'Sample Data Set'. Steps [here](https://docs.atlas.mongodb.com/sample-data/)

## Realm Set up

### Via the UI

- Create new stitch app, call it `graphqlServer`.  Steps [here](https://docs.mongodb.com/stitch/procedures/create-stitch-app/)
  - Write down app id `graphqlserver-?????`
- Enable `anonymous` authentication. Steps [here](https://docs.mongodb.com/stitch/authentication/anonymous/#configuration)
- Enable GraphQL and configure. Steps [here](https://docs.mongodb.com/stitch/graphql/expose-data/)
  - Generate 'stitch schemas'
  - Select any of the `sample-*` databases/collections and click 'Generate Schema'. 
    - Enable read/write access

### Via the CLI

- TODO

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

### (Optional)  Set up dummy GraphQL record

This is to verify our client can reach the graphql server and obtain a valid response.

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

- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
- [MongoDB GraphQL Docs](https://docs-mongodbcom-staging.corp.mongodb.com/stitch/nick/graphql/graphql.html)