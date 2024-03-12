<template>
  <VApp id="inspire">
    <VAppBar
      flat
      class="app-bar"
      :color="route.name === 'Home' ? 'transparent' : undefined"
    >
      <VContainer
        fluid
        class="fill-height d-flex align-center justify-space-between"
      >
        <!-- <VBtn
          class="route-links"
          v-for="{ route, text } in links"
          :key="route"
          :to="{ name: route }"
        >
          {{ text.toUpperCase() }}
        </VBtn> -->
        <RouterLink :to="{ name: 'Home' }" style="display: flex">
          <MainLogo />
          <!-- <VImg
            src="/images/main-logo.svg"
            :style="{
              filter: theme.current.value.dark
                ? 'invert(1) hue-rotate(180deg)'
                : '',
            }"
          /> -->
        </RouterLink>

        <VMenu :close-on-content-click="false" open-on-hover>
          <template #activator="{ props }">
            <VBtn
              :icon="mdiInvertColors"
              v-bind="props"
              color=""
              class="ml-auto"
            />
          </template>
          <VCard class="pa-1 theme-card">
            <VColorPicker
              elevation="0"
              v-model="color"
              hide-canvas
              :modes="['hex']"
              @update:model-value="debounce(setThemeColour, 500)()"
            />
            <VBtn
              class="col-select"
              @click="toggleTheme"
              :style="{
                color: `${
                  theme.current.value.dark ? 'black' : 'white'
                } !important`,
              }"
            >
              {{ theme.global.name.value === "darkTheme" ? "Light" : "Dark" }}
              mode
            </VBtn>
          </VCard>
        </VMenu>

        <VMenu
          :close-on-content-click="userStore.loggedIn ? true : false"
          @update:model-value="
            (opening) => {
              if (!opening) isRegistering = false;
            }
          "
        >
          <template #activator="{ props }">
            <VBtn v-if="userStore.loggedIn" class="avatar" icon v-bind="props">
              <VAvatar class="text-primary" :icon="mdiAccountCircle" />
            </VBtn>
            <VBtn v-else id="login-btn" v-bind="props" color="">Login</VBtn>
          </template>
          <VCard
            :min-width="userStore.loggedIn ? '15rem' : '20rem'"
            class="user-card"
          >
            <template v-if="userStore.loggedIn">
              <form>
                <v-list>
                  <v-list-item
                    :title="userStore.name"
                    :subtitle="
                      userStore.user?.firstname
                        ? userStore.user?.email
                        : undefined
                    "
                  />
                </v-list>

                <v-divider></v-divider>

                <v-list>
                  <v-list-item
                    v-for="{ title, icon } in profileLinks"
                    :key="title"
                    :to="{ name: title }"
                    active-color="primary"
                    :title="title"
                    :prepend-icon="icon"
                  />
                </v-list>
                <VBtn @click="userStore.logout">Logout</VBtn>
              </form>
            </template>
            <template v-else>
              <Transition name="login-fade" mode="out-in">
                <VForm v-if="isRegistering" @submit.prevent="register">
                  <VTextField
                    v-for="(_, label) in loginDetails"
                    :key="label"
                    :label="startCase(label)"
                    v-model="loginDetails[label]"
                    :name="label"
                    :type="
                      label.toLowerCase().includes('password')
                        ? 'password'
                        : 'text'
                    "
                    :rules="[
                      rules.required,
                      ...(label.toLowerCase().includes('password')
                        ? [rules.password, rules.passwordMatch]
                        : label === 'email'
                        ? [rules.email]
                        : []),
                    ]"
                  />
                  <VBtn type="submit">Register</VBtn>
                  <VBtn
                    variant="text"
                    color="gray"
                    @click="isRegistering = false"
                  >
                    Cancel
                  </VBtn>
                </VForm>

                <VForm v-else @submit.prevent="login">
                  <VTextField
                    label="Email"
                    name="email"
                    v-model="loginDetails.email"
                    type="email"
                    :rules="[rules.required, rules.email]"
                    autofocus
                  />
                  <VTextField
                    label="Password"
                    name="password"
                    v-model="loginDetails.password"
                    type="password"
                    :rules="[rules.required, rules.password]"
                  />
                  <VBtn type="submit">Login</VBtn>
                  <VBtn
                    variant="text"
                    color="gray"
                    @click="isRegistering = true"
                  >
                    Register
                  </VBtn>
                </VForm>
              </Transition>
            </template>
          </VCard>
        </VMenu>
      </VContainer>
    </VAppBar>

    <VMain
      class="bg-grey"
      :style="{
        background: `linear-gradient(0deg, ${theme.current.value.colors.background}, ${theme.current.value.colors.background}A1) center/cover fixed, center/cover fixed url(/images/retailers.jpg)`,
      }"
    >
      <RouterView v-slot="{ Component }">
        <template v-if="Component">
          <Transition name="fade" mode="out-in">
            <Suspense>
              <!-- main content -->
              <component :is="Component"></component>
            </Suspense>
          </Transition>
        </template>
      </RouterView>
    </VMain>
  </VApp>
