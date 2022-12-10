import { createApp, h } from "vue";
import App from "./App.vue";
import { ApolloClient, HttpLink, InMemoryCache } from "@apollo/client/core";
import { createApolloProvider } from "@vue/apollo-option";

// Http endpoint
const httpEndpoint = SERVER_HTTP || "http://localhost:8081/query";
// const httpEndpoint = 'http://localhost:8081/query'

const httpLink = new HttpLink({
  // You should use an absolute URL here
  uri: httpEndpoint,
});

// Create the apollo client
const apolloClient = new ApolloClient({
  link: httpLink,
  cache: new InMemoryCache(),
  connectToDevTools: true,
});

// Create a provider
const apolloProvider = createApolloProvider({
  defaultClient: apolloClient,
});

// Use the provider
const app = createApp({
  render: () => h(App),
});
app.use(apolloProvider);

// TODO(chris): replace the following code with something like here
// https://www.apollographql.com/docs/react/networking/authentication/

// // Manually call this when user log in
// export async function onLogin(apolloClient, token) {
//   if (typeof localStorage !== 'undefined' && token) {
//     localStorage.setItem(AUTH_TOKEN, token)
//   }
//   if (apolloClient.wsClient) restartWebsockets(apolloClient.wsClient)
//   try {
//     await apolloClient.resetStore()
//   } catch (e) {
//     // eslint-disable-next-line no-console
//     console.log('%cError on cache reset (login)', 'color: orange;', e.message)
//   }
// }

// // Manually call this when user log out
// export async function onLogout(apolloClient) {
//   if (typeof localStorage !== 'undefined') {
//     localStorage.removeItem(AUTH_TOKEN)
//   }
//   if (apolloClient.wsClient) restartWebsockets(apolloClient.wsClient)
//   try {
//     await apolloClient.resetStore()
//   } catch (e) {
//     // eslint-disable-next-line no-console
//     console.log('%cError on cache reset (logout)', 'color: orange;', e.message)
//   }
// }
