import { defineConfig } from 'orval';

// Orval 7.21 has no baseUrl.runtime; API_BASE_URL is applied in custom-fetch.ts.
export default defineConfig({
  api: {
    input: './openapi.json',
    output: {
      target: './src/generated.ts',
      client: 'fetch',
      mode: 'single',
      override: {
        mutator: {
          path: './src/custom-fetch.ts',
          name: 'customFetch',
        },
        fetch: {
          includeHttpResponseReturnType: true,
        },
      },
    },
  },
});
