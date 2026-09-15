# Page Metadata

Pages are App-owned bundles: YAML for dygo, and an optional Vue file for the screen. YAML contains metadata, not UI code.

## Bundle shape

```txt
apps/<app>/pages/<page>/<page>.page.yml
apps/<app>/pages/<page>/<page>.vue
```

Create that bundle with:

```sh
dygo generate page <app>/<page>
```

The generator writes Page metadata, a Vue starter, and `access/<page>.page.access.yml`. It does not overwrite an existing Vue file. `--force` refreshes dygo-generated YAML and access only.

## Framework renderer

Home uses the built-in entity index:

```yaml
label: Home
description: Start page for the entities available to the current Studio user.
icon: house
route:
  path: /
renderer: entity-index
options: {}
```

`path` is `/` or one kebab-case root path such as `/board`. `renderer` is `entity-index` or `vue`. `options` is a YAML map passed to the Page component.

## Vue renderer

`renderer: vue` loads the sibling `<page>.vue` file. Studio compiles that file into the Studio UI. The component receives a `page` prop with label, description, path, options, and owning App.

Page identity is app-scoped. The bundle directory is the Page key. The runtime name is `<app>.<key>`.

Page metadata is loaded with Entity metadata. Route ownership is checked with framework-reserved and Entity routes, so two Pages or a Page and Entity cannot claim the same public path.

The public Go contract is in `pkg/dygo`. Studio owns the Page host, routing, loading, and error states. Global CRUD remains framework-owned.

A Vue Page can use Studio design components and Record clients the same way other Studio screens do. Permissions stay on the server. Opening a Page does not grant access to its Records.

Studio does not put Pages in the Entity sidebar. Pin a Page from its header to keep it in the personal Pinned section.
