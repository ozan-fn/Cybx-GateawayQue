import { cpSync, existsSync, rmSync } from "node:fs";
import { execFileSync } from "node:child_process";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = join(dirname(fileURLToPath(import.meta.url)), "..");
const dashboardDir = join(root, "Dashboard");
const outputDir = join(dashboardDir, "out");
const backendDashboardDir = join(root, "Backend", "dashboard");

execFileSync("npx", ["--no-install", "next", "build"], {
  cwd: dashboardDir,
  stdio: "inherit",
  env: { ...process.env, NEXT_STATIC_EXPORT: "true" },
  shell: process.platform === "win32",
});

if (!existsSync(outputDir)) {
  throw new Error("Dashboard static export was not created");
}

rmSync(backendDashboardDir, { recursive: true, force: true });
cpSync(outputDir, backendDashboardDir, { recursive: true });
