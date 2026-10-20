const fs = require("fs")
const path = require("path")

const rawColor =
  /#(?:[0-9a-fA-F]{8}|[0-9a-fA-F]{6}|[0-9a-fA-F]{3,4})\b|\b(?:rgb|rgba|hsl|hsla|hwb|lab|lch|oklab|oklch)\s*\(|\b(?:white|black|red|green|blue|gray|grey|yellow|orange|purple|pink|brown|cyan|magenta|silver|maroon|navy|teal|olive|lime|aqua|fuchsia)\b/i

const cssColorProp =
  /(?:^|[^\w-])((?:background(?:-color)?|color|fill|stroke|border(?:-(?:top|right|bottom|left))?(?:-color)?|outline(?:-color)?|box-shadow|caret-color|accent-color|column-rule-color|text-decoration-color|-webkit-text-fill-color))\s*:\s*([^;}]+)/gi

const jsColorProp =
  /(?:style\.(?:color|background|backgroundColor|borderColor|borderTopColor|borderRightColor|borderBottomColor|borderLeftColor|outlineColor|fill|stroke|caretColor|accentColor|boxShadow|textDecorationColor|webkitTextFillColor)|(?:fillStyle|strokeStyle|shadowColor))\s*=\s*(['"`])([\s\S]*?)\1/g

const why =
  "Use a var(--token) (why: an inline colour is missed when the palette changes)."

function declarationValues(style) {
  return style
    .split(";")
    .map((part) => {
      const index = part.indexOf(":")
      return index === -1 ? "" : part.slice(index + 1)
    })
    .join(" ")
}

function hasRawColor(value) {
  const stripped = value
    .replace(/var\(\s*--[A-Za-z0-9_-]+\s*\)/g, "")
    .replace(/\b(?:transparent|currentcolor|inherit|initial|unset|none)\b/gi, "")
  return rawColor.test(stripped)
}

function lineOf(text, index) {
  return text.slice(0, index).split("\n").length
}

function findInlineColors(text) {
  const hits = []
  const seen = new Set()
  const report = (index, sample) => {
    const line = lineOf(text, index)
    const key = `${line}:${sample}`
    if (seen.has(key)) return
    seen.add(key)
    hits.push({ line, text: sample.trim(), message: `Inline colour "${sample.trim()}". ${why}` })
  }

  for (const match of text.matchAll(/style\s*=\s*"([^"]*)"/gi)) {
    if (hasRawColor(declarationValues(match[1]))) report(match.index, match[1])
  }
  for (const match of text.matchAll(/style\s*=\s*'([^']*)'/gi)) {
    if (hasRawColor(declarationValues(match[1]))) report(match.index, match[1])
  }
  for (const match of text.matchAll(cssColorProp)) {
    if (hasRawColor(match[2])) report(match.index, `${match[1]}: ${match[2]}`)
  }
  for (const match of text.matchAll(jsColorProp)) {
    if (hasRawColor(match[2])) report(match.index, match[2])
  }
  const svgAttr = /\b(?:fill|stroke|stop-color|flood-color|lighting-color)\s*=\s*(['"])([^'"]*)\1/gi
  for (const match of text.matchAll(svgAttr)) {
    if (hasRawColor(match[2])) report(match.index, match[0])
  }
  const setAttr =
    /\.setAttribute\(\s*(['"])(?:fill|stroke|stop-color|flood-color|lighting-color)\1\s*,\s*(['"`])([\s\S]*?)\2/gi
  for (const match of text.matchAll(setAttr)) {
    if (hasRawColor(match[3])) report(match.index, match[3])
  }
  return hits
}

function scanFile(file) {
  const text = fs.readFileSync(file, "utf8")
  return findInlineColors(text).map((hit) => ({ ...hit, file }))
}

function scanDir(root) {
  const hits = []
  const walk = (dir) => {
    for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
      const full = path.join(dir, entry.name)
      if (entry.isDirectory()) {
        walk(full)
        continue
      }
      if (!entry.name.endsWith(".html") && !entry.name.endsWith(".js")) continue
      hits.push(...scanFile(full))
    }
  }
  walk(root)
  return hits
}

module.exports = { findInlineColors, scanDir, hasRawColor, why }
