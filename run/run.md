# Run

Folder that contains the build docker-compose related stuff. Meant to be self contained in one host - three containers, but can be run in different places with no problems.

The password for encryption of the slog-server storage can be changed or updated in the `docker-compose.yml`, if running the binary in another place is desired, just point the `slog-server` to this docker-compose host by using it's config file.

The folder `data` is for long term storage, so the data stored and changes done to the Dashboards are persistent.

To start the stack, after running the `build-docker-image.sh` just run `docker-compose up`, this will start the `slog-server`, Grafana and Influx.

Grafana can be reached at <http://localhost:3000>, login credentials are `admin/admin`.

## Additional info

The official Grafana documentation and how-to's can be found here: <https://grafana.com/docs/grafana/latest/getting-started/getting-started/>

This uses Influx DB as a backend, as a time series database makes more sense for a logging program than a relational database.

Though it shouldn't be required if only Grafana is used, the official documentation for Influx can be found here: <https://www.docs.influxdata.com/influxdb/v1.8> if the need of fancy queries rises.
