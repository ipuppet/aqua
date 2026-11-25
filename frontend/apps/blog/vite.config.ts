import { FileSystemIconLoader } from "unplugin-icons/loaders"
import { createViteConfig } from "@aqua/vite-config"

export default createViteConfig({
  rootDir: import.meta.dirname,
  iconCollections: ["tabler", "custom"],
  icons: {
    customCollections: {
      custom: FileSystemIconLoader("./src/assets", svg =>
        svg
          .replace(/ width="[^"]*"/, "")
          .replace(/ height="[^"]*"/, "")
          .replace(/^<svg /, '<svg fill="currentColor" width="1em" height="1em" ')
      )
    }
  }
})
