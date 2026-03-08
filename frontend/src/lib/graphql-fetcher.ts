export const graphqlFetcher = <TData, TVariables>(
  query: string,
  variables?: TVariables,
) => {
  return async (): Promise<TData> => {
    const token = localStorage.getItem("token");
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };

    if (token && token !== "null") {
      headers["Authorization"] = `Bearer ${token}`;
    }

    const res = await fetch("/api/query", {
      method: "POST",
      headers,
      body: JSON.stringify({ query, variables }),
      credentials: "include",
    });

    const json = await res.json();

    if (!res.ok) {
      throw new Error(
        json.error || json.message || `HTTP error! status: ${res.status}`,
      );
    }

    if (json.errors) {
      const { message } = json.errors[0];
      throw new Error(message);
    }

    return json.data;
  };
};
