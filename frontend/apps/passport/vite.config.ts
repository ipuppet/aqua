import { createViteConfig } from "@aqua/vite-config"

export default createViteConfig({
  rootDir: import.meta.dirname,
  iconCollections: ["mdi", "ic", "heroicons-outline", "ep"]
})