</template>

<script setup lang="ts">
import {
  mdiAccountCircle,
  mdiAccountCogOutline,
  mdiInvertColors,
  mdiWardrobeOutline,
} from "@mdi/js";
import { useStorage } from "@vueuse/core";
import { debounce, startCase } from "lodash";
import { reactive, ref } from "vue";
import { RouterView, useRoute } from "vue-router";
import { useTheme } from "vuetify/lib/framework.mjs";
import { useUserStore } from "./stores/user";
import MainLogo from "@/assets/images/MainLogo.vue";

const theme = useTheme();
const color = ref("#4befa6");
const defaultColour = useStorage("primaryColour", "#4befa6");

const setThemeColour = () => {
  theme.themes.value.primaryTheme.colors.primary = color.value;
  theme.themes.value.darkTheme.colors.primary = color.value;
  defaultColour.value = color.value;
  const logo = document.getElementById("path0");

  if (logo) logo.style.fill = color.value;
};

const toggleTheme = () => {
  const newTheme = theme.current.value.dark ? "primaryTheme" : "darkTheme";

  theme.global.name.value = newTheme;

  localStorage.setItem("defaultTheme", newTheme);

  const meta = document.getElementsByName("theme-color")[0];
  meta.setAttribute("name", "theme-color");
  meta.setAttribute(
    "content",
    newTheme === "darkTheme" ? "#212121" : "#ffffff"
  );
};
const route = useRoute();
const userStore = useUserStore();
const links = [
  {
    route: "Home",
    text: "start now",
  },
];

const profileLinks = [
  {
    title: "Wardrobe",
    icon: mdiWardrobeOutline,
  },
  {
    title: "Profile",
    icon: mdiAccountCogOutline,
  },
];

const loginDetails = reactive({
  firstName: "",
  lastName: "",
  email: "",
  username: "",
  password: "",
  confirmPassword: "",
});

const isRegistering = ref(false);

const rules = {
  required: (v: string) => !!v || "Field is required",
  email: (v: string) => {
    const pattern =
      /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    return pattern.test(v) || "Invalid e-mail.";
  },
  password: (v: string) =>
    v.length > 6 || v.length < 20 || "Password must be between 6-20 chars",
  passwordMatch: (v: string) =>
    v === loginDetails.password || "Passwords do not match",
};
const login = () => {
  userStore.login({
    email: loginDetails.email,
    password: loginDetails.password,
  });
};

const register = () => {
  const { firstName, lastName, username, email, password } = loginDetails;

  userStore.register({
    firstName,
    lastName,
    username,
    email,
    password,
  });
};

userStore.loggedIn && userStore.getUser();
</script>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease !important;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.login-fade-enter-active,
.login-fade-leave-active {
  transition: opacity 0.1s ease !important;
}

.login-fade-enter-from,
.login-fade-leave-to {
  opacity: 0;
}

* {
  font-family: "Alexandria", sans-serif;
}
</style>

<style scoped lang="scss">
.theme-card {
  display: flex;
  flex-flow: column;
  gap: 5px;

  .col-select {
    background-color: rgba(
      var(--v-theme-on-background),
      var(--v-high-emphasis-opacity)
    ) !important;
  }
}
.app-bar {
  transition: 0.5s background !important;

  .route-links {
    font-weight: bold;

    &.v-btn--active {
      color: rgb(var(--v-theme-primary));

      &.v-theme--primaryTheme {
        filter: brightness(0.9);
      }
    }
  }
  .avatar {
    font-size: 6rem;
  }
}

.v-main {
  display: flex;
  align-items: center;
}

.user-card {
  padding: 10px;

  form {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }
}
</style>
