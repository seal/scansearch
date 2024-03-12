import { createVuetify, type ThemeDefinition } from "vuetify";
import { aliases, mdi } from "vuetify/iconsets/mdi-svg";

const primaryColour = localStorage.getItem("primaryColour") || "#4befa6";
const prefersDark = window.matchMedia?.("(prefers-color-scheme: dark)").matches;
let defaultTheme = localStorage.getItem("defaultTheme");

if (!defaultTheme) defaultTheme = prefersDark ? "darkTheme" : "primaryTheme";

const meta = document.createElement("meta");
meta.setAttribute("name", "theme-color");
meta.setAttribute(
  "content",
  defaultTheme === "darkTheme" ? "#212121" : "#ffffff"
);

const container = document.getElementsByTagName("head")[0];
container.insertBefore(meta, container.firstChild);

const primaryTheme: ThemeDefinition = {
  dark: false,
  colors: {
    primary: primaryColour,
    "on-primary": "#FFFFFF",
  },
};

const darkTheme: ThemeDefinition = {
  dark: true,
  colors: {
    primary: primaryColour,
    "on-primary": "#FFFFFF",
  },
};

export default createVuetify({
  theme: {
    defaultTheme,
    themes: {
      primaryTheme,
      darkTheme,
    },
  },
  icons: {
    defaultSet: "mdi",
    aliases,
    sets: {
      mdi,
    },
  },
  defaults: {
    VTextField: {
      color: "primary",
      variant: "outlined",
      hideDetails: "auto",
    },
    VBtn: {
      color: "primary",
    },
  },
});
