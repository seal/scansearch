<template>
  <VContainer>
    <VRow>
      <VSpacer />
      <VCol :cols="smAndDown ? 12 : 6">
        <VSheet class="search-bar-container" rounded="lg" elevation="2">
          <VForm @submit.prevent="search">
            <VTextField
              v-model="searchTerm"
              clearable
              label=""
              density="compact"
              :disabled="loading"
              placeholder="Search for a product"
              hide-details="auto"
            />
            <VBtn
              class="text-white font-weight-bold flex-sm-grow-0"
              color="primary"
              :disabled="!searchTerm"
              rounded="lg"
              type="submit"
              :loading="loading"
            >
              SEARCH
            </VBtn>
          </VForm>
        </VSheet>
      </VCol>
      <VSpacer />
    </VRow>
    <RouterView
      v-slot="{ Component, route }"
      v-on:searchComplete="loading = false"
    >
      <template v-if="Component">
        <Transition name="fade" mode="out-in">
          <Suspense>
            <!-- main content -->
            <component :is="Component" :key="route.fullPath"></component>
          </Suspense>
        </Transition>
      </template>
    </RouterView>
  </VContainer>
</template>

<script setup lang="ts">
import { ref, toRefs } from "vue";
import { useRouter } from "vue-router";
import { useDisplay } from "vuetify";

const { smAndDown } = useDisplay();
const router = useRouter();

const props = defineProps<{
  query: string;
}>();

const { query } = toRefs(props);

const searchTerm = ref(query.value);
const loading = ref(true);

const search = async () => {
  if (searchTerm.value === query.value) return;

  loading.value = true;
  router.push({ name: "Results", query: { query: searchTerm.value } });
};
</script>

<script lang="ts">
import { defineComponent } from "vue";
import type { RouteLocationNormalized } from "vue-router";
export default defineComponent({
  beforeRouteEnter(to: RouteLocationNormalized) {
    if (!to.query.query) {
      return { name: "Home" };
    } else {
      return true;
    }
  },
});
</script>

<style scoped lang="scss">
.v-container {
  margin-bottom: auto;
  max-width: none;
  .search-bar-container {
    padding: 16px;

    .v-form {
      display: flex;
      flex-wrap: wrap;
      align-items: center;
      gap: 10px;

      .v-input {
        min-width: 250px;
      }

      .v-btn {
        flex-grow: 1;
      }
    }
  }
}
</style>
