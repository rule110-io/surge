import Dashboard from "@/layout/Dashboard/Dashboard.vue";
import Search from "@/pages/Search.vue";

const routes = [
  {
    path: "/",
    component: Dashboard,
    redirect: "/search",
    children: [
      {
        path: "search",
        name: "search",
        component: Search,
      },
    ],
  },
];

/**
 * Asynchronously load view (Webpack Lazy loading compatible)
 * The specified component must be inside the Views folder
 * @param  {string} name  the filename (basename) of the view to load.
function view(name) {
   var res= require('../components/Dashboard/Views/' + name + '.vue');
   return res;
};**/

export default routes;
