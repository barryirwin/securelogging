
# Secure Logging

## Background
This code repository was developed as part of a final year research project by Franco Loyola.  The project was supervised by Prof Barry Irwin.

## Publications
This work has been published  as follows:
ICONIC 2022 -TBA

### Citation
Cite this work as:
TBA

## The project name: slog
Why? The initial name was `SecureLogging` but since it was too long to write and not pretty to use, decided to use a shortened version of it `sLog`, but having mix of small and capital letters was also cumbersome to write "properly" everywhere, settled with `slog`, which picks a bit on the user, a-la `Cockroach-DB`, or an Oddworld reference, who knows.

## Repo folders structure
First, how to find things in this repo: Each folder focus on a given aspect of the project, and they contain more information in the form of extra `README.md`  files within them.

 ### **code**
Contains all the Go source code for the project.
### build
Contains the build script for all the binaries. They can be built from within each `cmd` folder as well. It also contains some example configuration files and keys used during development.

### run
Contains the docker-compose file, Dockerfile, Docker build script, etc. It also doubles as the persistent storage for Influx and Grafana, it depends on the contents of the `build` folder.

  



## How the system works


  

It is composed by a client that will tail constantly log files and send data to a server, and the server that receives the data and stores it in `.txt` and `.slog` format.

Also, if the Influx output is enabled, the received data will be sent to the DB.

  

The main code for each binary can be found under the corresponding folder, under `code/cmd/`, most of the code is at `code/pkg`.

  


## Build process

It is very simple to get the platform up and running. Install Go, and then run the build script `build-bins.sh`, which is located under the `build` folder.


- Go installation: <https://golang.org/doc/install>

  

`go mod` should fetch all dependencies automatically, but if the Go version is too old, you might need to use `go get $some-package`

  

Then use the `build/build-bins.sh` to get stand-alone binaries or `run/build-docker-image.sh` to build the Docker image of the server, which can be used by `docker-compose`.

  



## Using the binaries (Stand-alone, no DB on the side)



  

Both binaries have built-in help, just run them with the `-h` flag, will provide all available flags, their defaults and a description.

  

**REPLACE THE KEYS** located at `build/keys` when using this programs, the ones in the repo are  a placeholder and for testing/development purposes only!

New keys can be generated from the `slog-server` binary using the following flags:

  

```bash

slog-server.bin -generate-comms-keys

slog-server.bin -generate-storage-keys

```

  

Copy the binaries to the desired hosts, and start them with the `[client|server].conf` file on the same folder.

  

**FOR THE CLIENT** only use the `comms/` (communications) public key! This can be configured in the config file.

  

**FOR THE SERVER** If the Influx output is enabled, make sure that Influx/Grafana are running somewhere and point it there, run it somewhere else or use the `docker-compose.yml` provided at the `run/` folder with the part of the `slog-server` commented out.

  

After some data has been stored, it can be exported. Instructions further below in this guide.



## Using the binaries (Grafana + Influx to see data arriving)



  

This requires that `docker` and `docker-compose` installed in the host that will run the binaries. Also, to build the Docker images, will require internet connection to fetch the base images.

  

- Docker installation: <https://docs.docker.com/get-docker/>

- Docker-compose installation: <https://docs.docker.com/compose/install/>

  

To build the Docker images, run the `build-docker-images.sh` located at the `run` folder, it will build the binaries, put them into the images for later use.

  

If the images are already built, then just run `docker-compose up` at the `run` folder (check that the user ID matches the one running the one in the `docker-compose.yml` file), and go to <http://localhost:3000>. There is already some folders mapped to Grafana and Influx to store the data, this is mounted at the container startup.

  

Grafana is still using the default login `admin`/`admin`.

  


## Decrypt data stored by the server



Data can only be retrieved as txt, comma separated after is received, and data can only be written once to the text-binary format, valid edition or deletion of it is not possible from this program, there is no code to decode and handle that, though it should be possible by compiling another version and using the encoders/decoders.

  

To replay the data, the `slog-server.bin` needs the private key and the common password used at the capture time, else the data will not be retrievable anymore. This is to assure that the only people able to retrieve the data is the people that setup the environment.

  

Also, the read-only nature of the data retrieval makes possible that any party interested to check that the logs that arrived to the server are the ones that were presented to them, they can simply ask for the key and password and compare that both txt files match, else someone manipulated them.

  

```bash

./slog-server.bin -password $password-used-at-store-time -priv-key path/to/private/storage/key -influx-ip http://localhost:8086 -read-file path/to/slog-logging.slog

```

  

### Compare replay vs. "old" realtime file

  

If the replay was done with the Influx output, the same data should be written to the DB, as the reference timestamp used is the one stored in the binary encrypted struct. The difference is in the `replayData` tag, which will be set to true. There is an example for this in Grafana already, check the Dashboard: TO-DO

  

Another way to compare both outputs, is to do a `diff` between the `slog-logging.txt`(default txt output on real time) and the `slog-reprocessed.txt` (default reprocessed)

  

This command should output nothing, if there is a difference in the files, will be shown, that means that the `.slog` file has been tampered with.

  

```bash

diff slog-logging.txt slog-reprocessed.txt

```

  

## Program flow


See each `client` or `server` folder, as their logic differs.