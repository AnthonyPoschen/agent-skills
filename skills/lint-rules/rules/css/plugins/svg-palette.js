const fs = require("fs")
const path = require("path")

const hexColor = /#(?:[0-9a-fA-F]{8}|[0-9a-fA-F]{6}|[0-9a-fA-F]{4}|[0-9a-fA-F]{3})\b/g
const attr = /\b(?:fill|stroke|stop-color|flood-color|lighting-color)\s*=\s*(['"])([^'"]*)\1/gi
const rawColor =
  /#(?:[0-9a-fA-F]{8}|[0-9a-fA-F]{6}|[0-9a-fA-F]{3,4})\b|\b(?:rgb|rgba|hsl|hsla|hwb|lab|lch|oklab|oklch)\s*\(|\b(?:white|black|red|green|blue|gray|grey|yellow|orange|purple|pink|brown|cyan|magenta|silver|maroon|navy|teal|olive|lime|aqua|fuchsia)\b/i

const why =
  "SVG paint must be a hex defined in the theme file (why: an off-palette fill is missed when the palette changes)."

function themeSolids(themeCSS) {
  const solids = new Set()
  for (const match of themeCSS.matchAll(hexColor)) solids.add(match[0].toLowerCase())
  return solids
}

function svgPaletteHits(text, solids) {
  const hits = []
  for (const match of text.matchAll(attr)) {
    const value = match[2]
    if (!rawColor.test(value)) continue
    const hexes = [...value.matchAll(hexColor)].map((item) => item[0].toLowerCase())
    const offPalette = hexes.some((hex) => !solids.has(hex))
    const nonHexLiteral = rawColor.test(value.replace(hexColor, ""))
    if (offPalette || nonHexLiteral || hexes.length === 0) {
      const line = text.slice(0, match.index).split("\n").length
      hits.push({ line, text: match[0], message: `${match[0].trim()}. ${why}` })
    }
  }
  return hits
}

function scanSvgFiles(root, themeFile) {
  const solids = themeSolids(fs.readFileSync(themeFile, "utf8"))
  const hits = []
  const walk = (dir) => {
    for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
      const full = path.join(dir, entry.name)
      if (entry.isDirectory()) {
        walk(full)
        continue
      }
      if (!entry.name.endsWith(".svg")) continue
      for (const hit of svgPaletteHits(fs.readFileSync(full, "utf8"), solids)) {
        hits.push({ ...hit, file: full })
      }
    }
  }
  walk(root)
  return hits
}

module.exports = { themeSolids, svgPaletteHits, scanSvgFiles, why }
