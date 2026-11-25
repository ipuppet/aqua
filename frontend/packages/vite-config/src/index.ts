import { defineConfig, type UserConfig } from "vite"
import Vue from "@vitejs/plugin-vue"
import AutoImport from "unplugin-auto-import/vite"
import Components from "unplugin-vue-components/vite"
import { ElementPlusResolver } from "unplugin-vue-components/resolvers"
import Icons from "unplugin-icons/vite"
import IconsResolver from "unplugin-icons/resolver"
import path from "node:path"

export type CreateViteConfigOptions = {
  /** Absolute path to the app root (usually import.meta.dirname). */
  rootDir: string
  /** unplugin-vue-components dirs; omit to use plugin defaults. */
  componentDirs?: string[]
  /** Iconify collections for IconsResolver. */
  iconCollections?: string[]
  /** Extra Icons() options (e.g. customCollections). */
  icons?: Record<string, unknown>
  /** Merge/override the generated Vite config. */
  override?: UserConfig
}

export function createViteConfig(options: CreateViteConfigOptions): UserConfig {
  const isProduction = process.env.NODE_ENV === "production"
  const iconCollections = options.iconCollections ?? ["mdi", "ic", "heroicons-outline", "ep"]

  const config: UserConfig = {
    envDir: "../../",
    build: {
      outDir: isProduction ? "dist" : "dist-dev"
    },
    resolve: {
      alias: {
        "#": path.resolve(options.rootDir, "src")
      }
    },
    plugins: [
      Vue(),
      Icons({
        compiler: "vue3",
        ...options.icons
      }),
      AutoImport({
        imports: ["vue"],
        resolvers: [ElementPlusResolver()],
        dts: "src/auto-imports.d.ts"
      }),
      Components({
        ...(options.componentDirs ? { dirs: options.componentDirs } : {}),
        resolvers: [
          ElementPlusResolver(),
          IconsResolver({
            enabledCollections: iconCollections
          })
        ],
        dts: "src/components.d.ts"
      })
    ]
  }

  return defineConfig({
    ...config,
    ...options.override,
    build: {
      ...config.build,
      ...options.override?.build
    },
    resolve: {
      ...config.resolve,
      ...options.override?.resolve,
      alias: {
        ...(config.resolve?.alias as Record<string, string> | undefined),
        ...(options.override?.resolve?.alias as Record<string, string> | undefined)
      }
    },
    plugins: [...(config.plugins ?? []), ...(options.override?.plugins ?? [])]
  })
}
