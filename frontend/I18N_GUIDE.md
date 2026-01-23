# 多语言国际化指南

## 概述

该项目使用 `vue-i18n` 实现多语言支持。目前支持中文（zh-CN）和英文（en-US）。

## 功能特性

- ✅ 根据浏览器语言自动加载对应语言
- ✅ 右上角语言切换按钮
- ✅ 语言偏好保存在 localStorage
- ✅ 所有UI文本都支持多语言

## 项目结构

```
frontend/src/
├── i18n/
│   ├── index.ts                 # i18n 配置入口
│   └── locales/
│       ├── zh-CN.ts             # 中文翻译文件
│       └── en-US.ts             # 英文翻译文件
└── composables/
    └── useI18n.ts              # i18n 组合函数
```

## 使用方法

### 1. 在模板中使用翻译

```vue
<template>
  <div>
    <!-- 使用 $t 调用翻译 -->
    <h1>{{ $t('home.title') }}</h1>
    <p>{{ $t('nav.articles') }}</p>
  </div>
</template>
```

### 2. 在脚本中使用翻译

```typescript
<script setup lang="ts">
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const title = t('home.title');
const description = t('home.noArticles');
</script>
```

### 3. 动态切换语言

```typescript
import { useI18n } from "vue-i18n";

const { locale } = useI18n();

// 切换到英文
locale.value = 'en-US';

// 切换到中文
locale.value = 'zh-CN';

// 保存到 localStorage
localStorage.setItem('language', 'en-US');
```

## 添加新的翻译

### 1. 编辑翻译文件

在 `src/i18n/locales/zh-CN.ts` 中添加新的翻译：

```typescript
export default {
  // ... 现有翻译
  newFeature: {
    title: "新功能",
    description: "这是一个新功能",
  },
};
```

在 `src/i18n/locales/en-US.ts` 中添加对应的英文翻译：

```typescript
export default {
  // ... 现有翻译
  newFeature: {
    title: "New Feature",
    description: "This is a new feature",
  },
};
```

### 2. 在模板中使用

```vue
<template>
  <div>
    <h2>{{ $t('newFeature.title') }}</h2>
    <p>{{ $t('newFeature.description') }}</p>
  </div>
</template>
```

## 翻译键的命名规范

- 使用点号（.）分隔层级
- 使用小驼峰命名（camelCase）
- 按功能分组（如 nav, home, admin, errors 等）

示例：
```
nav.articles        # 导航菜单 - 文章
home.title         # 主页 - 标题
admin.users        # 管理后台 - 用户
errors.notFound    # 错误 - 未找到
```

## 当前支持的语言

| 语言 | 代码 |
|------|------|
| 中文(简体) | zh-CN |
| 英文 | en-US |

## 浏览器语言检测

应用启动时会自动检测浏览器语言：
- 浏览器语言为中文 → 加载中文
- 浏览器语言为英文 → 加载英文
- 其他语言 → 默认使用英文

## 本地化持久化

用户选择的语言会保存在 localStorage 中，键为 `language`。下次访问时会自动加载用户之前的语言选择。

## 添加新语言

要添加新语言（如日文），需要：

1. 在 `src/i18n/locales/` 下创建新文件 `ja-JP.ts`
2. 在 `src/i18n/index.ts` 中导入并注册
3. 在 `TopBar.vue` 的 `languages` 对象中添加新语言

示例：

```typescript
// src/i18n/index.ts
import jaJP from "./locales/ja-JP";

// 在 messages 中添加
messages: {
  "zh-CN": zhCN,
  "en-US": enUS,
  "ja-JP": jaJP,  // 新增
}

// TopBar.vue
const languages = {
  "zh-CN": "中文",
  "en-US": "English",
  "ja-JP": "日本語",  // 新增
};
```

## 常见问题

### Q: 如何在组件中访问当前语言？
A: 使用 `useI18n()` 钩子获取 `locale` 值：
```typescript
const { locale } = useI18n();
console.log(locale.value); // 'zh-CN' 或 'en-US'
```

### Q: 翻译中如何使用变量？
A: 使用 `t` 函数的第二个参数传递变量：
```typescript
// zh-CN.ts
greeting: "你好，{name}！"

// 使用
t('greeting', { name: '张三' })
```

### Q: 如何处理复数形式？
A: vue-i18n 支持复数，详见[官方文档](https://vue-i18n.intlify.dev/guide/essentials/pluralization.html)

## 相关文档

- [vue-i18n 官方文档](https://vue-i18n.intlify.dev/)
- [项目多语言配置](./src/i18n/index.ts)
