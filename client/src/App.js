import React from "react";
import {BrowserRouter as Router, Route} from "react-router-dom";

import routes from "./routes";

import "bootstrap/dist/css/bootstrap.min.css";
import "./styles/css/cloud-tides.css";
import "./styles/css/tides-icons/tides_icons.css"
import {PrivateRoute} from "./utils/routerExt";


export default () => (
    <Router basename={process.env.REACT_APP_BASENAME || ""}>
        <div>
            {routes.map((route, index) => {
                return (
                    <div>
                        {(route.path == "/signup" || route.path == "/login") ? (
                            <Route
                                key={index}
                                path={route.path}
                                exact={route.exact}
                                component={(props) => (
                                    <route.layout {...props}>
                                        <route.component {...props} />
                                    </route.layout>
                                )}
                            />
                        ) : (
                            <PrivateRoute
                                key={index}
                                path={route.path}
                                exact={route.exact}
                                component={(props) => (
                                    <route.layout {...props}>
                                        <route.component {...props} />
                                    </route.layout>
                                )}
                            />
                        )}
                    </div>
                );
            })}
        </div>
    </Router>
);
