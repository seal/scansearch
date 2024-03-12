<template>
  <VRow>
    <VCol
      v-if="filters.length > 0"
      :cols="mdAndUp ? 3 : 12"
      class="filters-list"
      v-resize="() => (dialogOpen = false)"
    >
      <component
        :is="mdAndUp ? 'div' : VDialog"
        v-model="dialogOpen"
        scrollable
      >
        <template #activator="{ props }">
          <VBtn
            class="open-filters-btn"
            v-bind="props"
            :width="xs ? '100%' : 'unset'"
          >
            Show Filters
          </VBtn>
        </template>

        <VSheet :class="{ 'pa-1': !mdAndUp }" color="transparent" elevation="0">
          <div v-if="!mdAndUp" style="display: flex">
            <VBtn
              class="ma-auto mb-2"
              color="gray"
              icon
              @click="dialogOpen = false"
            >
              <VIcon>{{ mdiClose }}</VIcon>
            </VBtn>
          </div>

          <VExpansionPanels variant="popout">
            <VBtn
              class="sort-btn"
              :class="{ 'sort-up': !sortDown }"
              :prepend-icon="sortDown === null ? mdiMinus : mdiArrowDown"
              @click="
                sortDown = sortDown ? false : sortDown === null ? true : null
              "
            >
              Sort
            </VBtn>
            <VExpansionPanel
              v-for="(filter, i) in results.filters"
              :key="i"
              :title="filter.type"
            >
              <VExpansionPanelText
                class="filter-category"
                :id="`${filter.type}-filter`"
              >
                <VCheckbox
                  :key="j"
                  v-for="({ text, tbs }, j) in filter.options"
                  :label="text"
                  :value="tbs"
                  color="primary"
                  density="compact"
                  v-model="filters[i]"
                  multiple
                  hide-details
                />
              </VExpansionPanelText>
            </VExpansionPanel>
            <div class="filter-btns">
              <VBtn color="grey" @click="clearFilters">Clear</VBtn>
              <VBtn :loading="loading" @click="saveFilters">Save</VBtn>
            </div>
          </VExpansionPanels>
        </VSheet>
      </component>
    </VCol>
    <VCol
      :class="hasResults ? 'results-list' : 'no-results'"
      :cols="filters.length > 0 && mdAndUp ? 9 : 12"
    >
      <VSheet v-if="!hasResults" elevation="2" rounded="lg">
        <h2>No Results Found</h2>
      </VSheet>
      <VHover
        v-else
        v-for="result in results.ShoppingResults?.sort(sortFunction)"
        :key="result.position"
        v-slot="{ isHovering, props }"
      >
        <VSheet min-height="250px" color="transparent">
          <component
            :is="result.position < 10 ? 'div' : VLazy"
            v-model="result.isActive"
          >
            <VCard
              class="result-card"
              v-bind="props"
              :elevation="isHovering ? 14 : 2"
            >
              <VImg
                class="product-thumbnail"
                :src="result.thumbnail"
                height="50%"
              />
              <VCardTitle tag="h6">
                {{ result.title }}
              </VCardTitle>
              <VCardSubtitle>
                {{ result.source }}
              </VCardSubtitle>
              <div class="content">
                <VCardText>
                  {{ result.price }}
                </VCardText>
                <VRating
                  v-model="result.rating"
                  size="x-small"
                  color="primary"
                  readonly
                  half-increments
                />
              </div>
              <div class="btns">
                <VBtn
                  :icon="mdiOpenInNew"
                  color="grey"
                  size="small"
                  :href="result.link"
                />
                <VTooltip
                  :disabled="userStore.loggedIn"
                  text="Login to save items"
                >
                  <template #activator="{ props }">
                    <VBtn
                      :color="userStore.loggedIn ? 'primary' : 'grey'"
                      v-bind="props"
                      :icon="mdiHeart"
                      size="small"
                      @click.stop="saveProduct(result)"
                    />
                  </template>
                </VTooltip>
              </div>
            </VCard>
          </component>
        </VSheet>
      </VHover>
    </VCol>
  </VRow>
</template>

<script setup lang="ts">
import axiosInstance from "@/axiosInstance";
import type {
  ISearchFilter,
  ISearchResponse,
  IShoppingResult,
} from "@/types/responses.model";
import { computed, ref, toRefs } from "vue";
import { useDisplay } from "vuetify";
import { VDialog, VLazy } from "vuetify/components";
import {
  mdiArrowDown,
  mdiClose,
  mdiHeart,
  mdiMinus,
  mdiOpenInNew,
} from "@mdi/js";
import sizeFilters from "@/assets/sizeTbs.json";
import retailerFilters from "@/assets/retailersTbs.json";

const userStore = useUserStore();
const { mdAndUp, xs } = useDisplay();
const dialogOpen = ref(false);

const results = ref<ISearchResponse>({
    ShoppingResults: [],
    filters: [],
    message: "",
    success: false,
  }),
  filters = ref<string[][]>([]),
  sortDown = ref<boolean | null>(true),
  loading = ref(false),
  hasResults = computed(() => !!results.value.ShoppingResults?.length),
  props = defineProps<{
    query: string;
  }>(),
  emit = defineEmits<{
    (e: "searchComplete"): void;
  }>(),
  { query } = toRefs(props);

