import { client } from "#/shared/lib/api"
import { sessionData } from "#/shared/lib/storage"
import { ref } from "vue"

export function useTags() {
  const tags = ref<any[]>([])

  async function loadTags() {
    const tagsFromSession = sessionData.get("tags")
    if (!tagsFromSession) {
      const { data, error } = await client.GET("/api/tags")
      if (!error && data?.tags) {
        sessionData.set("tags", data.tags)
        tags.value = data.tags
      }
    } else {
      tags.value = tagsFromSession ?? []
    }
  }

  return { tags, loadTags }
}
