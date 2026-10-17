# CSS colour-token rules

Ban hard-coded hex, named colours and colour functions (`rgb`, `hsl`, and the
rest) everywhere except the theme file. Colours live in one place so a
recolour or theme switch cannot miss values.

## Wire into a repo

1. Copy this directory's config into the repo (or merge the rules). Point the
   theme file path at your palette file:

   ```js
   // stylelint.config.cjs
   module.exports = require("./lint/css/create-config.cjs")({
     themeFiles: ["web/styles/theme.css"],
   })
   ```

   The default `themeFiles` is `["**/theme.css"]`.

2. Install Stylelint:

   ```sh
   npm install -D stylelint@16
   ```

3. Run in CI so a failure blocks merge:

   ```sh
   npx stylelint "src/**/*.css" "web/**/*.css"
   ```

   On an existing codebase, allowlist current hits first, then burn them down.

## What it catches

| rule | bad |
| --- | --- |
| `color-no-hex` | `color: #f4f6f8` |
| `color-named` | `background-color: red` |
| `function-disallowed-list` | `border-color: rgb(10, 20, 30)` |

The theme file named in `themeFiles` may still use raw colours to define
tokens. Everywhere else, use `var(--token)`.

## Test

```sh
npm --prefix skills/lint-rules/rules/css test
```

Or run [../test.sh](../test.sh), which includes this suite.
