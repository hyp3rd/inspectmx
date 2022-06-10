# INSPECTMX

A simple http service to validate an email address against a subset of pre-configured [providers](./src/.config.yml).

## HOW TO USE IT

This is a microservice written in `golang`.
On 'nix based systems, build the `Docker` image as follows:

```bash
./tools/build.bash
```

And run it locally leveraging the docker-compose stack:

```bash
docker-compose up
```

The requests are described in an http file [here](./requests/all.http).
***I recommend testing http services with [rest-clint](https://github.com/Huachao/vscode-restclient) if you're using VSCode.***

## NEXT VERSION

The service is ready for **Redis** to allow more intense workloads and broadened distribution.
In the future, I will add the features to verify the address's existence.

## DISCLAIMER

The software is provided `as is.` Use it at your own risk.
Pull requests are more than welcome if you find bugs or security flaws.
