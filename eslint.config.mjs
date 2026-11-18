import { defineConfig } from "eslint/config";
import raycastConfig from "@raycast/eslint-config";

export default defineConfig([
  {
    ignores: [
      "raycast-env.d.ts",
      // raycast:strip-start
      "vicinae-env.d.ts",
      // raycast:strip-end
    ],
  },
  ...raycastConfig,
]);
