import type { CodegenConfig } from "@graphql-codegen/cli";

const config: CodegenConfig = {
  schema: "../backend/internal/graph/schema.graphqls",
  documents: ["src/graphql/**/*.graphql"],
  generates: {
    "src/gen-gql/graphql.ts": {
      plugins: [
        {
          add: {
            content: 'import { graphqlFetcher } from "../lib/graphql-fetcher";',
          },
        },
        "typescript",
        "typescript-operations",
        "typescript-react-query",
      ],
      config: {
        useTypeImports: true,
        reactQueryVersion: 5,
        fetcher: {
          func: "graphqlFetcher",
          isReactHook: false,
        },
      },
    },
  },
};

export default config;
