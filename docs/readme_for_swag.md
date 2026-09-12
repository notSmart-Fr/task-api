### Swagger Documentation
```bash
swag init --v3.1 -g cmd/server/main.go
```

### Known bug: stale docs in Docker
`docker-compose.yml` used to mount a named volume at `/app` (the WORKDIR), which
shadowed the `docs/` folder and `server` binary baked into the image. Docker only
seeds a named volume from image content on first creation, so regenerating docs
with `swag init` and rebuilding never showed up in the running container (Scalar
appeared to have no routes). Fixed by mounting the volume at `/app/data` instead
and pointing the SQLite `DB_PATH` there. After changing volume mounts, always run
`docker-compose down -v` before `up --build` to drop stale volumes.
