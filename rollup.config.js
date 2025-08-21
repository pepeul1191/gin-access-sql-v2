import svelte from 'rollup-plugin-svelte';
import commonjs from '@rollup/plugin-commonjs';
import terser from '@rollup/plugin-terser';
import resolve from '@rollup/plugin-node-resolve';
import css from 'rollup-plugin-css-only';
import copy from 'rollup-plugin-copy';
import postcss from 'rollup-plugin-postcss';
import { nodeResolve } from '@rollup/plugin-node-resolve';

const production = !process.env.ROLLUP_WATCH;

const Vendor = {
  input: 'src/entries/vendor.js',
  output: {
    sourcemap: true,
    format: 'iife',
    name: 'vendor',
    file: production ? 'static/dist/vendor.min.js' : 'static/dist/vendor.js',
    globals: {
      'axios': 'axios',
      'bootstrap': 'bootstrap'
    }
  },
  plugins: [
    svelte({
      compilerOptions: {
        dev: !production
      }
    }),

    css({
      output: 'vendor.min.css',
      minify: true
    }),

    resolve({
      browser: true,
      dedupe: ['svelte'],
      exportConditions: ['svelte']
    }),

    commonjs(),

    production && terser(),

    copy({
      hook: 'writeBundle',
      targets: [
        {
          src: 'node_modules/font-awesome/fonts/*',
          dest: 'static/fonts/'
        },
        {
          src: 'node_modules/simplemde/dist/*',
          dest: 'static/dist/'
        }
      ]
    })
  ],
  watch: {
    clearScreen: false
  }
};

const App = {
  input: 'src/entries/app.js', // Tu archivo de entrada principal
  output: {
    file: 'static/dist/app.min.js',
    format: 'iife',
    sourcemap: !production
  },
  plugins: [
    // Procesamiento de CSS
    postcss({
      extract: true, // Extrae CSS a archivo separado
      minimize: production,
      plugins: [
        require('autoprefixer')() // Opcional: agrega prefijos vendor
      ]
    }),

    // Resolución de módulos
    nodeResolve({
      browser: true
    }),

    // Soporte para CommonJS
    commonjs(),

    // Minificación en producción
    production && terser(),

    // Copia archivos estáticos (opcional)
    copy({
      targets: [
        { src: 'src/assets/*', dest: 'static/dist/assets' }
      ]
    })
  ]
};

export default [App, Vendor, ];