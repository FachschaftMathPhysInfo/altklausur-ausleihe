# Hacking on this project

## Setup a dev environment

1. Clone this repository:
   ```
   $ git clone https://github.com/FachschaftMathPhysInfo/altklausur-ausleihe.git
   ```

2. Install [`docker`](https://docs.docker.com/engine/install/) and the [`docker compose`](https://docs.docker.com/compose/install/) plugin.
3. Change into the `env_files` directory and copy the example env files:
   ```
   $ cd env_files
   $ cp backend.env.example backend.env
   $ cp exam_marker.env.example exam_marker.env
   $ cp minio.env.example minio.env
   $ cp postgres-credentials.env.example postgres-credentials.env
   ```
4. Build and start the project:
   ```
   $ docker compose up --build
   ```
4. Visit http://localhost:8080 to see the project
5. For the test instance you can go to http://localhost:8080/testlogin to obtain a valid JWT token

## Uploading exams

See the documentation in [`./tools/README.md`](./tools/README.md).
