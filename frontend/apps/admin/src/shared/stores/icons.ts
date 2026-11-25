import { defineStore } from "pinia"

// 使用 Material Design Icons (mdi)
import IAboutMe from "virtual:icons/mdi/account-details"
import IArchive from "virtual:icons/mdi/archive"
import IArticles from "virtual:icons/mdi/file-document-multiple"
import IBanner from "virtual:icons/mdi/image"
import IBlog from "virtual:icons/mdi/post"
import ICheck from "virtual:icons/mdi/check"
import IClose from "virtual:icons/mdi/close"
import IDarkMode from "virtual:icons/mdi/theme-light-dark"
import IDown from "virtual:icons/mdi/chevron-down"
import IFold from "virtual:icons/mdi/menu-open"
import IHeart from "virtual:icons/mdi/heart"
import IHome from "virtual:icons/mdi/home"
import IInfo from "virtual:icons/mdi/information"
import ILink from "virtual:icons/mdi/link"
import IList from "virtual:icons/mdi/format-list-bulleted"
import ILock from "virtual:icons/mdi/lock"
import ILogout from "virtual:icons/mdi/logout"
import IPermission from "virtual:icons/mdi/shield-account"
import IPlus from "virtual:icons/mdi/plus"
import IRefresh from "virtual:icons/mdi/refresh"
import IRole from "virtual:icons/mdi/account-group"
import ISearch from "virtual:icons/mdi/magnify"
import ITag from "virtual:icons/mdi/tag"
import ITrash from "virtual:icons/mdi/delete"
import IUnfold from "virtual:icons/mdi/menu"
import IUnlock from "virtual:icons/mdi/lock-open"
import IUp from "virtual:icons/mdi/chevron-up"
import IUserSet from "virtual:icons/mdi/account-cog"
import IUser from "virtual:icons/mdi/account"
import IFileManager from "virtual:icons/mdi/folder-multiple"
import IImageManager from "virtual:icons/mdi/image-multiple"

export const useIconsStore = defineStore("icons", {
  state: () => ({
    AboutMe: IAboutMe,
    Archive: IArchive,
    Articles: IArticles,
    Banner: IBanner,
    Blog: IBlog,
    Check: ICheck,
    Close: IClose,
    DarkMode: IDarkMode,
    Down: IDown,
    Fold: IFold,
    Heart: IHeart,
    Home: IHome,
    Info: IInfo,
    Link: ILink,
    List: IList,
    Lock: ILock,
    Logout: ILogout,
    Permission: IPermission,
    Plus: IPlus,
    Refresh: IRefresh,
    Role: IRole,
    Search: ISearch,
    Tag: ITag,
    Trash: ITrash,
    Unfold: IUnfold,
    Unlock: IUnlock,
    Up: IUp,
    UserSet: IUserSet,
    User: IUser,
    FileManager: IFileManager,
    ImageManager: IImageManager
  })
})
