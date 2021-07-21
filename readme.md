# ACC Race Engineer hub

This is the hub application that receives data and forwards it to the appropriate connected websocket clients.

## Compiling and running
Your `go version` must support modules in order for `go build` to obtain the necessary dependencies. Currently `mongodb` is the only supported database.

1. Install mongodb and create a new database.
2. `cd cmd/are_hub/`.
3. Create a copy of `config.json.example`.
4. Rename the copy to `config.json`.
5. Edit the file to reflect your environment.
6. `go build`.
7. `./are_hub` or `./are_hub --config <path_to_json_configuration_file>` (if the config file is not in the same directory as the executable and/or not named `config.json`).

## Deployment
1. Create `config.json.docker` in the application's root directory. This is the file that will be copied and used as the deployments `config.json` and the keys are identical as described in `cmd/are_hub/config.json.example`
2. `docker build .`

## Licence
BSD-3-clause
