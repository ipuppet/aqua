import globals from "globals"
import tsParser from "@typescript-eslint/parser"
import tsPlugin from "@typescript-eslint/eslint-plugin"
import vueParser from "vue-eslint-parser"
import vuePlugin from "eslint-plugin-vue"

export function createEslintConfig(projectDirname: string, extraIgnores: string[] = []) {
  const baseConfig = {
    ignores: ["src/auto-imports.d.ts", "src/components.d.ts", "**/*.d.ts"],
    languageOptions: {
      parserOptions: {
        projectService: true,
        tsconfigRootDir: projectDirname
      },
      globals: {
        ...globals.browser,
        ...globals.node
      }
    },
    plugins: {
      "@typescript-eslint": tsPlugin
    },
    rules: {
      "no-console": "error",
      "no-debugger": "error",
      "prefer-const": "error",
      "no-var": "error",
      "@typescript-eslint/no-unused-vars": [
        "error",
        {
          argsIgnorePattern: "^_",
          varsIgnorePattern: "^_"
        }
      ],
      "@typescript-eslint/prefer-nullish-coalescing": "warn"
    }
  }

  const defaultIgnores = [
    "src/auto-imports.d.ts",
    "src/components.d.ts",
    "src/vite-env.d.ts",
    "**/*.d.ts",
    "dist/",
    "dist-dev/"
  ]

  return [
    {
      ignores: Array.from(new Set([...defaultIgnores, ...extraIgnores]))
    },
    {
      ...baseConfig,
      files: ["src/**/*.{js,mjs,cjs,ts,tsx}"],
      languageOptions: {
        ...baseConfig.languageOptions,
        parser: tsParser
      },
      rules: {
        ...baseConfig.rules,
        "no-duplicate-imports": "off"
      }
    },
    {
      ...baseConfig,
      files: ["src/**/*.vue"],
      languageOptions: {
        ...baseConfig.languageOptions,
        parser: vueParser,
        parserOptions: {
          ...baseConfig.languageOptions.parserOptions,
          parser: tsParser,
          extraFileExtensions: [".vue"]
        }
      },
      plugins: {
        ...baseConfig.plugins,
        vue: vuePlugin
      },
      rules: baseConfig.rules
    },
    {
      files: ["src/shared/**/*.{js,mjs,cjs,ts,tsx,vue}"],
      rules: {
        "no-restricted-imports": [
          "error",
          {
            patterns: [
              {
                group: ["#/app/**", "#/pages/**", "#/widgets/**", "#/features/**", "#/entities/**"],
                message: "shared 层不能依赖上层分层（app/pages/widgets/features/entities）。"
              }
            ]
          }
        ]
      }
    },
    {
      files: ["src/entities/**/*.{ts,tsx,vue}"],
      rules: {
        "no-restricted-imports": [
          "error",
          {
            patterns: [
              {
                group: ["#/app/**", "#/pages/**", "#/widgets/**", "#/features/**"],
                message: "entities 层只能依赖 shared（以及同层）。"
              }
            ]
          }
        ]
      }
    },
    {
      files: ["src/features/**/*.{ts,tsx,vue}"],
      rules: {
        "no-restricted-imports": [
          "error",
          {
            patterns: [
              {
                group: ["#/app/**", "#/pages/**", "#/widgets/**"],
                message: "features 层不能依赖 app/pages/widgets。"
              }
            ]
          }
        ]
      }
    },
    {
      files: ["src/widgets/**/*.{ts,tsx,vue}"],
      rules: {
        "no-restricted-imports": [
          "error",
          {
            patterns: [
              {
                group: ["#/app/**", "#/pages/**"],
                message: "widgets 层不能依赖 app/pages。"
              }
            ]
          }
        ]
      }
    },
    {
      files: ["src/pages/**/*.{ts,tsx,vue}"],
      rules: {
        "no-restricted-imports": [
          "error",
          {
            patterns: [
              {
                group: ["#/app/**"],
                message: "pages 层禁止依赖 app。"
              }
            ]
          }
        ]
      }
    }
  ]
}
