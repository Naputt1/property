import { defineConfig } from '@kubb/core'
import { pluginOas } from '@kubb/plugin-oas'
import { pluginClient } from '@kubb/plugin-client'
import { pluginReactQuery } from '@kubb/plugin-react-query'
import { pluginTs } from '@kubb/plugin-ts'
import { pluginZod } from '@kubb/plugin-zod'

export default defineConfig({
  root: '.',
  input: {
    path: '../backend/docs/swagger.json',
  },
  output: {
    path: './src/gen',
    clean: true,
  },
  plugins: [
    pluginOas(),
    pluginTs({
      output: {
        path: './models',
      },
    }),
    pluginClient({
      output: {
        path: './clients',
      },
      importPath: '../../services/kubb-client',
    }),
    pluginReactQuery({
      output: {
        path: './hooks',
      },
      client: {
        importPath: '../../services/kubb-client',
      },
    }),
    pluginZod({
      output: {
        path: './zod',
      },
    }),
  ],
})
