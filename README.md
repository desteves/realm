# mongodb realm

Contains [mongoDB realm](https://stitch.mongodb.com/) (formerly known as stitch) [Authentication](pkg/auth) and a [GraphQL Client](pkg/graphql) go packages.


See [examples/](examples/)


## Usage 

`go get https://github.com/desteves/realm`
`import  "github.com/desteves/realm/pkg/graphql"`

## TODO (work in progress)

- Add stitch config steps
  - creating webhook for ping test
  - creating dummy graphql record for testing
- Test Refresh Token
- Test all auth mechanisms, so far only test anonymous login.
- Add `*_test` files/packages
- Implement `Disconnect` for auth and graphql packages.

## References

- [Mongo Go Driver](https://github.com/mongodb/mongo-go-driver)
- 