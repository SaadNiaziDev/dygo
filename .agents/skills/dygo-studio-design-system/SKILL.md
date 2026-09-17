---
name: dygo-studio-design-system
description: Create, extend, or review dygo Studio design tokens, Tailwind v4 styles, primitives, components, and interaction patterns. Use when visual consistency or a reusable Studio UI pattern is the main task.
---

# dygo Studio Design System

Extend one coherent Studio system. Do not build a generic dashboard theme.

## Start

Inspect `apps/studio/ui/src/design`, the shell, existing CSS tokens, and representative Record and Page renderers. Read `docs/studio.md`, plus `docs/dialogs.md` or `docs/toasts.md` when relevant.

## Rules

- Use Tailwind v4 and the repository's CSS-first token approach.
- Prefer semantic tokens over raw repeated color, spacing, shadow, and state values.
- Reuse or extend atoms, primitives, molecules, and organisms at the canonical layer.
- Preserve reusable UI components and exports even when they have no current consumers. Remove them only when the user authorizes their removal.
- Keep component variants small and named by purpose or state.
- Design dense business workflows for clarity, scan speed, and keyboard use.
- Include hover, focus, active, disabled, loading, error, empty, and destructive states when applicable.
- Only `apps/studio/ui/src/design/` imports Reka UI. Feature, shell, renderer, and page code imports the styled primitives (`Popover`, `Dialog`, `Combobox`, `DropdownMenu` and its parts, `Tree`, `Select`, `Switch`, `RadioGroup`) from `@dygo/ui`. Extend the primitive layer instead of importing Reka elsewhere.
- Use Lucide for icons.
- Keep labels and actions in sentence case and dygo vocabulary.
- Make responsive behavior intentional. Do not hide essential business actions without an alternate path.
- Avoid decorative novelty that makes Studio surfaces inconsistent.

Panel classes passed to a primitive (`panel-class`, `overlay-class`, `input-class`) render outside the consumer's scoped CSS. Style them with `:global(...)` or a non-scoped style block. Primitives own structure and interaction states in `styles/base.css`; consumer panel and item classes own surface, radius, shadow, and layout. Keep primitive defaults free of properties a consumer is expected to set, because `base.css` is emitted after component styles.

Verify reusable components in representative consumers. Use browser inspection when visual acceptance matters.
