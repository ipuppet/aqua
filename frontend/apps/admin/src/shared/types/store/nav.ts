export type NavButtonSize = "default" | "small" | "large" | string

export interface NavButton {
  text: string
  click: () => void
  type?: string
  icon?: string
  size?: NavButtonSize
}

export interface NavState {
  buttons: NavButton[]
  showSearch: boolean
  searchKeyword: string
  onSearch: ((keyword: string) => void) | null
  onClearSearch: (() => void) | null
}
