const fs = require("fs")
const path = require("path")
const assert = require("node:assert/strict")
const { test } = require("node:test")
const stylelint = require("stylelint")
const createConfig = require("./create-config.cjs")
const { findInlineColors } = require("./plugins/inline-color.js")
const { svgPaletteHits, themeSolids } = require("./plugins/svg-palette.js")

const fixtures = path.join(__dirname, "fixtures")
const why = /why:|theme file|var\(--token\)|recolours|theming/i

async function lintFile(file) {
  return stylelint.lint({
    files: [path.join(fixtures, file)],
    config: createConfig(),
    configBasedir: __dirname,
  })
}

async function lintCode(code, codeFilename, config) {
  return stylelint.lint({
    code,
    codeFilename,
    config: config || createConfig(),
    configBasedir: __dirname,
  })
}

function warnings(result) {
  return (result.results[0] && result.results[0].warnings) || []
}

function rulesOf(result) {
  return new Set(warnings(result).map((w) => w.rule))
}

async function expectRule(file, rule) {
  const result = await lintFile(file)
  const found = warnings(result)
  assert.ok(
    found.some((w) => w.rule === rule),
    `${file} missing ${rule}: ${JSON.stringify(found)}`,
  )
  for (const w of found.filter((item) => item.rule === rule)) {
    assert.match(w.text, why, `${rule} message should explain why: ${w.text}`)
  }
  console.log(`css-rule fail ${rule}`)
}

async function expectClean(file) {
  const result = await lintFile(file)
  assert.equal(warnings(result).length, 0, `${file}: ${JSON.stringify(warnings(result))}`)
  console.log(`css-rule pass ${file}`)
}

test("failing fixture: raw hex, named colour and colour functions", async () => {
  const result = await lintFile("bad.css")
  const rules = rulesOf(result)
  assert.ok(rules.has("color-no-hex"))
  assert.ok(rules.has("color-named"))
  assert.ok(rules.has("function-disallowed-list"))
  for (const w of warnings(result)) {
    assert.match(w.text, why, `message should explain why: ${w.text}`)
  }
  console.log("css-rule fail color-no-hex")
  console.log("css-rule fail color-named")
  console.log("css-rule fail function-disallowed-list")
})

test("passing fixture: var()-only CSS", async () => {
  await expectClean("good.css")
})

test("theme file path is configurable and exempt", async () => {
  const code = ":root { --ink: #f4f6f8; --danger: red; --glow: rgb(1, 2, 3); }"
  const blocked = await lintCode(code, "src/app.css")
  assert.ok(warnings(blocked).length > 0, "raw colours outside theme must fail")

  const allowed = await lintCode(
    code,
    "web/styles/theme.css",
    createConfig({
      themeFiles: ["web/styles/theme.css"],
      themeFile: "web/styles/theme.css",
      deprecated: [],
      colorPrefixes: [],
    }),
  )
  assert.equal(warnings(allowed).length, 0, JSON.stringify(warnings(allowed)))
})

const pairs = [
  ["old-scheme-bad.css", "old-scheme-good.css", "lint/no-old-scheme"],
  ["hex-bad.css", "hex-good.css", "lint/no-hex-color"],
  ["builtin-color-bad.css", "builtin-color-good.css", "color-no-hex"],
  ["colors-in-theme-bad.css", "colors-in-theme-good.css", "lint/colors-in-theme"],
  ["placeholder-bad.css", "placeholder-good.css", "lint/no-placeholder-color"],
  ["clipped-heading-bad.css", "clipped-heading-good.css", "lint/no-clipped-heading"],
  ["layout-animation-bad.css", "layout-animation-good.css", "lint/no-layout-animation"],
  ["unknown-token-bad.css", "unknown-token-good.css", "lint/unknown-token"],
  ["page-component-bad.css", "page-component-good.css", "lint/no-page-component-css"],
  ["hygiene-bad.css", "hygiene-good.css", "declaration-no-important"],
]

for (const [bad, good, rule] of pairs) {
  test(`${rule} fails ${bad} and passes ${good}`, async () => {
    await expectRule(bad, rule)
    await expectClean(good)
  })
}

test("hygiene, font-size and z-index rules fail the bad fixture", async () => {
  const result = await lintFile("hygiene-bad.css")
  const rules = rulesOf(result)
  for (const rule of [
    "declaration-no-important",
    "declaration-block-no-duplicate-properties",
    "no-duplicate-selectors",
    "declaration-block-no-shorthand-property-overrides",
    "property-no-unknown",
    "declaration-property-value-allowed-list",
  ]) {
    assert.ok(rules.has(rule), `missing ${rule}: ${[...rules].join(",")}`)
    for (const w of warnings(result).filter((item) => item.rule === rule)) {
      assert.match(w.text, why, w.text)
    }
    console.log(`css-rule fail ${rule}`)
  }
})

test("builtin colour rules each explain why", async () => {
  const result = await lintFile("builtin-color-bad.css")
  for (const rule of [
    "color-no-hex",
    "color-named",
    "function-disallowed-list",
    "scale-unlimited/declaration-strict-value",
  ]) {
    assert.ok(rulesOf(result).has(rule), rule)
    console.log(`css-rule fail ${rule}`)
  }
})

test("inline-color scanner", () => {
  const bad = findInlineColors(fs.readFileSync(path.join(fixtures, "inline-color-bad.html"), "utf8"))
  const good = findInlineColors(fs.readFileSync(path.join(fixtures, "inline-color-good.html"), "utf8"))
  assert.ok(bad.length >= 4, JSON.stringify(bad))
  for (const hit of bad) assert.match(hit.message, why)
  assert.equal(good.length, 0, JSON.stringify(good))
  console.log("css-rule fail inline-color")
  console.log("css-rule pass inline-color-good.html")
})

test("svg-palette scanner", () => {
  const solids = themeSolids(fs.readFileSync(path.join(fixtures, "theme.css"), "utf8"))
  const bad = svgPaletteHits(fs.readFileSync(path.join(fixtures, "svg-palette-bad.svg"), "utf8"), solids)
  const good = svgPaletteHits(fs.readFileSync(path.join(fixtures, "svg-palette-good.svg"), "utf8"), solids)
  assert.equal(bad.length, 1)
  assert.match(bad[0].message, why)
  assert.equal(good.length, 0)
  console.log("css-rule fail svg-palette")
  console.log("css-rule pass svg-palette-good.svg")
})
