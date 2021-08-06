# auth-blog-service

![Actions Workflow](https://github.com/joaomarcuslf/auth-blog-service/workflows/go/badge.svg)

## How to Start Development

1. Copy ```sample.env``` to ```.env``` and rename the variables if you need
2. Build the images and run the containers:

```sh
$ yarn install
$ cp sample.env .env
$ docker-compose up --build
```

- API: [http://localhost:5000](http://localhost:5000)
