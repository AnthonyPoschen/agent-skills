const path = require("path")
const assert = require("node:assert/strict")
const { test } = require("node:test")
const stylelint = require("stylelint")
const createConfig = require("./create-config.cjs")

const fixtures = path.join(__dirname, "fixtures")

async function lintFile(file, config) {
  return stylelint.lint({
    files: [path.join(fixtures, file)],
    config: config || createConfig(),
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

test("failing fixture: raw hex, named colour and colour functions", async () => {
  const result = await lintFile("bad.css")
  const rules = new Set(warnings(result).map((w) => w.rule))
  assert.ok(rules.has("color-no-hex"), "expected color-no-hex")
  assert.ok(rules.has("color-named"), "expected color-named")
  assert.ok(
    rules.has("function-disallowed-list"),
    "expected function-disallowed-list",
  )
  for (const w of warnings(result)) {
    assert.match(
      w.text,
      /theme file|var\(--token\)|recolours|theming/i,
      `message should explain why: ${w.text}`,
    )
  }
})

test("passing fixture: var()-only CSS", async () => {
  const result = await lintFile("good.css")
  assert.equal(warnings(result).length, 0, JSON.stringify(warnings(result)))
})

test("theme file path is configurable and exempt", async () => {
  const code = ":root { --ink: #f4f6f8; --danger: red; --glow: rgb(1, 2, 3); }"
  const blocked = await lintCode(code, "src/app.css")
  assert.ok(warnings(blocked).length > 0, "raw colours outside theme must fail")

  const allowed = await lintCode(
    code,
    "web/styles/theme.css",
    createConfig({ themeFiles: ["web/styles/theme.css"] }),
  )
  assert.equal(
    warnings(allowed).length,
    0,
    JSON.stringify(warnings(allowed)),
  )
})
