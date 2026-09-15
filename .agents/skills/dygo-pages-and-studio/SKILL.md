---
name: dygo-pages-and-studio
description: Define dygo App Pages and metadata-driven Studio experiences for Business Apps. Use when an App needs a Page, navigation entry, renderer choice, or Studio presentation beyond default Entity routes.
---

# dygo Pages And Studio

Use metadata to select `entity-index` or a sibling Vue file. Do not place UI implementation inside Page metadata.

## Start

Read `docs/pages.md`, `docs/studio.md`, and the Page-related parts of `docs/app-model.md` and `docs/dir.md`. Confirm which Page features exist now.

## Rules

- Keep the Page bundle app-owned under the documented `pages/` path.
- Use `dygo generate page <app>/<page>` to create YAML, Vue, and access files.
- Use `renderer: entity-index` for the framework home index, or `renderer: vue` with a sibling `<page>.vue` file.
- Let Studio own rendering, layout, permissions, loading, errors, and shared interaction patterns.
- Prefer normal Entity and Record surfaces when they solve the task.
- Add a Page only when it represents a useful business Space or cross-Entity entry point.
- Keep navigation labels and actions in dygo vocabulary.
- Do not encode arbitrary frontend code, SQL, or unvalidated behavior in metadata.
- Put Vue UI in `<page>.vue` beside the YAML. Do not implement a separate App frontend bundle contract.

## Check

Validate Page metadata through the current project validation path. Confirm that boot data, navigation, route resolution, permissions, and renderer registration agree.
