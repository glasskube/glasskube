import globals from "globals";
import pluginJs from "@eslint/js";
import * as tseslint from "typescript-eslint";
import pluginReactConfig from "eslint-plugin-react/configs/recommended.js";
import prettierEslintPlugin from "eslint-plugin-prettier";
import docusaurusEslintPlugin from "@docusaurus/eslint-plugin";
import pluginPrettierConfig from "eslint-config-prettier";

/** @type {import('eslint').Linter.FlatConfig[]} */
export default [
  {ignores: [".docusaurus/*", "build/*"]},
  {
    plugins: {
      prettier: prettierEslintPlugin,
      "@docusaurus": docusaurusEslintPlugin,
    },
  },
  {
    languageOptions: {
      parser: tseslint.parser, 
      parserOptions: {
        ecmaFeatures: {jsx: true},
        sourceType: "module"
      },
      globals: {
        ...globals.browser,
        ...globals.es2021,
        ...globals.node 
      },
    },
  },
  pluginJs.configs.recommended,
  ...tseslint.configs.recommended,
  pluginPrettierConfig,
  pluginReactConfig,
  {
    settings: {
      react: {
        version: "detect",
      },
    },
  },
  {
    rules: {
      "prettier/prettier": "warn",
      "@docusaurus/string-literal-i18n-messages": "warn",
      "@docusaurus/no-untranslated-text": "warn",
      "@docusaurus/no-html-links": "warn",
      "react/react-in-jsx-scope": "off"
    },
  },
];