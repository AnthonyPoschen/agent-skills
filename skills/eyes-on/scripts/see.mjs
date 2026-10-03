#!/usr/bin/env node
// Save a viewport × color-scheme matrix of PNGs for the eyes-on skill.
// The agent must open each printed path with an image reader before judging.

import { createRequire } from "node:module"
import { spawnSync } from "node:child_process"
import { existsSync, mkdirSync, writeFileSync } from "node:fs"
import { join } from "node:path"

const args = process.argv.slice(2)

if (args.includes("--help") || args.includes("-h")) {
  process.stdout.write(`Usage: node see.mjs --url URL --out DIR [options]

Options:
  --size WxH     Viewport. Repeat for more than one. Default: 1280x800 and 390x844
  --scheme LIST  Comma-separated light and/or dark. Default: light,dark
  --wait CSS     Wait until this selector is visible
  --scroll CSS   Scroll this selector into view before the shot
  --settle MS    Pause after load so fonts and layout settle. Default: 400
  --full         Also save a full-page PNG per size and scheme

Reads Playwright from NODE_PATH or EYES_ON_PLAYWRIGHT (package root).
Uses chromium on PATH, or EYES_ON_CHROMIUM.
`)
  process.exit(0)
}

function option(name) {
  const values = []
  for (let i = 0; i < args.length; i += 1) {
    if (args[i] === name) {
      const value = args[i + 1]
      if (value == null || value.startsWith("--")) {
        fail(`missing value for ${name}`)
      }
      values.push(value)
      i += 1
    }
  }
  return values
}

function fail(message) {
  process.stderr.write(`see.mjs: ${message}\n`)
  process.exit(1)
}

function parseSize(raw) {
  const match = /^(\d+)x(\d+)$/.exec(raw)
  if (match == null) {
    fail(`size must look like 1280x800, got ${raw}`)
  }
  return { width: Number(match[1]), height: Number(match[2]), label: raw }
}

function findPlaywrightRoot() {
  if (process.env.EYES_ON_PLAYWRIGHT) {
    return process.env.EYES_ON_PLAYWRIGHT
  }
  const starts = [process.cwd()]
  if (process.env.NODE_PATH) {
    starts.push(...process.env.NODE_PATH.split(":"))
  }
  const home = process.env.HOME
  if (home) {
    starts.push(join(home, ".local/share/mise/installs/npm-playwright/latest/node_modules"))
  }
  for (const start of starts) {
    let dir = start
    for (let depth = 0; depth < 6; depth += 1) {
      const root = join(dir, "node_modules", "playwright")
      const direct = join(dir, "playwright")
      if (existsSync(join(root, "package.json"))) {
        return root
      }
      if (existsSync(join(direct, "package.json"))) {
        return direct
      }
      const parent = join(dir, "..")
      if (parent === dir) {
        break
      }
      dir = parent
    }
  }
  return ""
}

function findChromium() {
  if (process.env.EYES_ON_CHROMIUM) {
    return process.env.EYES_ON_CHROMIUM
  }
  for (const name of ["chromium", "chromium-browser", "google-chrome", "google-chrome-stable"]) {
    const found = spawnSync("sh", ["-c", `command -v ${name}`], { encoding: "utf8" })
    const path = found.stdout.trim()
    if (found.status === 0 && path !== "") {
      return path
    }
  }
  return ""
}

const urlValues = option("--url")
const outValues = option("--out")
if (urlValues.length !== 1) {
  fail("pass one --url")
}
if (outValues.length !== 1) {
  fail("pass one --out")
}
const url = urlValues[0]
const outDir = outValues[0]

const sizeArgs = option("--size")
const sizes = (sizeArgs.length === 0 ? ["1280x800", "390x844"] : sizeArgs).map(parseSize)

const schemeArgs = option("--scheme")
const schemes = (schemeArgs.length === 0 ? ["light,dark"] : schemeArgs).join(",").split(",").map((item) => item.trim()).filter((item) => item !== "")
for (const scheme of schemes) {
  if (scheme !== "light" && scheme !== "dark") {
    fail(`scheme must be light or dark, got ${scheme}`)
  }
}

const waitArgs = option("--wait")
const scrollArgs = option("--scroll")
const settleArgs = option("--settle")
if (waitArgs.length > 1 || scrollArgs.length > 1 || settleArgs.length > 1) {
  fail("pass --wait, --scroll, and --settle at most once")
}
const waitSelector = waitArgs[0] ?? ""
const scrollSelector = scrollArgs[0] ?? ""
const settle = settleArgs.length === 0 ? 400 : Number(settleArgs[0])
if (!Number.isFinite(settle) || settle < 0) {
  fail("--settle must be a number of milliseconds")
}
const fullPage = args.includes("--full")

const playwrightRoot = findPlaywrightRoot()
if (playwrightRoot === "") {
  fail("Playwright was not found. Set EYES_ON_PLAYWRIGHT to the playwright package root.")
}

const require = createRequire(import.meta.url)
const { chromium } = require(join(playwrightRoot, "index.js"))
const executablePath = findChromium()

mkdirSync(outDir, { recursive: true })

const launchOptions = { headless: true }
if (executablePath !== "") {
  launchOptions.executablePath = executablePath
}

const browser = await chromium.launch(launchOptions)
const written = []

try {
  for (const scheme of schemes) {
    for (const size of sizes) {
      const context = await browser.newContext({
        colorScheme: scheme,
        viewport: { width: size.width, height: size.height },
        deviceScaleFactor: 1,
      })
      const page = await context.newPage()
      try {
        await page.goto(url, { waitUntil: "load", timeout: 20000 })
        if (waitSelector !== "") {
          await page.waitForSelector(waitSelector, { state: "visible", timeout: 20000 })
        }
        if (scrollSelector !== "") {
          await page.locator(scrollSelector).first().scrollIntoViewIfNeeded()
        }
        await new Promise((resolve) => {
          setTimeout(resolve, settle)
        })
        const viewportPath = join(outDir, `${size.label}-${scheme}.png`)
        await page.screenshot({ path: viewportPath, fullPage: false })
        written.push(viewportPath)
        if (fullPage) {
          const fullPath = join(outDir, `${size.label}-${scheme}-full.png`)
          await page.screenshot({ path: fullPath, fullPage: true })
          written.push(fullPath)
        }
      } finally {
        await context.close()
      }
    }
  }
} finally {
  await browser.close()
}

const lines = ["Read each PNG with an image reader before you judge.", ...written]
writeFileSync(join(outDir, "eyes.txt"), `${lines.join("\n")}\n`)
process.stdout.write(`${lines.join("\n")}\n`)
