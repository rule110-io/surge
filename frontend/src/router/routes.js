import Dashboard from "@/layout/Dashboard/Dashboard.vue";
import Explore from "@/pages/Explore.vue";
import Search from "@/pages/Search.vue";
import Download from "@/pages/Download.vue";
import Settings from "@/pages/Settings.vue";

const routes = [
  {
    path: "/",
    component: Dashboard,
    redirect: "/explore",
    children: [
      {
        path: "explore",
        name: "explore",
        component: Explore,
      },
      {
        path: "search",
        name: "search",
        component: Search,
      },
      {
        path: "download",
        name: "download",
        component: Download,
      },
      {
        path: "settings",
        name: "settings",
        component: Settings,
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
