# Studio App

Studio is the first-party UI App for dygo.

It will provide the global interface for Spaces, Records, Forms, Lists, Saved Views, jobs, audit logs, settings, and metadata-driven rendering.

The Studio app now includes design context and the first UI scaffold:

- `PRODUCT.md` defines Studio's product register, users, purpose, anti-references, and design principles.
- `DESIGN.md` defines Studio's visual direction, Shell model, Page vocabulary, Page Types, and Dygo UI component architecture.
- `ui/` contains the Vue/Vite Studio frontend, router, route guards, and the first Dygo UI components.

Feature code should use Dygo UI components from `ui/src/design/`. Only the design system imports Reka UI; feature, shell, renderer, and page code imports the styled primitives from `@dygo/ui` instead of Reka directly.

## Development

Studio should be opened through the dygo server port:

```sh
go run ./cmd/dygo dev
```

When the source checkout contains `apps/studio/ui/package.json`, `dygo dev` starts Studio's development asset server internally and proxies it through dygo's configured server address, normally `http://127.0.0.1:6790/`.

Use `dygo dev --studio-dev-url` only when you want to run the Studio asset server yourself. `Ctrl-C` stops the dygo server and the auto-started Studio dev server together.

## Routes

Studio is root-mounted by default. Global pages and record pages both live at root paths:

```txt
/login
/
/:entity
/:entity/new
/:entity/:recordName
```

Dynamic Entity route slugs are authenticated. `/login` is public and redirects authenticated users back to `/`, which resolves through the boot default `home` route when configured. Activity appears inside the Record page.
