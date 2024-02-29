import { build } from 'esbuild'

await build({
  logLevel: 'info',
  entryPoints: ['web/index.*'],
  bundle: true,
  minify: true,
  loader: {
    '.woff': 'file',
    '.woff2': 'file'
  },
  outExtension: {
    '.css': '.min.css',
    '.js': '.min.js'
  },
  outdir: 'internal/web/root/static/bundle'
})
