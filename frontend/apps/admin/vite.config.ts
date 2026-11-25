import { createViteConfig } from "@aqua/vite-config"

export default createViteConfig({
  rootDir: import.meta.dirname,
  componentDirs: ["src/shared/ui", "src/widgets/layout"],
  iconCollections: ["mdi", "ic", "heroicons-outline", "ep"]
})
