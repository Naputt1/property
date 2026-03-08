import { graphqlFetcher } from "../lib/graphql-fetcher";
import { useQuery, useMutation, type UseQueryOptions, type UseMutationOptions } from '@tanstack/react-query';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Map: { input: any; output: any; }
  Time: { input: any; output: any; }
};

export type Job = {
  __typename?: 'Job';
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  message: Scalars['String']['output'];
  progress: Scalars['Int']['output'];
  status: Scalars['String']['output'];
  taskType: Scalars['String']['output'];
  total: Scalars['Int']['output'];
  updatedAt: Scalars['Time']['output'];
};

export type JobConnection = {
  __typename?: 'JobConnection';
  items: Array<Job>;
  total: Scalars['Int']['output'];
};

export type Mutation = {
  __typename?: 'Mutation';
  createProperty: Property;
  deleteProperty: Scalars['Boolean']['output'];
  updateProperty: Property;
};


export type MutationCreatePropertyArgs = {
  input: PropertyInput;
};


export type MutationDeletePropertyArgs = {
  id: Scalars['ID']['input'];
};


export type MutationUpdatePropertyArgs = {
  id: Scalars['ID']['input'];
  input: PropertyInput;
};

export type Property = {
  __typename?: 'Property';
  county: Scalars['String']['output'];
  createdAt: Scalars['Time']['output'];
  dateOfTransfer: Scalars['Time']['output'];
  district: Scalars['String']['output'];
  duration: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  locality: Scalars['String']['output'];
  oldNew: Scalars['String']['output'];
  paon: Scalars['String']['output'];
  postcode: Scalars['String']['output'];
  ppdCategoryType: Scalars['String']['output'];
  price: Scalars['Int']['output'];
  propertyType: Scalars['String']['output'];
  recordStatus: Scalars['String']['output'];
  saon: Scalars['String']['output'];
  street: Scalars['String']['output'];
  townCity: Scalars['String']['output'];
  updatedAt: Scalars['Time']['output'];
};

export type PropertyConnection = {
  __typename?: 'PropertyConnection';
  items: Array<Property>;
  total: Scalars['Int']['output'];
};

export type PropertyInput = {
  county: Scalars['String']['input'];
  dateOfTransfer: Scalars['Time']['input'];
  district: Scalars['String']['input'];
  duration: Scalars['String']['input'];
  locality: Scalars['String']['input'];
  oldNew: Scalars['String']['input'];
  paon: Scalars['String']['input'];
  postcode: Scalars['String']['input'];
  ppdCategoryType: Scalars['String']['input'];
  price: Scalars['Int']['input'];
  propertyType: Scalars['String']['input'];
  recordStatus: Scalars['String']['input'];
  saon: Scalars['String']['input'];
  street: Scalars['String']['input'];
  townCity: Scalars['String']['input'];
};

export type Query = {
  __typename?: 'Query';
  jobs: JobConnection;
  properties: PropertyConnection;
  property?: Maybe<Property>;
};


export type QueryJobsArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
};


export type QueryPropertiesArgs = {
  county?: InputMaybe<Scalars['String']['input']>;
  limit?: InputMaybe<Scalars['Int']['input']>;
  maxPrice?: InputMaybe<Scalars['Int']['input']>;
  minPrice?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
  postcode?: InputMaybe<Scalars['String']['input']>;
  propertyType?: InputMaybe<Scalars['String']['input']>;
  townCity?: InputMaybe<Scalars['String']['input']>;
};


export type QueryPropertyArgs = {
  id: Scalars['ID']['input'];
};

export type GetAdminJobsQueryVariables = Exact<{
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
}>;


export type GetAdminJobsQuery = { __typename?: 'Query', jobs: { __typename?: 'JobConnection', total: number, items: Array<{ __typename?: 'Job', id: string, createdAt: any, updatedAt: any, taskType: string, status: string, message: string, progress: number, total: number }> } };

export type GetPropertiesQueryVariables = Exact<{
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
  minPrice?: InputMaybe<Scalars['Int']['input']>;
  maxPrice?: InputMaybe<Scalars['Int']['input']>;
  postcode?: InputMaybe<Scalars['String']['input']>;
  townCity?: InputMaybe<Scalars['String']['input']>;
  county?: InputMaybe<Scalars['String']['input']>;
  propertyType?: InputMaybe<Scalars['String']['input']>;
}>;


export type GetPropertiesQuery = { __typename?: 'Query', properties: { __typename?: 'PropertyConnection', total: number, items: Array<{ __typename?: 'Property', id: string, price: number, dateOfTransfer: any, postcode: string, propertyType: string, oldNew: string, paon: string, saon: string, street: string, townCity: string, county: string }> } };

export type CreatePropertyMutationVariables = Exact<{
  input: PropertyInput;
}>;


export type CreatePropertyMutation = { __typename?: 'Mutation', createProperty: { __typename?: 'Property', id: string, price: number, street: string } };

