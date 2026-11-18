#!/usr/bin/env node
/**
 * Turns a full checkout (master) into the Raycast Store version, in place.
 *
 *   node scripts/raycast-release.mjs [dir]
 *
 * - deletes paths that don't belong in raycast/extensions (REMOVE_PATHS)
 * - drops Vicinae-only scripts and devDependencies from package.json
 * - removes every block between `raycast:strip-start` and `raycast:strip-end`
 *   marker lines (both markers included) in any text file. The markers can sit
 *   in any comment syntax: `<!-- raycast:strip-start -->`, `# raycast:strip-start`,
 *   `// raycast:strip-start`, ...
 *
 * Run `npm install --package-lock-only` afterwards so the lockfile matches.
 * Used by .github/workflows/raycast-branch.yml.
 */
import { existsSync, readdirSync, readFileSync, rmSync, statSync, writeFileSync } from "node:fs";
import { basename, join, relative } from "node:path";

const root = process.argv[2] ?? process.cwd();

const REMOVE_PATHS = ["example", ".vscode", ".github", "scripts"];
const REMOVE_DEV_DEPENDENCY = /^@vicinae\//;
const REMOVE_SCRIPT = /vicinae/i;

const STRIP_START = "raycast:strip-start";
const STRIP_END = "raycast:strip-end";
const TEXT_FILE = /\.(md|ts|tsx|js|jsx|mjs|cjs|json)$|^\.gitignore$/;
const SKIP_DIRS = new Set([".git", "node_modules"]);

function log(message) {
  console.log(`[raycast-release] ${message}`);
}

function removePaths() {
  for (const path of REMOVE_PATHS) {
    const full = join(root, path);
    if (existsSync(full)) {
      rmSync(full, { recursive: true, force: true });
      log(`removed ${path}`);
    }
  }
}

function cleanPackageJson() {
  const file = join(root, "package.json");
  const pkg = JSON.parse(readFileSync(file, "utf8"));
  for (const section of ["dependencies", "devDependencies"]) {
    for (const name of Object.keys(pkg[section] ?? {})) {
      if (REMOVE_DEV_DEPENDENCY.test(name)) {
        delete pkg[section][name];
        log(`package.json: removed ${section}.${name}`);
      }
    }
  }
  for (const [name, command] of Object.entries(pkg.scripts ?? {})) {
    if (REMOVE_SCRIPT.test(name) || REMOVE_SCRIPT.test(command)) {
      delete pkg.scripts[name];
      log(`package.json: removed scripts.${name}`);
    }
  }
  writeFileSync(file, JSON.stringify(pkg, null, 2) + "\n");
}

function stripMarkers(text, file) {
  const lines = text.split(/\r?\n/);
  const kept = [];
  let startLine = 0;
  lines.forEach((line, index) => {
    if (line.includes(STRIP_START)) {
      if (startLine) throw new Error(`${file}:${index + 1}: nested ${STRIP_START} (open since line ${startLine})`);
      startLine = index + 1;
    } else if (line.includes(STRIP_END)) {
      if (!startLine) throw new Error(`${file}:${index + 1}: ${STRIP_END} without ${STRIP_START}`);
      startLine = 0;
    } else if (!startLine) {
      kept.push(line);
    }
  });
  if (startLine) throw new Error(`${file}:${startLine}: ${STRIP_START} is never closed`);
  // Removing a block can leave a run of blank lines behind; collapse it.
  return kept.join("\n").replace(/\n{3,}/g, "\n\n");
}

function* walk(dir) {
  for (const name of readdirSync(dir)) {
    if (SKIP_DIRS.has(name)) continue;
    const full = join(dir, name);
    if (statSync(full).isDirectory()) yield* walk(full);
    else yield full;
  }
}

function stripAllFiles() {
  const leftovers = [];
  for (const full of walk(root)) {
    if (!TEXT_FILE.test(basename(full))) continue;
    const file = relative(root, full).replaceAll("\\", "/");
    const text = readFileSync(full, "utf8");
    if (text.includes(STRIP_START) || text.includes(STRIP_END)) {
      writeFileSync(full, stripMarkers(text, file));
      log(`stripped marked blocks in ${file}`);
    }
    if (file !== "package-lock.json" && /vicinae/i.test(readFileSync(full, "utf8"))) {
      leftovers.push(file);
    }
  }
  if (leftovers.length) {
    log(`note: these files still mention Vicinae: ${leftovers.join(", ")}`);
  }
}

removePaths();
cleanPackageJson();
stripAllFiles();
log("done");
