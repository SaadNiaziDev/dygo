# Directory Shape

## Project

```txt
project/                         - Generated dygo project root
  .gitignore                     - Project ignore rules
  dygo.yml                       - Project config and marker
  cmd/                           - Project binaries live here
    dygo/                        - Project dygo runner
      main.go                    - Hook-aware CLI entrypoint
  apps/                          - App packages live here
    <app>/                       - App package bundle
      app.yml                    - App manifest and paths
      entities/                  - App Entity definitions
        <entity>/                - Normal Entity bundle
          <entity>.entity.yml    - Entity metadata definition
          hooks.go               - Entity hook scaffold
          fixtures.yml           - Entity fixture records
          views.yml              - Reserved Entity view metadata (not loaded yet)
        _collections/            - Collection row definitions
          <collection>.yml       - Single-file collection metadata
          <collection>/          - Folder-form collection bundle
            <collection>.entity.yml - Collection metadata definition
      access/                    - App access metadata
        _roles.yml               - App-contributed global roles
        <entity>.access.yml      - Entity access policy contribution
        <page>.page.access.yml   - Page access policy contribution
      jobs/                      - App background jobs
        <job>/                   - Job bundle
          job.yml                - Job metadata definition
          run.go                 - Job runner code
        _schedules.yml           - Recurring job schedules
      pages/                     - Custom app pages
        <page>/                  - Custom page bundle
          <page>.page.yml        - Page metadata definition
          <page>.vue             - Vue Page view for renderer: vue
      reports/                   - Cross-Entity report definitions (not loaded yet)
        <report>.yml             - Single-file report metadata
        <report>/                - Folder-form report bundle
          report.yml             - Report metadata definition
  db/                            - Database generated artifacts
    schema.sql                   - PostgreSQL schema snapshot
  docs/                          - Project documentation files
  config/                        - Project config files
    secrets/                     - Encrypted environment secrets
    storage.yml                  - Future storage config
    queues.yml                   - Queue registry and concurrency config
    logging.yml                  - Future logging config
  .dygo/                         - Framework-managed Apps and local runtime state
    apps/
      core/                      - Tracked Core metadata required by fresh checkouts
      studio/                    - Studio metadata and cached UI assets
    files/                       - Local uploaded files
    logs/                        - Local runtime logs
    tmp/                         - Local temporary files
    secrets/                     - Local private secret keys
```

## Runtime

A deployed runtime layout is not yet implemented. `dygo deploy` is proposed in the issue tracker and will define the deployed layout when it ships.

## Framework

The framework repository uses this working tree:

```txt
dygo/                           - Framework repository root
  cmd/                          - Framework binaries live here
    dygo/                       - Stock dygo CLI
  internal/                     - Private framework packages
    access/                     - App access metadata load, apply, and export
    accesspolicy/               - Conditional access policy AST
    actions/                    - Entity action registry and execution
    app/                        - App manifest loading and app registry
    auth/                       - Password hashing and session authentication
    cli/                        - Cobra command implementation
    config/                     - Config loading defaults
    corevalues/                 - Built-in metadata constants
    db/                         - PostgreSQL schema sync, records, and metadata
    dygodata/                   - Internal adapters for the public SDK
    entity/                     - Entity metadata catalog and field types
    files/                      - Private file service
    fixtures/                   - Fixture loading runtime
    frameworkapp/               - Framework-managed Core app installation
    fsutil/                     - Shared filesystem helpers
    generate/                   - App, Entity, hook, job, and page scaffolding
    health/                     - Health check handlers
    hookevents/                 - Hook event definitions
    hookgen/                    - Hook scaffold generator
    hooks/                      - Hook runtime registry
    imports/                    - Durable CSV import service
    jobgen/                     - Job scaffold generator
    jobs/                       - Job metadata, runtime, and queue store
    migration/                  - App lifecycle migration transaction
    naming/                     - Record naming strategies
    notifications/              - Notification inbox and email delivery
    pages/                      - App page metadata loading
    patches/                    - Explicit patch runtime
    permissions/                - Permission evaluation logic
    project/                    - Project root discovery and metadata loading
    projectgen/                 - Project scaffold generator
    queues/                     - Queue config loading
    recordfilter/               - Record filter parsing
    recordquery/                - Record query helpers
    recordsecret/               - Record field encryption and key ring
    reserved/                   - Reserved name registry
    routes/                     - Route registry and validation
    runnergen/                  - Project runner wiring generator
    schedules/                  - Schedule metadata loading
    secrets/                    - Encrypted secrets runtime
    server/                     - HTTP server runtime
    shape/                      - Filesystem path and filename conventions
    studio/                     - Studio asset handling
    studiostate/                - Private Studio state storage
    upgrade/                    - Project upgrade runtime
    yamlmeta/                   - YAML metadata helpers
  apps/                         - First-party dygo apps
    core/                       - Core platform app
    studio/                     - Studio web app
  pkg/                          - Public Go API surface
    dygo/                       - App hook, Job, and logging API
  config/                       - Framework runtime config files
    secrets/                    - Encrypted dev secrets
    github.yml                  - GitHub repository and issue-tracking metadata
  db/                           - Framework DB artifacts
    schema.sql                  - Framework schema snapshot
  docs/                         - Framework documentation
  schemas/                      - Editor JSON Schemas
  scripts/                      - Release and helper scripts
```
