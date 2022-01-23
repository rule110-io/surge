import Dashboard from "@/layout/Dashboard/Dashboard.vue";
import Transfers from "@/pages/Transfers.vue";
import Discover from "@/pages/Discover.vue";

const routes = [
  {
    path: "/",
    component: Dashboard,
    redirect: "/transfers",
    children: [
      {
        path: "transfers",
        name: "transfers",
        component: Transfers,
      },
      {
        path: "discover",
        name: "discover",
        component: Discover,
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
