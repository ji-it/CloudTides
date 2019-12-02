import React, {Component} from "react";
import {BrowserRouter as Router, Route, Switch} from "react-router-dom";

import "bootstrap/dist/css/bootstrap.min.css";
import "./styles/css/cloud-tides.css";
import "./styles/css/tides-icons/tides_icons.css"
import PrivateRoute from "./utils/routerExt";

import {STATE_LOGIN, STATE_SIGNUP} from "./components/auth/AuthForm";

import {DefaultLayout, EmptyLayout} from "./layouts";

import Dashboard from "./views/Dashboard";
import Errors from "./views/Errors";
import AuthPage from "./views/AuthPage";
import ManageResources from "./views/ManageResources";
import Contribution from "./views/Contribution";
import Templates from "./views/Templates";
import Settings from "./views/Settings";

class App extends Component {
    render() {
        return (
            <Router>
                <div className="App">
                    <Switch>
                        {/* A user can't go to the HomePage if is not authenticated */}
                        <PrivateRoute
                            path="/"
                            component={(props) => (
                                <DefaultLayout {...props}>
                                    <Dashboard {...props} />
                                </DefaultLayout>
                            )} exact
                        />
                        <Route
                            path="/login"
                            component={(props) => (
                                <EmptyLayout {...props}>
                                    <AuthPage {...props} authState={STATE_LOGIN}/>
                                </EmptyLayout>
                            )}
                        />
                        <Route
                            path="/signup"
                            component={(props) => (
                                <EmptyLayout {...props}>
                                    <AuthPage {...props} authState={STATE_SIGNUP}/>
                                </EmptyLayout>
                            )}
                        />
                        <PrivateRoute
                            path="/home"
                            component={(props) => (
                                <DefaultLayout {...props}>
                                    <Dashboard {...props} />
                                </DefaultLayout>
                            )}
                        />
                        <PrivateRoute
                            path="/manage-resources"
                            component={(props) => (
                                <DefaultLayout {...props}>
                                    <ManageResources {...props} />
                                </DefaultLayout>
                            )}
                        />
                        <PrivateRoute
                            path="/contribution"
                            component={(props) => (
                                <DefaultLayout {...props}>
                                    <Contribution {...props} />
                                </DefaultLayout>
                            )}
                        />
                        <PrivateRoute
                            path="/templates"
                            component={(props) => (
                                <DefaultLayout {...props}>
                                    <Templates {...props} />
                                </DefaultLayout>
                            )}
                        />
                        <PrivateRoute
                            path="/settings"
                            component={(props) => (
                                <DefaultLayout {...props}>
                                    <Settings {...props} />
                                </DefaultLayout>
                            )}
                        />
                        <Route path="" component={Errors}/>
                    </Switch>
                </div>
            </Router>
        );
    }
}

export default App;