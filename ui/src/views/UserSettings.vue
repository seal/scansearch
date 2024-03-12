<template>
  <VContainer>
    <VRow>
      <VCol>
        <VSheet class="sheet" rounded="xl" elevation="10">
          <h1>Account</h1>
          <div class="forms">
            <VForm @submit.prevent="userStore.updateUser(userDetails)">
              <h2>Your Details:</h2>
              <VTextField
                label="First Name"
                :rules="[rules.required]"
                type="text"
                v-model="userDetails.firstname"
              />
              <VTextField
                label="Last Name"
                :rules="[rules.required]"
                type="text"
                v-model="userDetails.lastname"
              />
              <VTextField
                label="Username"
                :rules="[rules.required]"
                type="text"
                v-model="userDetails.username"
              />
              <VTextField
                label="Email"
                :rules="[rules.required, rules.email]"
                type="email"
                v-model="userDetails.email"
              />
              <VBtn type="submit">Save</VBtn>
            </VForm>

            <VForm
              @submit.prevent="
                userStore.updatePassword(
                  userDetails.password,
                  userDetails.newPassword
                )
              "
            >
              <h2>Password</h2>
              <VTextField
                label="Old Password"
                :rules="[rules.required, rules.password]"
                type="password"
                v-model="userDetails.password"
              />
              <VTextField
                label="New Password"
                :rules="[rules.required, rules.password]"
                type="password"
                v-model="userDetails.newPassword"
              />
              <VTextField
                label="Confirm Password"
                :rules="[rules.required, rules.password, rules.passwordMatch]"
                type="password"
                v-model="userDetails.confirmPassword"
              />
              <VBtn type="submit">Save</VBtn>
            </VForm>
          </div>
        </VSheet>
      </VCol>
    </VRow>
  </VContainer>
</template>

<script setup lang="ts">
import { useUserStore } from "@/stores/user";
import { reactive } from "vue";

const userStore = useUserStore();

const userDetails = reactive({
  firstname: userStore.user?.firstname || "",
  lastname: userStore.user?.lastname || "",
  email: userStore.user?.email || "",
  username: userStore.user?.username || "",
  password: "",
  confirmPassword: "",
  newPassword: "",
});

const rules = {
  required: (v: string) => !!v || "Field is required",
  email: (v: string) => {
    const pattern =
      /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    return pattern.test(v) || "Invalid e-mail.";
  },
  password: (v: string | any[]) =>
    v.length > 6 || v.length < 20 || "Password must be between 6-20 chars",
  passwordMatch: (v: string) =>
    v === userDetails.newPassword || "Passwords do not match",
};
</script>

<style scoped lang="scss">
.sheet {
  padding: 1.5rem;

  h1 {
    padding-bottom: 1rem;
    text-align: center;
  }

  .forms {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-around;
    gap: 10px;
    .v-form {
      width: 30ch;
      display: flex;
      flex-flow: column;
      gap: 10px;

      .v-input {
        flex: 0;
      }
    }
  }
}
</style>
