const fs = require("fs")
const path = require("path")
const stylelint = require("stylelint")

const ruleName = "lint/no-hex-color"
const messages = stylelint.utils.ruleMessages(ruleName, {
  rejected: (hex) =>
    `Unexpected hex colour ${hex}. Use a token from the theme file (why: a raw hex is missed when the palette changes).`,
})
const hexColor = /#(?:[0-9a-fA-F]{8}|[0-9a-fA-F]{6}|[0-9a-fA-F]{4}|[0-9a-fA-F]{3})\b/g

const rule = (primary, secondaryOptions) => {
  return (root, result) => {
    const validOptions = stylelint.utils.validateOptions(
      result,
      ruleName,
      { actual: primary },
      {
        actual: secondaryOptions,
        possible: {
          themeFile: (value) => typeof value === "string" && value.length > 0,
          allowlistFile: (value) => typeof value === "string",
        },
        optional: true,
      },
    )
    if (!validOptions || primary === false) return

    const options = secondaryOptions || {}
    const themeFile = options.themeFile || "theme.css"
    const file = (root.source && root.source.input && root.source.input.file) || ""
    const normalized = file.replaceAll("\\", "/")
    if (normalized.endsWith(themeFile.replaceAll("\\", "/"))) return

    const allowed = options.allowlistFile ? allowlistLines(options.allowlistFile) : new Set()
    root.walkDecls((decl) => {
      for (const match of String(decl.value).matchAll(hexColor)) {
        const line = decl.source && decl.source.start ? decl.source.start.line : 0
        const key = `${normalized}:${line}`
        if (allowed.has(key)) continue
        stylelint.utils.report({
          message: messages.rejected(match[0]),
          node: decl,
          index: match.index,
          endIndex: match.index + match[0].length,
          result,
          ruleName,
          word: match[0],
        })
      }
    })
  }
}

function allowlistLines(file) {
  const abs = path.isAbsolute(file) ? file : path.join(process.cwd(), file)
  if (!fs.existsSync(abs)) return new Set()
  const lines = new Set()
  for (const raw of fs.readFileSync(abs, "utf8").split("\n")) {
    const trimmed = raw.trim()
    if (!trimmed || trimmed.startsWith("#")) continue
    lines.add(trimmed.split(/\s+/)[0])
  }
  return lines
}

rule.ruleName = ruleName
rule.messages = messages
module.exports = stylelint.createPlugin(ruleName, rule)
