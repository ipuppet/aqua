<template>
  <div class="router-view-container">
    <loading-ball-atom v-show="!isLoaded" class="loading" />
    <div class="block">
      <el-timeline>
        <el-timeline-item
          v-for="article in archive"
          :key="article.aid"
          :timestamp="new Date(article.createdAt).toLocaleString()"
          placement="top"
        >
          <router-link
            :to="'/detail/' + article.aid + '/'"
            :class="{ dark: darkModeStore.isDarkMode }"
            class="timeline-title link"
            >{{ article.title }}</router-link
          >
        </el-timeline-item>
      </el-timeline>
    </div>
  </div>
</template>
<script>
import LoadingBallAtom from "#/shared/ui/LoadingBallAtom.vue"
import { useDarkModeStore } from "#/shared/stores/darkMode"
import { sessionData } from "#/shared/lib/storage"
import { client } from "#/shared/lib/api"
export default {
  components: { LoadingBallAtom },
  setup() {
    const darkModeStore = useDarkModeStore()
    return { darkModeStore }
  },
  data: () => ({ archive: [], isLoaded: false }),
  methods: {
    async getArchive() {
      const archive = sessionData.get("archive")
      if (!archive) {
        const { data, error } = await client.GET("/api/archive")
        if (!error && data) {
          this.archive = data.articles
          sessionData.set("archive", this.archive)
          this.isLoaded = true
        } else {
          this.isLoaded = true
        }
      } else {
        this.archive = archive
        this.isLoaded = true
      }
    }
  },
  created() {
    document.title = "Archive"
    this.getArchive()
  }
}
</script>
<style scoped>
.timeline-title {
  font-size: 16px;
  margin-left: 3px;
}

.timeline-title:hover {
  margin-left: 7px;
}

.tags {
  font-size: 12px;
  margin: 0 5px;
}

.tags:hover {
  margin: 0 10px;
}
</style>
