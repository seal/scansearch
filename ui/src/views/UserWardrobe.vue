<template>
  <VContainer>
    <VSheet class="pa-2" elevation="2" rounded="lg">
      <VRow>
        <VCol>
          <h1 class="title">Wardrobe</h1>
        </VCol>
      </VRow>
      <VRow>
        <VCol v-if="wardrobe.length === 0" class="no-items">
          <h2>There's nothing here!</h2>
          <p>Add something to your wardrobe to get started</p>
        </VCol>
        <VCol v-else class="results-list">
          <VHover
            v-for="item in wardrobe"
            :key="item.ID"
            v-slot="{ isHovering, props }"
          >
            <VCard
              class="result-card"
              v-bind="props"
              :elevation="isHovering ? 8 : 0"
              hover
            >
              <VImg
                class="product-thumbnail"
                :src="item.thumbnail"
                height="50%"
              />
              <VCardTitle tag="h6">
                {{ item.title }}
              </VCardTitle>
              <VCardSubtitle>
                {{ item.source }}
              </VCardSubtitle>
              <div class="content">
                <VCardText>
                  {{ item.price }}
                </VCardText>
                <VRating
                  v-model="item.rating"
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
                  :href="item.link"
                />
                <VBtn
                  :icon="mdiClose"
                  size="small"
                  @click.stop="deleteProduct(item.ID)"
                />
              </div>
            </VCard>
          </VHover>
        </VCol>
      </VRow>
    </VSheet>
  </VContainer>
</template>

<script setup lang="ts">
import axiosInstance from "@/axiosInstance";
import displayAlert from "@/helpers/displayAlert";
import handleAxiosError from "@/helpers/handleAxiosError";
import type { IWardrobeItem } from "@/types/responses.model";
import { mdiClose, mdiOpenInNew } from "@mdi/js";
import { ref } from "vue";

const wardrobe = ref<IWardrobeItem[]>([]);

try {
  const { data } = await axiosInstance.get<IWardrobeItem[]>("/wardrobe");
  wardrobe.value = data;
} catch (err) {
  handleAxiosError(err);
}

const deleteProduct = async (id: number) => {
  try {
    await axiosInstance.delete("/wardrobe", { data: [id] });
    displayAlert("Product removed from wardrobe", { type: "success" });
    wardrobe.value = wardrobe.value?.filter(({ ID }) => id !== ID);
  } catch (err) {
    handleAxiosError(err);
  }
};
</script>

<style lang="scss" scoped>
.v-container {
  margin: auto;

  .v-sheet {
    .no-items {
      text-align: center;
      padding: 100px 16px;
    }

    .title {
      text-align: center;
    }

    .results-list {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
      gap: 10px;

      .result-card {
        aspect-ratio: 1;
        padding: 16px;
        display: flex;
        flex-direction: column;

        .product-thumbnail {
          background: #fff;
        }

        & > * {
          padding-left: 0;
          padding-right: 0;
        }

        &:hover {
          .v-card-title {
            color: rgb(var(--v-theme-primary));
          }
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
  }
}
</style>
