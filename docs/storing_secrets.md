# <img height="25" src="https://user-images.githubusercontent.com/10101283/66178622-8f14d480-e62b-11e9-8db7-d18cc7885fb3.png"> &ensp; Storing Secrets

:red_circle: __Never store secrets in your committed code base__ :red_circle: 

Instead, store secrets in the environment variables.

## Development Secrets

`a1` make this easy because it uses [dotenv](https://github.com/joho/godotenv). By default, `a1` will load a file called `.env` from the directory you run the executable from.

### Specify the `dotenv` path

When you execute the binary, you can set the `ENV_FILE` variable to specify the path to the `dotenv` file. If you do not specify the `ENV_FILE`, it defaults to just `.env`. __The path is always relative to the current directory__. So if your folder structure looks something like this:

- `project/`
  - `out/`
    - __`server`__
  - __`server.go`__
  - __`.env`__

Below are some examples of how to specify the `ENV_FILE` when you execute the binary from different directories.

```bash
# /path/to/project
ENV_FILE=.env ./out/server

# /path/to/project - Not specified, defaults to ENV_FILE=.env
./out/server

# /path/to/project/out
ENV_FILE=../.env ./out/server
```

> If the `ENV_FILE` file does not exist or you specified the wrong path, `a1` will tell you that the file could not be found during startup. It __DOES NOT__ stop the server from starting.

### Example `.env`

`a1` uses a couple of environment variables internally, so it's a good idea to set them explicitly. That way you won't override them with your own environment variables by default

```dotenv
DEV=true|false
VERBOSE=true|false
STARTUP_ONLY=true|false
```

To learn more about how to setup your `.env` files, checkout the [dotenv](https://github.com/joho/godotenv) documentation.

## Production Secrets

__You should not use a `.env` on your production instance__. Instead, set the environment variables directly using whatever tools your host gives you.

Checkout the [deployment docs](https://github.com/aklinker1/a1/tree/master/docs/deployment.md) to learn more about deploying your server to production.
