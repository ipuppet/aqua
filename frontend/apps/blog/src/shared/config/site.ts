export interface SiteConfig {
  title: string
  description: string
  author: string
  authorImage: string
  github: {
    name: string
    link: string
  }
}

export const siteConfig: SiteConfig = {
  title: "Aqua Blog",
  description: "A blog powered by Aqua.",
  author: "Aqua",
  authorImage: "/favicon-96x96.png",
  github: {
    name: "Aqua",
    link: "https://github.com/ipuppet/aqua"
  }
}
