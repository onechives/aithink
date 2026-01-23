import { createI18n } from "vue-i18n";
import zhCN from "./locales/zh-CN";
import enUS from "./locales/en-US";

// 检测浏览器语言
const getDefaultLocale = (): string => {
  const stored = localStorage.getItem("language");
  if (stored) return stored;

  const browserLang = navigator.language.toLowerCase();
  if (browserLang.startsWith("en")) return "en-US";
  if (browserLang.startsWith("zh")) return "zh-CN";
  return "en-US"; // 默认英文
};

const i18n = createI18n({
  legacy: false,
  locale: getDefaultLocale(),
  fallbackLocale: "en-US",
  messages: {
    "zh-CN": zhCN,
    "en-US": enUS,
  },
});

export default i18n;
