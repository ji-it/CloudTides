import React from "react";
import {Redirect} from "react-router-dom";

// Layout Types
import {DefaultLayout} from "./layouts";

// Route Views
import {STATE_LOGIN, STATE_SIGNUP} from "./components/auth/AuthForn";
import AuthPage from "./views/AuthPage";
import Dashboard from "./views/Dashboard";
import ManageResources from "./views/ManageResources";
import Errors from "./views/Errors";
import EmptyLayout from "./layouts/Empty";
// import UserProfileLite from "./views/UserProfileLite";
// import AddNewPost from "./views/AddNewPost";
// import ComponentsOverview from "./views/ComponentsOverview";
// import Tables from "./views/Tables";
// import BlogPosts from "./views/BlogPosts";
export default [
    // {
    //     path: "/",
    //     exact: true,
    //     layout: DefaultLayout,
    //     component: () => <Redirect to="/blog-overview"/>
    // },
    {
        path: "/",
        exact: true,
        layout: DefaultLayout,
        component: Errors
    },
    {
        path: "/login",
        layout: EmptyLayout,
        component: props => <AuthPage {...props} authState={STATE_LOGIN}/>
    },
    {
        path: "/signup",
        layout: EmptyLayout,
        component: props => <AuthPage {...props} authState={STATE_SIGNUP}/>
    },
    {
        path: "/home",
        layout: DefaultLayout,
        component: Dashboard
    },
    {
        path: "/manage-resources",
        layout: DefaultLayout,
        component: ManageResources
    },
    {
        path: "/errors",
        layout: DefaultLayout,
        component: Errors
    },
    // {
    //   path: "/components-overview",
    //   layout: DefaultLayout,
    //   component: ComponentsOverview
    // },
    // {
    //   path: "/tables",
    //   layout: DefaultLayout,
    //   component: Tables
    // },
    // {
    //   path: "/blog-posts",
    //   layout: DefaultLayout,
    //   component: BlogPosts
    // }
];
