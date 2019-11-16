import React from "react";
import {Route, Redirect} from "react-router-dom";

import {ISAUTHENTICATED} from "../components/auth/AuthForn";

export const PrivateRoute = ({component: Component, ...rest}) => (
    <Route {...rest} render={(props) => (
         <Component {...props} />
        // ISAUTHENTICATED === true
        //     ? <Component {...props} />
        //     : <Redirect to={{
        //         pathname: '/login',
        //         state: {from: props.location}
        //     }}/>
    )}/>
);
