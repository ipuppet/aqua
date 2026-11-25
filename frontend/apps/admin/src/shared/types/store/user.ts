export interface UserInfo {
  avatar: string
  nickname: string
  username?: string
  email?: string
  [key: string]: unknown
}

export interface UserStructure {
  [key: string]: unknown
}