const clearFilters = () => {
  filters.value = results.value.filters?.map((x) => []) || [];
  searchProducts();
};

const saveFilters = () => {
  dialogOpen.value = false;
  searchProducts();
};

const formatFilters = (filtersArr: ISearchFilter[]) => {
  const getColourFromTbs = (tbs: string): string => {
    const colourStr = tbs.split("color_val:").pop()?.split(",")[0] || "unknown";
    return colourStr[0].toUpperCase() + colourStr.substring(1);
  };

  const colourIndex = filtersArr.findIndex(({ type }) => type === "Colour"),
    retailerIndex = filtersArr.findIndex(({ type }) => type === "Seller");

  filtersArr[colourIndex]?.options.forEach(
    ({ tbs }, i, arr) => (arr[i].text = getColourFromTbs(tbs))
  );

  filtersArr[retailerIndex] = retailerFilters;

  filtersArr = filtersArr.filter(({ options }) =>
    options.some(({ text }) => !!text)
  );

  filtersArr.splice(1, 0, sizeFilters);

  return filtersArr;
};

const sortFunction = (a: IShoppingResult, b: IShoppingResult) =>
  sortDown.value === null
    ? a.position - b.position
    : sortDown.value
    ? a.extracted_price - b.extracted_price
    : b.extracted_price - a.extracted_price;

const searchProducts = async () => {
  if (query.value) {
    loading.value = true;
    let data: ISearchResponse;
    const filterArray =
      filters.value.length > 0
        ? filters.value
            .reduce((a, b) => a.concat(b))
            .map((str) => str.substring(5))
        : [];

    try {
      if (filterArray.length === 0) {
        ({ data } = await axiosInstance.get("/tbs", {
          params: {
            query: query.value,
          },
        }));

        if (data.filters) {
          data.filters = formatFilters(data.filters);
          filters.value = data.filters.map(() => []);
        }
        results.value = data;
      } else {
        ({ data } = await axiosInstance.post("/tbs", filterArray, {
          params: {
            query: query.value,
          },
        }));

        results.value.ShoppingResults = data.ShoppingResults;
      }

      window.scrollTo(0, 0);
      loading.value = false;
    } catch (err) {
      handleAxiosError(err);
      loading.value = false;
    }
    emit("searchComplete");
  }
};

await searchProducts();

const saveProduct = async (product: IShoppingResult) => {
  if (userStore.loggedIn) {
    try {
      await axiosInstance.post("/wardrobe", product);

      displayAlert("Product saved to wardrobe", { type: "success" });
    } catch (err) {
      handleAxiosError(err);
    }
  } else {
    const loginCard = document.getElementsByClassName("user-card")[0];

    !loginCard && document.getElementById("login-btn")?.click();
  }
};
</script>

<script lang="ts">
import { defineComponent } from "vue";
import type { RouteLocationNormalized } from "vue-router";
import handleAxiosError from "@/helpers/handleAxiosError";
import displayAlert from "@/helpers/displayAlert";
import { useUserStore } from "@/stores/user";

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
.v-overlay {
  :deep(.v-overlay__content) {
    margin: 0;
    max-height: 100%;
    max-width: 100%;
    width: auto;
  }
}

.filters-list {
  display: flex;

  .open-filters-btn {
    margin: auto;
  }
}

.filter-category {
  max-height: 20ch;
  overflow-y: auto;

  &#Retailers-filter {
    :deep(.v-expansion-panel-text__wrapper) {
      padding: 0;
      padding-bottom: 16px;
    }
  }
}

.sort-btn {
  margin-bottom: 10px;
  max-width: calc(100% - 32px);
  flex: 1 0 100%;

  :deep(.v-icon) {
    transition: transform 0.5s;
  }

  &.sort-up {
    :deep(.v-icon) {
      transform: rotate(-180deg);
    }
  }
}
.filter-btns {
  display: flex;
  justify-content: space-around;
  flex-grow: 1;
  max-width: calc(100% - 32px);
  margin-top: 10px;
}

.no-results {
  display: flex;

  .v-sheet {
    margin: auto;
    margin-top: 20%;
    padding: 16px;
  }
}

.results-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(275px, 1fr));
  gap: 10px;

  .result-card {
    aspect-ratio: 1;
    padding: 16px;
    display: flex;
    flex-direction: column;

    & > * {
      padding-left: 0;
      padding-right: 0;
    }

    &:hover {
      .v-card-title {
        color: rgb(var(--v-theme-primary));
      }
    }

    .product-thumbnail {
      background: #fff;
    }

    .v-card-title {
      padding-bottom: 0;
      transition: 0.28s color;
    }

    .content {
      display: flex;
      padding: 0;

      .v-card-text {
        padding: 0;
        align-self: center;
      }
    }

    .btns {
      display: flex;
      justify-content: space-between;
    }
  }
}
</style>
