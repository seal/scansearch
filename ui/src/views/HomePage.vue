<template>
  <VContainer>
    <VRow>
      <VSheet class="sheet" rounded="xl" elevation="10">
        <h3>SCAN: <span class="text-primary">SEARCH</span></h3>
        <h2 class="my-6">
          COMPARE PRICES ON 2.5 MILLION PRODUCTS FROM THE LARGEST UK RETAILERS
        </h2>
        <VForm @submit="search">
          <VTextField
            v-model="searchTerm"
            clearable
            label=""
            density="compact"
            autofocus
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
    </VRow>
  </VContainer>
</template>

<script setup lang="ts">
import displayAlert from "@/helpers/displayAlert";
import { ref } from "vue";
import { useRouter } from "vue-router";
import axiosInstance from "@/axiosInstance";
import handleAxiosError from "@/helpers/handleAxiosError";

const router = useRouter();

const props = defineProps<{ code?: string }>();

const searchTerm = ref("");
const loading = ref(false);

const search = async (e: Event) => {
  e.preventDefault();

  loading.value = true;
  router.push({ name: "Results", query: { query: searchTerm.value } });
};

if (props.code) {
  try {
    const { data } = await axiosInstance.get("/verifyemail", {
      params: {
        verificationCode: props.code,
      },
    });
    displayAlert(data.message, { type: "success" });
  } catch (err) {
    handleAxiosError(err);
  }
}
</script>

<style scoped lang="scss">
.sheet {
  padding: 3.5rem;
  max-width: 1000px;
  margin: auto;

  h2 {
    line-height: initial;
  }

  .v-form {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 10px;

    .v-input {
      min-width: 160px;
      flex-grow: 1000;
    }

    .v-btn {
      flex-grow: 1;
    }
  }
}
</style>
