import { createEslintConfig } from "@aqua/eslint-config"

export default createEslintConfig(import.meta.dirname, ["src/env.d.ts"])
