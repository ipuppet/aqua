export interface SiteConfig {
  title: string
  links: { title: string; href: string }[]
}

export const siteConfig: SiteConfig = {
  title: "Aqua",
  links: [{ title: "Home", href: "https://example.com" }]
}
