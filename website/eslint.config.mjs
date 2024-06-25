import globals from "globals";
import pluginJs from "@eslint/js";
import * as tseslint from "typescript-eslint";
import pluginReactConfig from "eslint-plugin-react/configs/recommended.js";
import prettierEslintPlugin from "eslint-plugin-prettier";
import docusaurusEslintPlugin from "@docusaurus/eslint-plugin";

export default [
  // TODO: eslint configure for .md, .mdx files
  {files: ["*/.{js,mjs,cjs,ts,jsx,tsx}"]},
  {ignores: ["babel.config.js", "eslint.config.mjs", "docusaurus.config.ts", ".docusaurus/*", "build/*"]},
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
      },
    },
  },
  pluginJs.configs.recommended,
  ...tseslint.configs.recommended,
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
      "prettier/prettier": "off",
      "@docusaurus/string-literal-i18n-messages": "warn",
      "@docusaurus/no-untranslated-text": "warn",
      "@docusaurus/no-html-links": "warn",
      "react/react-in-jsx-scope": "off",
    },
  },
];