export type UpdatePropertyMutationVariables = Exact<{
  id: Scalars['ID']['input'];
  input: PropertyInput;
}>;


export type UpdatePropertyMutation = { __typename?: 'Mutation', updateProperty: { __typename?: 'Property', id: string, price: number, street: string } };

export type DeletePropertyMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeletePropertyMutation = { __typename?: 'Mutation', deleteProperty: boolean };



export const GetAdminJobsDocument = `
    query GetAdminJobs($limit: Int, $offset: Int) {
  jobs(limit: $limit, offset: $offset) {
    items {
      id
      createdAt
      updatedAt
      taskType
      status
      message
      progress
      total
    }
    total
  }
}
    `;

export const useGetAdminJobsQuery = <
      TData = GetAdminJobsQuery,
      TError = unknown
    >(
      variables?: GetAdminJobsQueryVariables,
      options?: Omit<UseQueryOptions<GetAdminJobsQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<GetAdminJobsQuery, TError, TData>['queryKey'] }
    ) => {
    
    return useQuery<GetAdminJobsQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['GetAdminJobs'] : ['GetAdminJobs', variables],
    queryFn: graphqlFetcher<GetAdminJobsQuery, GetAdminJobsQueryVariables>(GetAdminJobsDocument, variables),
    ...options
  }
    )};

export const GetPropertiesDocument = `
    query GetProperties($limit: Int, $offset: Int, $minPrice: Int, $maxPrice: Int, $postcode: String, $townCity: String, $county: String, $propertyType: String) {
  properties(
    limit: $limit
    offset: $offset
    minPrice: $minPrice
    maxPrice: $maxPrice
    postcode: $postcode
    townCity: $townCity
    county: $county
    propertyType: $propertyType
  ) {
    items {
      id
      price
      dateOfTransfer
      postcode
      propertyType
      oldNew
      paon
      saon
      street
      townCity
      county
    }
    total
  }
}
    `;

export const useGetPropertiesQuery = <
      TData = GetPropertiesQuery,
      TError = unknown
    >(
      variables?: GetPropertiesQueryVariables,
      options?: Omit<UseQueryOptions<GetPropertiesQuery, TError, TData>, 'queryKey'> & { queryKey?: UseQueryOptions<GetPropertiesQuery, TError, TData>['queryKey'] }
    ) => {
    
    return useQuery<GetPropertiesQuery, TError, TData>(
      {
    queryKey: variables === undefined ? ['GetProperties'] : ['GetProperties', variables],
    queryFn: graphqlFetcher<GetPropertiesQuery, GetPropertiesQueryVariables>(GetPropertiesDocument, variables),
    ...options
  }
    )};

export const CreatePropertyDocument = `
    mutation CreateProperty($input: PropertyInput!) {
  createProperty(input: $input) {
    id
    price
    street
  }
}
    `;

export const useCreatePropertyMutation = <
      TError = unknown,
      TContext = unknown
    >(options?: UseMutationOptions<CreatePropertyMutation, TError, CreatePropertyMutationVariables, TContext>) => {
    
    return useMutation<CreatePropertyMutation, TError, CreatePropertyMutationVariables, TContext>(
      {
    mutationKey: ['CreateProperty'],
    mutationFn: (variables?: CreatePropertyMutationVariables) => graphqlFetcher<CreatePropertyMutation, CreatePropertyMutationVariables>(CreatePropertyDocument, variables)(),
    ...options
  }
    )};

export const UpdatePropertyDocument = `
    mutation UpdateProperty($id: ID!, $input: PropertyInput!) {
  updateProperty(id: $id, input: $input) {
    id
    price
    street
  }
}
    `;

export const useUpdatePropertyMutation = <
      TError = unknown,
      TContext = unknown
    >(options?: UseMutationOptions<UpdatePropertyMutation, TError, UpdatePropertyMutationVariables, TContext>) => {
    
    return useMutation<UpdatePropertyMutation, TError, UpdatePropertyMutationVariables, TContext>(
      {
    mutationKey: ['UpdateProperty'],
    mutationFn: (variables?: UpdatePropertyMutationVariables) => graphqlFetcher<UpdatePropertyMutation, UpdatePropertyMutationVariables>(UpdatePropertyDocument, variables)(),
    ...options
  }
    )};

export const DeletePropertyDocument = `
    mutation DeleteProperty($id: ID!) {
  deleteProperty(id: $id)
}
    `;

export const useDeletePropertyMutation = <
      TError = unknown,
      TContext = unknown
    >(options?: UseMutationOptions<DeletePropertyMutation, TError, DeletePropertyMutationVariables, TContext>) => {
    
    return useMutation<DeletePropertyMutation, TError, DeletePropertyMutationVariables, TContext>(
      {
    mutationKey: ['DeleteProperty'],
    mutationFn: (variables?: DeletePropertyMutationVariables) => graphqlFetcher<DeletePropertyMutation, DeletePropertyMutationVariables>(DeletePropertyDocument, variables)(),
    ...options
  }
    )};
