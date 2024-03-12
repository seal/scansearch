import { createPinia } from "pinia";
import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";

import "mosha-vue-toastify/dist/style.css";
import "vuetify/styles";

import vuetify from "./plugins/vuetify";

const pinia = createPinia();
const app = createApp(App);

app.use(router);
app.use(vuetify);
app.use(pinia);

app.mount("#app");
