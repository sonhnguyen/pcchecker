
# go-getting-started

A barebones Go app, which can easily be deployed to Heroku.

This application supports the [Getting Started with Go on Heroku](https://devcenter.heroku.com/articles/getting-started-with-go) article - check it out.
## API DOCS FOR OUR FRONTEND MASTERPIECE

- POST /createBuild
  - params: Array of product ID
  
- GET /getBuildById?id=123123
  - URL params: id
  
- GET /build/:encodedURL
  - URL params: encodedURL
  
- GET /getBuildRecent?limit=10
  - URL params: limit, if not input default is 10

- GET /getProducts/:category
  - returns all product of that category

- GET /getProducts?query=asdasd
  - URL params: query
  - Performs search product by text, return a list with relevant sorted

- GET /product/:id/
  - Param id: return product of that id

## Running Locally

Make sure you have [Go](http://golang.org/doc/install) and the [Heroku Toolbelt](https://toolbelt.heroku.com/) installed.

```sh
$ go get -u github.com/heroku/go-getting-started
$ cd $GOPATH/src/github.com/heroku/go-getting-started
$ heroku local
```

Your app should now be running on [localhost:5000](http://localhost:5000/).

You should also install [Godep](https://github.com/tools/godep) if you are going to add any dependencies to the sample app.

## Deploying to Heroku

```sh
$ heroku create
$ git push heroku master
$ heroku open
```

or

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)


## Documentation

For more information about using Go on Heroku, see these Dev Center articles:

- [Go on Heroku](https://devcenter.heroku.com/categories/go)
